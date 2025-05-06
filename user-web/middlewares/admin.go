package middlewares

import (
	"github.com/gin-gonic/gin"
	"joyshop_api/user-web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 2 {
			c.JSON(200, gin.H{
				"code": 403,
				"msg":  "没有权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
