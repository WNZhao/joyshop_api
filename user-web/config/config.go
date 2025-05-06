/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 15:50:09
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-06 14:10:09
 * @FilePath: /joyshop_api/user-web/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServeConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	Lang        string        `mapstructure:"lang"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo     JwtConfig     `mapstructure:"jwt"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}
