package initialize

import (
	"fmt"
	"goshop_api/order_web/global"
	"goshop_api/order_web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	//拨号连接商品grpc服务器
	connect, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServeConfig.ConsulInfo.Host, global.ServeConfig.ConsulInfo.Port, global.ServeConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorf("[InitSrvConn] 连接 [商品服务失败]", "msg", err.Error())
	}
	global.GoodsClient = proto.NewGoodsClient(connect)

	//连接订单服务
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServeConfig.ConsulInfo.Host, global.ServeConfig.ConsulInfo.Port, global.ServeConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务失败】", "msg", err.Error())
	}
	global.OrderClient = proto.NewOrderClient(orderConn)

	invConnect, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServeConfig.ConsulInfo.Host, global.ServeConfig.ConsulInfo.Port, global.ServeConfig.InventoryInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorf("[InitSrvConn] 连接 [库存服务失败]", "msg", err.Error())
	}
	global.InventoryClient = proto.NewInventoryClient(invConnect)
}
