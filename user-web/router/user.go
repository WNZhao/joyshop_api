/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-04 12:57:25
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-07 18:03:45
 * @FilePath: /joyshop_api/user-web/router/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package router

import (
	"joyshop_api/user-web/api"
	"joyshop_api/user-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user") // .Use(middlewares.JWTAuth())
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("create", api.CreateUser)
		UserRouter.PUT("update", api.UpdateUser)
		UserRouter.DELETE("delete", api.DeleteUser)
		UserRouter.POST("password_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}
}
