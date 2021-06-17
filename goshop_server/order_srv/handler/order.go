package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goshop/order_srv/global"
	"goshop/order_srv/model"
	"goshop/order_srv/proto"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

//获得当前用户购物车中所有的商品
func (o *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shoppingcart []model.ShoppingCart
	var cartListInfo proto.CartItemListResponse
	if res := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shoppingcart); res.Error != nil {
		return nil, res.Error
	} else {
		cartListInfo.Total = int32(res.RowsAffected)
	}
	var shopCarInfo []*proto.ShopCartInfoResponse
	for _, cartInfo := range shoppingcart {
		shopCarInfo = append(shopCarInfo, &proto.ShopCartInfoResponse{
			Id:      cartInfo.ID,
			UserId:  cartInfo.User,
			Num:     cartInfo.Nums,
			GoodsId: cartInfo.Goods,
			Checked: cartInfo.Checked,
		})
	}
	cartListInfo.Data = shopCarInfo
	return &cartListInfo, nil
}

//新增购物车商品
func (o *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//1.如果之前购物车没有此商品那就新建 如果已经存在则需要合并购物车记录
	var shopCart model.ShoppingCart
	if res := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&shopCart); res.RowsAffected == 1 {
		shopCart.Nums += req.Num
	} else {
		shopCart.Checked = req.Checked
		shopCart.Goods = req.GoodsId
		shopCart.User = req.UserId
		shopCart.Nums = req.Num
	}
	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

//更新购物车商品
func (o *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	//主要更新选中状态 / 商品数量
	var shopCart model.ShoppingCart
	shopCart.Nums = req.Num
	shopCart.Checked = req.Checked
	if req := global.DB.Model(&shopCart).Where("user=? AND goods=? ", req.UserId, req.GoodsId).Updates(&shopCart); req.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "更新失败")
	}
	return &empty.Empty{}, nil
}

func (o *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	if res := global.DB.Where("goods=? AND user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &empty.Empty{}, nil
}

//订单相关的功能

type orderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
}

//RocketMQ本地事务消息执行逻辑
func (o *orderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	//通过message的body获取订单相关信息
	_ = json.Unmarshal(msg.Body, &orderInfo)
	var shopCart []model.ShoppingCart
	var goodIds []int32
	//检查购物车是否有选择商品
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCart); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算商品"
		return primitive.RollbackMessageState
	}
	goodsNumMap := make(map[int32]int32)
	for _, shopCartInfo := range shopCart {
		goodIds = append(goodIds, shopCartInfo.Goods)
		goodsNumMap[shopCartInfo.Goods] = shopCartInfo.Nums
	}
	//跨服务调用（商品查询）
	rsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodIds})
	if err != nil {
		o.Code = codes.InvalidArgument
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range rsp.Data {
		orderAmount += good.ShopPrice * float32(goodsNumMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumMap[good.Id],
		})
	}
	//跨微服务调用（扣减库存）
	//TODO: 分布式事务
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInvInfo: goodsInvInfo, OrderSn: orderInfo.OrderSn}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "批量扣减失败"
		return primitive.RollbackMessageState
	}
	//创建订单的基本信息
	tx := global.DB.Begin()
	//订单总金额
	orderInfo.OrderMount = orderAmount
	if res := tx.Save(&orderInfo); res.RowsAffected == 0 {
		o.Code = codes.ResourceExhausted
		o.Detail = "保存订单失败"
		return primitive.CommitMessageState
	}
	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	//批量插入
	if res := tx.CreateInBatches(&orderGoods, 100); res.RowsAffected == 0 {
		o.Code = codes.Internal
		o.Detail = "批量插入失败"
		return primitive.CommitMessageState
	}

	//删除下单了的购物车记录
	if res := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); res.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}
	//发送订单延时消息,超过三十分钟归还库存
	p, err := NewProducer()
	if err != nil {
		tx.Rollback()
		o.Code = codes.Internal
		return primitive.CommitMessageState
	}
	layoutMsg := primitive.NewMessage("order_timeout", msg.Body)
	layoutMsg.WithDelayTimeLevel(3)
	_, err = p.SendSync(context.Background(), layoutMsg)
	if err != nil {
		zap.S().Errorf("发送失败:%v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延迟消息失败"
		return primitive.CommitMessageState
	}
	// defer p.Shutdown()
	//提交本地事务
	tx.Commit()
	o.Code = codes.OK
	//本地事务执行成功 取消归还库存事务操作
	return primitive.RollbackMessageState
}

