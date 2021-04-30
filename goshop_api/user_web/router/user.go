package router

import (
	"goshop_api/user_web/api"
	"goshop_api/user_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserGroup := Router.Group("user")
	{
		UserGroup.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserGroup.POST("pwd_login", api.PasswordLogin)
		UserGroup.POST("register", api.Register)
	}

}
