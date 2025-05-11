/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 15:50:09
 * @LastEditors: Will zw37520@gmail.com
 * @LastEditTime: 2025-05-11 13:13:10
 * @FilePath: /joyshop_api/user-web/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServeConfig struct {
	Name        string          `mapstructure:"name" json:"name"`
	Host        string          `mapstructure:"host" json:"host"`
	Port        int             `mapstructure:"port" json:"port"`
	Lang        string          `mapstructure:"lang" json:"lang"`
	UserSrvInfo UserSrvConfig   `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JwtConfig       `mapstructure:"jwt" json:"jwt"`
	AliyunSms   AliyunSmsConfig `mapstructure:"aliyun_sms" json:"aliyun_sms"`
	RedisInfo   RedisConfig     `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig    `mapstructure:"consul" json:"consul"`
	NacosInfo   NacosConfig     `mapstructure:"nacos" json:"nacos"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"signing_key" json:"signing_key"`
	ExpireTime int    `mapstructure:"expire_time" json:"expire_time"` // token过期时间（小时）
}

// 阿里云短信配置
type AliyunSmsConfig struct {
	SignName      string `mapstructure:"sign_name" json:"sign_name"`
	TemplateCode  string `mapstructure:"template_code" json:"template_code"`
	PhoneNumbers  string `mapstructure:"phone_numbers" json:"phone_numbers"`
	TemplateParam string `mapstructure:"template_param" json:"template_param"`
	AccessKeyId   string `mapstructure:"access_key_id" json:"access_key_id"`
	AccessSecret  string `mapstructure:"access_secret" json:"access_secret"`
}

// redis配置
type RedisConfig struct {
	Host       string `mapstructure:"host" json:"host"`
	Port       int    `mapstructure:"port" json:"port"`
	ExpireTime int    `mapstructure:"expire_time" json:"expire_time"` // 验证码过期时间（分钟）
}

// Consul配置
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

// NacosConfig Nacos配置
type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	Timeout   uint64 `mapstructure:"timeout" json:"timeout"`
	LogDir    string `mapstructure:"logDir" json:"logDir"`
	CacheDir  string `mapstructure:"cacheDir" json:"cacheDir"`
	LogLevel  string `mapstructure:"logLevel" json:"logLevel"`
	DataId    string `mapstructure:"dataId" json:"dataId"`
	Group     string `mapstructure:"group" json:"group"`
}
