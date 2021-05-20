package initialize

import (
	"goshop_api/userop_web/middlewares"
	apiRouter "goshop_api/userop_web/router"
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
	ApiGroup := router.Group("/up/v1")
	apiRouter.InitAddressRouter(ApiGroup)
	apiRouter.InitMessageRouter(ApiGroup)
	apiRouter.InitUserFavRouter(ApiGroup)
	return router
}
