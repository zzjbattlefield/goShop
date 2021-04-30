package router

import (
	"goshop_api/goods_web/api/goods"
	"goshop_api/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsGroup := Router.Group("goods")
	{
		GoodsGroup.GET("", goods.List)
		GoodsGroup.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsGroup.GET("/:id", goods.Detail)                                                      //获取商品详情
		GoodsGroup.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) //删除商品
		GoodsGroup.GET("/:id/stocks", goods.Stocks)
		GoodsGroup.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsGroup.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}

}
