/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 15:50:09
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-07 15:49:53
 * @FilePath: /joyshop_api/user-web/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServeConfig struct {
	Name        string          `mapstructure:"name"`
	Port        int             `mapstructure:"port"`
	Lang        string          `mapstructure:"lang"`
	UserSrvInfo UserSrvConfig   `mapstructure:"user_srv"`
	JWTInfo     JwtConfig       `mapstructure:"jwt"`
	AliyunSms   AliyunSmsConfig `mapstructure:"aliyun_sms"`
	RedisInfo   RedisConfig     `mapstructure:"redis"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"signing_key"`
	ExpireTime int    `mapstructure:"expire_time"` // token过期时间（小时）
}

// 阿里云短信配置
type AliyunSmsConfig struct {
	SignName      string `mapstructure:"sign_name"`
	TemplateCode  string `mapstructure:"template_code"`
	PhoneNumbers  string `mapstructure:"phone_numbers"`
	TemplateParam string `mapstructure:"template_param"`
	AccessKeyId   string `mapstructure:"access_key_id"`
	AccessSecret  string `mapstructure:"access_secret"`
}

// redis配置
type RedisConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	ExpireTime int    `mapstructure:"expire_time"` // 验证码过期时间（分钟）
}
