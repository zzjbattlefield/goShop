package main

import (
	"fmt"
	"goshop_api/userop_web/global"
	"goshop_api/userop_web/initialize"
	"goshop_api/userop_web/utils"
	"goshop_api/userop_web/utils/register/consul"
	myvalidator "goshop_api/userop_web/validator"
	"os"
	"os/signal"
	"syscall"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
)

func main() {
	//初始化日志
	initialize.InitLogger()
	_ = initialize.InitTrans("zh")
	initialize.InitConfig()
	//初始化srv连接
	initialize.InitSrvConn()
	//初始化路由
	router := initialize.Routers()
	//如果是本地开发环境端口固定 线上环境自动获取端口号
	if global.IsDebug != 1 {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServeConfig.Port = port
		}
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	//服务注册
	register_client := consul.NewRegistryClient(global.ServeConfig.ConsulInfo.Host, global.ServeConfig.ConsulInfo.Port)
	serviceId, _ := uuid.NewV4()
	serviceIdStr := fmt.Sprintf("%s", serviceId)
	err := register_client.Register(global.ServeConfig.Host, global.ServeConfig.Port, global.ServeConfig.Name, global.ServeConfig.Tags, serviceIdStr)
	if err != nil {
		zap.S().Panic("注册服务失败", err.Error())
	}
	zap.S().Infof("启动服务,端口:%d", global.ServeConfig.Port)

	go func() {
		if err := router.Run(fmt.Sprintf(":%d", global.ServeConfig.Port)); err != nil {
			zap.S().Panic("启动失败", err.Error())
		}
	}()
	//监听终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = register_client.DeRegister(serviceIdStr)
	if err != nil {
		zap.S().Panic("注销失败", err)
	} else {
		zap.S().Info("注销成功")
	}
}
