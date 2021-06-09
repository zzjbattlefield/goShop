package main

import (
	"flag"
	"fmt"
	"goshop/inventory_srv/global"
	"goshop/inventory_srv/handler"
	"goshop/inventory_srv/utils"
	"goshop/inventory_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	"goshop/inventory_srv/initialize"
	"goshop/inventory_srv/proto"

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
	port := flag.Int("port", 0, "端口号")
	flag.Parse()
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedis()
	zap.S().Info(global.ServerConfig)
	if *port == 0 {
		//没有指定端口号时自动生成
		*port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip", *ip)
	zap.S().Info("port", *port)
	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.Inventort{})
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
	go func() {
		//server会阻塞
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	//监听库存归还topic
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.58.130:9876"}),
		consumer.WithGroupName("mxshop-inventory"))
	if err != nil {
		panic("新建consumer失败")
	}
	//订阅消息
	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Shutdown()
	if err = register_client.DeRegister(uuid.String()); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
