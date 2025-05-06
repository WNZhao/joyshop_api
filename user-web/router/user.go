package router

import (
	"github.com/gin-gonic/gin"
	"joyshop_api/user-web/api"
	"joyshop_api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user") // .Use(middlewares.JWTAuth())
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("create", api.CreateUser)
		UserRouter.PUT("update", api.UpdateUser)
		UserRouter.DELETE("delete", api.DeleteUser)
		UserRouter.POST("password_login", api.PassWordLogin)
	}
}
