package initialize

import (
	"goshop_api/goods_web/middlewares"
	apiRouter "goshop_api/goods_web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	router := gin.Default()
	//配置健康检查的路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	//配置跨域
	router.Use(middlewares.Cors())
	ApiGroup := router.Group("/g/v1")
	apiRouter.InitGoodsRouter(ApiGroup)
	apiRouter.InitCategoryRouter(ApiGroup)
	apiRouter.InitBannerRouter(ApiGroup)
	return router
}
