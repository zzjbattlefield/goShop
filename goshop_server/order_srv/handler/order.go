package handler

import (
	"context"
	"goshop/order_srv/global"
	"goshop/order_srv/model"
	"goshop/order_srv/proto"

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
	var shopCart model.ShoppingCart
	if res := global.DB.Where("id=?", req.Id).Delete(&shopCart); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &empty.Empty{}, nil
}

//订单相关的功能

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
	var shopCart []model.ShoppingCart
	var goodIds []int32
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选择结算商品")
	}
	goodsNumMap := make(map[int32]int32)
	for _, shopCartInfo := range shopCart {
		goodIds = append(goodIds, shopCartInfo.Goods)
		goodsNumMap[shopCartInfo.Goods] = shopCartInfo.Nums
	}
	//跨服务调用（商品查询）
	rsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodIds})
	if err != nil {
		zap.S().Errorf("【order】-【CreateOrder】-批量查询商品信息失败:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
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
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInvInfo: goodsInvInfo}); err != nil {
		zap.S().Errorf("【order】-【CreateOrder】-批量扣减失败:%s", err.Error())
		return nil, status.Errorf(codes.ResourceExhausted, "批量扣减失败")
	}
	//创建订单的基本信息
	tx := global.DB.Begin()
	orderModel := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	if res := tx.Save(&orderModel); res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.ResourceExhausted, "保存订单失败")
	}

	for _, orderGood := range orderGoods {
		orderGood.Order = orderModel.ID
	}

	//批量插入
	if res := tx.CreateInBatches(&orderGoods, 100); res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.ResourceExhausted, "批量插入失败")
	}

	//删除下单了的购物车记录
	if res := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.ResourceExhausted, "删除购物车记录失败")
	}
	tx.Commit()
	return &proto.OrderInfoResponse{Id: orderModel.ID, OrderSn: orderModel.OrderSn, Total: orderModel.OrderMount}, nil
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
