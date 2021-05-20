package router

import (
	"goshop_api/userop_web/api/message"
	"goshop_api/userop_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message")
	{
		MessageRouter.GET("", middlewares.JWTAuth(), message.List) // 留言列表页
		MessageRouter.POST("", middlewares.JWTAuth(), message.New) //新建留言
	}
}
