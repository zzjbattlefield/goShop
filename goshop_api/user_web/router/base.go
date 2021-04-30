package router

import (
	"goshop_api/user_web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptCha)
		BaseRouter.POST("send_sms", api.SendSms)
	}
}
