package router

import (
	"goshop_api/goods_web/api/banner"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	bannerGroup := Router.Group("banners")
	{
		bannerGroup.GET("", banner.List)
		bannerGroup.POST("", banner.New)
		bannerGroup.DELETE("/:id", banner.Delete)
		bannerGroup.PUT("/:id", banner.Update)
	}
}
