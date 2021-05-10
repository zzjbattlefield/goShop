package router

import (
	"goshop_api/order_web/api/shop_cart"
	"goshop_api/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	shopCartGroup := Router.Group("shop_cart").Use(middlewares.JWTAuth())
	{
		shopCartGroup.GET("", shop_cart.List)          //获取购物车列表
		shopCartGroup.DELETE("/:id", shop_cart.Delete) //删除购物车条目
		shopCartGroup.POST("", shop_cart.New)          //获取购物车列表
		shopCartGroup.PATCH("/:id", shop_cart.Update)  //删除购物车条目
	}
}
