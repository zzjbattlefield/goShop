package initialize

import (
	"fmt"
	"goshop_api/user_web/global"
	proto "goshop_api/user_web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	//从注册中心获取用户服务信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServeConfig.ConsulInfo.Host, global.ServeConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, global.ServeConfig.UserSrvInfo.Name))
	// data, err := client.Agent().ServicesWithFilter(`Service == "user-srv"`)
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接用户服务失败")
		return
	}

	//拨号连接用户grpc服务器
	connect, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorf("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
	}
	userServerClient := proto.NewUserClient(connect)
	global.UserClinet = userServerClient
}
