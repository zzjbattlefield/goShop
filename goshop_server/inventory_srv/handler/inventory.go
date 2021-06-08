package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goshop/inventory_srv/global"
	"goshop/inventory_srv/model"
	"goshop/inventory_srv/proto"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-redsync/redsync/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Inventort struct {
	proto.UnimplementedInventoryServer
}

//设置(更新)库存
func (i *Inventort) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*empty.Empty, error) {
	var invModel model.Inventory
	if res := global.DB.Where("goods = ?", req.GoodsId).First(&invModel); res.RowsAffected == 0 {
		invModel.Goods = req.GoodsId
	}
	invModel.Stocks = req.Num
	global.DB.Save(&invModel)
	return &empty.Empty{}, nil
}

func (i *Inventort) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var invModel model.Inventory
	if res := global.DB.Where("goods=?", req.GoodsId).First(&invModel); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存信息不存在")
	} else {
		return &proto.GoodsInvInfo{GoodsId: invModel.ID, Num: invModel.Stocks}, nil
	}
}

//扣减库存(乐观锁和悲观锁)
// func (i *Inventort) Sell(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
// 	tx := global.DB.Begin()
// 	for _, goodInfo := range req.GoodsInvInfo {
// 		for {
// 			var invModel model.Inventory
// 			//悲观锁 在查询条件是索引的情况下是行级锁,不然是表级锁
// 			// if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("goods=?", goodInfo.GoodsId).Find(&invModel); res.RowsAffected == 0 {
// 			// 	tx.Rollback()
// 			// 	return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
// 			// }
// 			// 乐观锁版本 使用表中的Level字段来控制并发
// 			if res := global.DB.Where("goods=?", goodInfo.GoodsId).Find(&invModel); res.RowsAffected == 0 {
// 				tx.Rollback()
// 				return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
// 			}
// 			if invModel.Stocks < goodInfo.Num {
// 				tx.Rollback()
// 				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
// 			}
// 			//扣减
// 			invModel.Stocks -= goodInfo.Num
// 			if res := tx.Model(&model.Inventory{}).Select("stocks", "version").Where("goods = ?", goodInfo.GoodsId).Where("version=?", invModel.Version).Updates(model.Inventory{Stocks: invModel.Stocks, Version: invModel.Version + 1}); res.RowsAffected == 0 {
// 				zap.S().Info("库存扣减失败开始重试")
// 			} else {
// 				break
// 			}
// 		}
// 	}
// 	tx.Commit()
// 	return &emptypb.Empty{}, nil
// }

//扣减库存(分布式redis锁)
func (i *Inventort) Sell(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	tx := global.DB.Begin()
	pool := global.RedisPool
	rs := redsync.New(pool)
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1, //默认已经扣减
		Detail:  nil,
	}
	var detail []model.GoodsDetail
	for _, goodInfo := range req.GoodsInvInfo {
		detail = append(detail, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})
		//遍历每个订单内的货物如果有一个库存不足时整个订单库存的扣减都要回滚
		mutxName := fmt.Sprintf("lock_goods_%d", goodInfo.GoodsId)
		mutex := rs.NewMutex(mutxName)
		if err := mutex.Lock(); err != nil {
			return nil, status.Error(codes.Internal, "初始化分布式锁失败")
		}
		var invModel model.Inventory
		if res := global.DB.Where("goods=?", goodInfo.GoodsId).Find(&invModel); res.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		}
		if invModel.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减
		invModel.Stocks -= goodInfo.Num
		tx.Save(&invModel)
		if ok, err := mutex.Unlock(); err != nil || !ok {
			return nil, status.Error(codes.Internal, "解除分布式锁失败")
		}
	}
	sellDetail.Detail = detail
	//写入sellDetail表
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

//库存归还
func (i *Inventort) Reback(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInvInfo {
		var invModel model.Inventory
		if res := global.DB.Where("goods=?", goodInfo.GoodsId).Find(&invModel); res.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		}
		//扣减
		invModel.Stocks += goodInfo.Num
		tx.Updates(&invModel)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

//自动归还库存(通过rocketMq)
func AutoReback(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	//定义在orderServer下 不建议跨微服务引用 重新定义一个
	type orderInfo struct {
		OrderSn string
	}
	for i := range msg {
		//归还库存 需要知道订单id和归还件数
		var orderInfo orderInfo
		if err := json.Unmarshal(msg[i].Body, &orderInfo); err != nil {
			zap.S().Errorf("解析json失败:%s", err)
			return consumer.ConsumeSuccess, nil
		}
		//将inv的库存加回去,将sellDetail的状态设置为2 防止重复归还 不同的表 放在事务运行
		tx := global.DB.Begin()
		var sellDetail = model.StockSellDetail{}
		if result := global.DB.Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
			//已经归还过了
			return consumer.ConsumeSuccess, nil
		}
		//如果查询出记录 则逐个归还库存
		for _, orderGood := range sellDetail.Detail {
			if res := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); res.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		if res := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status", 2); res.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
