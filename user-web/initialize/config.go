/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 15:52:24
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-05 16:18:11
 * @FilePath: /joyshop_api/user-web/initialize/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"github.com/fsnotify/fsnotify"
	"joyshop_api/user-web/global"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	// 设置配置文件名称
	configName := "config-debug.yaml"
	if os.Getenv("APP_ENV") == "production" {
		configName = "config-prod.yaml"
	}

	// 设置配置文件路径
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./user-web")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Panicf("读取配置文件失败: %v", err)
	}

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 解析配置到结构体
	if err := viper.Unmarshal(&global.ServerConfig); err != nil {
		zap.S().Panicf("解析配置文件失败: %v", err)
	}
	zap.S().Infof("配置信息: %v", global.ServerConfig)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", in.Name)
		if err := viper.Unmarshal(&global.ServerConfig); err != nil {
			zap.S().Panicf("重新解析配置文件失败: %v", err)
		}
		zap.S().Infof("重新加载配置信息: %v", global.ServerConfig)
	})
}

func GetEnvInfo(env string) bool {
	return viper.GetBool(env)
}
