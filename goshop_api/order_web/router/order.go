package router

import (
	"goshop_api/order_web/api/order"
	"goshop_api/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderGroup := Router.Group("orders").Use(middlewares.JWTAuth())
	{
		OrderGroup.GET("", order.List) //获取订单列表
		OrderGroup.POST("", order.New) //新建订单
		OrderGroup.GET("/:id", order.Detail)
	}

}
