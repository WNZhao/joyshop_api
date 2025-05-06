/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-06 10:13:13
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-06 10:33:29
 * @FilePath: /joyshop_api/user-web/utils/validator.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package utils

import (
	"joyshop_api/user-web/global"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// RemoveFormName 移除表单验证错误信息中的表单名称前缀
func RemoveFormName(errors validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		// 获取字段名，去掉表单名称前缀
		field := err.Field()
		if idx := strings.Index(field, "."); idx != -1 {
			field = field[idx+1:]
		}
		// 将字段名转换为小写
		field = strings.ToLower(field)
		errorMap[field] = err.Translate(global.Trans)
	}
	return errorMap
}

// HandleValidatorError 处理表单验证错误处理
func HandleValidatorError(ctx *gin.Context, err error, formName string) bool {
	if err != nil {
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.S().Errorw("[%s] 参数校验失败", formName, "msg", err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "参数错误:" + err.Error(),
			})
			return true
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  RemoveFormName(errors),
		})
		return true
	}
	return false
}
