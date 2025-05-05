package router

import (
	"joyshop_api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("create", api.CreateUser)
		UserRouter.PUT("update", api.UpdateUser)
		UserRouter.DELETE("delete", api.DeleteUser)
		UserRouter.POST("password_login", api.PassWordLogin)
	}
}
