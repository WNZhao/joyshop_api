package router

import (
	"github.com/gin-gonic/gin"
	"joyshop_api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("create", api.CreateUser)
		UserRouter.PUT("update", api.UpdateUser)
		UserRouter.DELETE("delete", api.DeleteUser)
	}

}
