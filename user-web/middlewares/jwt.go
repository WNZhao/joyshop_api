/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-06 14:03:12
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-06 16:02:59
 * @FilePath: /joyshop_api/user-web/middlewares/jwt.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package middlewares

import (
	"joyshop_api/user-web/global"
	"joyshop_api/user-web/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT结构体
type JWT struct {
	SigningKey []byte
}

// NewJWT 创建JWT对象
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				zap.S().Errorw("[JWT] 不合法的token", "msg", err.Error())
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				zap.S().Errorw("[JWT] token过期", "msg", err.Error())
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				zap.S().Errorw("[JWT] token未生效", "msg", err.Error())
				return nil, err
			}
		}
		zap.S().Errorw("[JWT] 解析token出错", "msg", err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	zap.S().Errorw("[JWT] token无效", "msg", err.Error())
	return nil, err
}

// JWTAuth 登录验证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("x-token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请登录",
			})
			c.Abort()
			return
		}
		jwtObj := NewJWT()
		claims, err := jwtObj.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "无效的token"})
			c.Abort()
			return
		}
		// 校验token合法性
		c.Set("user_id", claims.ID)
		c.Set("user_nick", claims.NickName)
		c.Set("role_id", claims.AuthorityId)
		c.Set("claims", claims) // ★★★
		c.Next()
	}
}
