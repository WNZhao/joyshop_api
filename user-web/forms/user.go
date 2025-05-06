package forms

type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码 mobile是一个自定义的验证器
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}
