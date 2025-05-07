/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 21:34:36
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-07 13:43:37
 * @FilePath: /joyshop_api/user-web/forms/user.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package forms

type PassWordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码 mobile是一个自定义的验证器
	Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` // 验证码
	CaptchaId string `form:"captchaId" json:"captchaId" binding:"required"`
}

// 用户注册
type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码 mobile是一个自定义的验证器
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	//Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` // 图片验证码
	Code string `form:"code" json:"code" binding:"required,min=6,max=6"` // 短信验证码
}
