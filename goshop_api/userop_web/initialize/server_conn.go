package initialize

import (
	"fmt"
	"goshop_api/userop_web/global"
	"goshop_api/userop_web/proto"

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

	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServeConfig.ConsulInfo.Host, global.ServeConfig.ConsulInfo.Port, global.ServeConfig.UserOpSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorf("[InitSrvConn] 连接 [用户操作服务失败]", "msg", err.Error())
	}
	global.UserFavClient = proto.NewUserFavClient(userOpConn)
	global.MessageClient = proto.NewMessageClient(userOpConn)
	global.AddressClient = proto.NewAddressClient(userOpConn)
}
