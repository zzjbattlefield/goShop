package router

import (
	"goshop_api/goods_web/api/category"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys")
	{
		CategoryRouter.GET("", category.List)
		CategoryRouter.GET("/:id", category.Detail)
		CategoryRouter.POST("", category.New)
		CategoryRouter.PUT("/:id", category.Update)
		CategoryRouter.DELETE("/:id", category.Delete)
	}
}
