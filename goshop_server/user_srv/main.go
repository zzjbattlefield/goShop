package main

import (
	"flag"
	"fmt"
	"goshop/user_srv/global"
	"goshop/user_srv/handler"
	"goshop/user_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"goshop/user_srv/initialize"
	proto "goshop/user_srv/proto"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "ip地址")
	port := flag.Int("port", 0, "端口号")
	flag.Parse()
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)
	if *port == 0 {
		//没有指定端口号时自动生成
		*port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip", *ip)
	zap.S().Info("port", *port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	check := new(api.AgentServiceCheck)
	check.GRPC = fmt.Sprintf("%s:%d", global.ServerConfig.Host, *port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "15s"
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceID,
		Port:    *port,
		Tags:    []string{"imooc", "user-srv"},
		Address: global.ServerConfig.Host,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	go func() {
		//server会阻塞
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
