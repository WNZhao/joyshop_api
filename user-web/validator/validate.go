package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// 自定义表单验证

func ValidateMobile(f1 validator.FieldLevel) bool {
	mobile := f1.Field().String()
	// 使用正则表达式验证手机号
	// 这里使用的是中国大陆的手机号格式
	reg := `^1[3-9]\d{9}$`
	regexp := regexp.MustCompile(reg)
	return regexp.MatchString(mobile)
}
