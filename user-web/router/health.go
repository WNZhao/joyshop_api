package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitHealthRouter 初始化健康检查路由
func InitHealthRouter(Router *gin.Engine) {
	// 健康检查路由
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
}
