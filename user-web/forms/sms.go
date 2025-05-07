package forms

// SendSmsForm 发送短信表单
type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type   int    `form:"type" json:"type" binding:"required,oneof=1 2"` // 1: 注册 2: 登录
}
