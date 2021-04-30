package initialize

import (
	"fmt"
	"goshop/order_srv/global"
	"goshop/order_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//初始化第三方微服务连接

func InitSrvConn() {
	conslutInfo := global.ServerConfig.ConsulInfo
	connect, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", conslutInfo.Host, conslutInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		zap.S().Errorf("[InitSrvConn] 连接 [商品服务失败]", "msg", err.Error())
	}
	global.GoodsSrvClient = proto.NewGoodsClient(connect)
	//初始化库存服务

	invConnect, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", conslutInfo.Host, conslutInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		zap.S().Errorf("[InitSrvConn] 连接 [库存服务失败]", "msg", err.Error())
	}
	global.InventorySrvClient = proto.NewInventoryClient(invConnect)
}
