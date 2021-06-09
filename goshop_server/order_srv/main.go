package main

import (
	"flag"
	"fmt"
	"goshop/order_srv/global"
	"goshop/order_srv/handler"
	"goshop/order_srv/utils"
	"goshop/order_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	"goshop/order_srv/initialize"
	"goshop/order_srv/proto"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "ip地址")
	port := flag.Int("port", 50051, "端口号")
	flag.Parse()
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvConn()
	initialize.InitRedis()
	zap.S().Info(global.ServerConfig)
	if *port == 0 {
		//没有指定端口号时自动生成
		*port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip", *ip)
	zap.S().Info("port", *port)
	server := grpc.NewServer()
	proto.RegisterOrderServer(server, &handler.OrderServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	uuid := uuid.NewV4()
	err = register_client.Register(global.ServerConfig.Host, *port, global.ServerConfig.Name, global.ServerConfig.Tags, uuid.String())
	//服务注册
	if err != nil {
		panic(err)
	}
	//监听订单topic
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.58.130:9876"}),
		consumer.WithGroupName("mxshop-order"))
	if err != nil {
		panic("新建consumer失败")
	}
	//订阅消息
	if err := c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeOut); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
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
	_ = c.Shutdown()
	if err = register_client.DeRegister(uuid.String()); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