//RocketMQ事务消息回查逻辑
func (o *orderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderModel model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderModel)
	res := global.DB.Where(&model.OrderInfo{OrderSn: orderModel.OrderSn}).First(&model.OrderInfo{})
	if res.RowsAffected == 0 {
		//没找到订单号数据说明扣减的时候回滚了,不一定代表扣减过库存,要在归还库存的接口保证幂等性(需要commit归还库存接口)
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

//创建订单
func (o *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单TODO:
			1.先从购物车中拿到选择的商品
			2.确认商品的金额(跨微服务)
			3.扣减商品的库存(跨微服务)
			4.创建订单的基本信息
			5.从公务车中删除下单的商品
	*/
	var orderListener orderListener
	/*
		使用rocketMq的基于可靠消息的分布式事务:
		先发送归还库存的half消息,如果本地事务执行失败则commit归还库存的消息,
		如果本地事务执行成功,则回滚归还库存的消息
	*/
	//创建事务生产者
	p, err := rocketmq.NewTransactionProducer(&orderListener, producer.WithNameServer([]string{"192.168.58.130:9876"}))
	if err != nil {
		zap.S().Errorf("生成producer失败:%s", err.Error())
		return nil, err
	}
	if err := p.Start(); err != nil {
		zap.S().Errorf("producer开启失败:%s", err.Error())
		return nil, err
	}
	//将order有关的信息通过json传给ExecuteLocalTransaction
	orderModel := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	jsonString, _ := json.Marshal(orderModel)
	_, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		zap.S().Errorf("发送事务消息失败:%s", err.Error())
		return nil, status.Error(codes.Internal, "发送消息失败")
	}
	if orderListener.Code != codes.OK {
		//当扣减失败时才需要commit归还库存的消息
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: orderModel.OrderSn, Total: orderListener.OrderAmount}, nil
}

//获取订单列表
func (o *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orderModel []model.OrderInfo
	var orderListResponse proto.OrderListResponse
	//此接口给管理后台和查询用户订单列表接口同时使用
	//如果不传用户ID则认为给管理后台使用
	//当UserId为零值时Where语句不会执行
	var total int64
	res := global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	if res.Error != nil {
		return nil, res.Error
	}
	orderListResponse.Total = int32(total)
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orderModel)
	var orderInfoResponse []*proto.OrderInfoResponse
	for _, orderInfo := range orderModel {
		orderInfoResponse = append(orderInfoResponse, &proto.OrderInfoResponse{
			Id:      orderInfo.ID,
			UserId:  orderInfo.User,
			OrderSn: orderInfo.OrderSn,
			PayType: orderInfo.PayType,
			Status:  orderInfo.Status,
			Total:   orderInfo.OrderMount,
			Post:    orderInfo.Post,
			Address: orderInfo.Address,
			Name:    orderInfo.SignerName,
			Mobile:  orderInfo.SingerMobile,
			AddTime: orderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	orderListResponse.Data = orderInfoResponse
	return &orderListResponse, nil
}

//获取订单的详情
func (o *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailRes, error) {
	//先查询订单相关的信息
	var orderModel model.OrderInfo
	var goodsModel []model.OrderGoods
	var orderInfoDetailRes proto.OrderInfoDetailRes
	if res := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&orderModel); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有相关订单记录")
	}
	orderInfo := proto.OrderInfoResponse{
		Id:      orderModel.ID,
		UserId:  orderModel.User,
		OrderSn: orderModel.OrderSn,
		PayType: orderModel.PayType,
		Status:  orderModel.Status,
		Post:    orderModel.Post,
		Total:   orderModel.OrderMount,
		Address: orderModel.Address,
		Name:    orderModel.SignerName,
		Mobile:  orderModel.SingerMobile,
	}
	orderInfoDetailRes.OrderInfo = &orderInfo
	if rsp := global.DB.Where(&model.OrderGoods{Order: orderInfo.Id}).Find(&goodsModel); rsp.Error != nil {
		return nil, rsp.Error
	}
	for _, goodsInfo := range goodsModel {
		orderInfoDetailRes.Goods = append(orderInfoDetailRes.Goods, &proto.OrderItemResponse{
			Id:         goodsInfo.ID,
			OrderId:    goodsInfo.Order,
			GoodsId:    goodsInfo.Goods,
			GoodsName:  goodsInfo.GoodsName,
			GoodsImage: goodsInfo.GoodsImage,
			GoodsPrice: goodsInfo.GoodsPrice,
		})
	}
	return &orderInfoDetailRes, nil
}

//更新订单状态
func (o *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*empty.Empty, error) {
	if res := global.DB.Model(&model.OrderInfo{}).Where("order_sn=?", req.OrderSn).Update("status", req.Status); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &empty.Empty{}, nil
}

//订单超时归还（基于rocketMq的延时消息）
func OrderTimeOut(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)
		fmt.Printf("获取订单超时时间：%v\n", time.Now())
		//查询订单支付状态，如果没有支付，归还库存
		var orderModel model.OrderInfo
		if res := global.DB.Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderModel); res.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if orderModel.Status != "TRADE_SUCCESS" {
			//修改订单状态为已关闭
			tx := global.DB.Begin()
			tx.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).Update("status", "TRADE_CLOSE")
			//归还库存（直接往rocketMq发送一个普通的order_reback消息）
			p, err := NewProducer()
			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				zap.S().Errorf("发送失败:%v\n", err.Error())
				return consumer.ConsumeRetryLater, nil
			}
			tx.Commit()
		}
	}
	return consumer.ConsumeSuccess, nil
}
