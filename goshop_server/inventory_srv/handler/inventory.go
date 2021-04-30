package handler

import (
	"context"
	"fmt"
	"goshop/inventory_srv/global"
	"goshop/inventory_srv/model"
	"goshop/inventory_srv/proto"

	"github.com/go-redsync/redsync/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	for _, goodInfo := range req.GoodsInvInfo {
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
