package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

// 图片验证码
var storage = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, storage)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("验证码生成失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "验证码生成失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码生成成功",
		"data": b64s,
		"id":   id,
	})

}
