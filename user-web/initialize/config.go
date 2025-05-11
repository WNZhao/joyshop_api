/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 15:52:24
 * @LastEditors: Will zw37520@gmail.com
 * @LastEditTime: 2025-05-11 10:34:59
 * @FilePath: /joyshop_api/user-web/initialize/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"joyshop_api/user-web/config"
	"joyshop_api/user-web/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var NacosConfig config.NacosConfig

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

	// 尝试从 Nacos 读取配置
	if err := initNacosConfig(); err != nil {
		zap.S().Warnf("从 Nacos 读取配置失败，使用本地配置: %v", err)
		if err := viper.Unmarshal(&global.ServerConfig); err != nil {
			zap.S().Panicf("解析配置文件失败: %v", err)
		}
	}

	// 打印详细的配置信息
	zap.S().Infof("服务名称: %s", global.ServerConfig.Name)
	zap.S().Infof("服务端口: %d", global.ServerConfig.Port)

	// JWT配置
	// zap.S().Infof("JWT密钥: %s", global.ServerConfig.JWTInfo.SigningKey)
	zap.S().Infof("JWT过期时间: %d小时", global.ServerConfig.JWTInfo.ExpireTime)

	// 阿里云短信配置
	zap.S().Infof("阿里云短信签名: %s", global.ServerConfig.AliyunSms.SignName)
	zap.S().Infof("阿里云短信模板: %s", global.ServerConfig.AliyunSms.TemplateCode)
	zap.S().Infof("阿里云短信手机号: %s", global.ServerConfig.AliyunSms.PhoneNumbers)

	// Redis配置
	zap.S().Infof("Redis地址: %s", global.ServerConfig.RedisInfo.Host)
	zap.S().Infof("Redis端口: %d", global.ServerConfig.RedisInfo.Port)

	// Consul配置
	zap.S().Infof("Consul地址: %s", global.ServerConfig.ConsulInfo.Host)
	zap.S().Infof("Consul端口: %d", global.ServerConfig.ConsulInfo.Port)

	// Nacos配置
	zap.S().Infof("Nacos地址: %s", global.ServerConfig.NacosInfo.Host)
	zap.S().Infof("Nacos端口: %d", global.ServerConfig.NacosInfo.Port)
	zap.S().Infof("Nacos命名空间: %s", global.ServerConfig.NacosInfo.Namespace)
	zap.S().Infof("Nacos超时时间: %d", global.ServerConfig.NacosInfo.Timeout)
	zap.S().Infof("Nacos日志目录: %s", global.ServerConfig.NacosInfo.LogDir)
	zap.S().Infof("Nacos缓存目录: %s", global.ServerConfig.NacosInfo.CacheDir)
	zap.S().Infof("Nacos日志级别: %s", global.ServerConfig.NacosInfo.LogLevel)
	zap.S().Infof("Nacos数据ID: %s", global.ServerConfig.NacosInfo.DataId)
	zap.S().Infof("Nacos分组: %s", global.ServerConfig.NacosInfo.Group)

	// 监听本地配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", in.Name)
		if err := viper.Unmarshal(&global.ServerConfig); err != nil {
			zap.S().Panicf("重新解析配置文件失败: %v", err)
		}
		zap.S().Infof("重新加载配置信息: %v", global.ServerConfig)
	})
}

func initNacosConfig() error {
	// 设置 Nacos 配置文件名称
	configName := "nacos-dev.yaml"
	if os.Getenv("APP_ENV") == "production" {
		configName = "nacos-prod.yaml"
	}

	// 设置配置文件路径
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		zap.S().Errorf("读取 Nacos 配置文件失败: %v", err)
		return err
	}

	// 解析 Nacos 配置
	if err := v.UnmarshalKey("nacos", &NacosConfig); err != nil {
		zap.S().Errorf("解析 Nacos 配置失败: %v", err)
		return err
	}

	zap.S().Infof("成功读取 Nacos 配置: %+v", NacosConfig)

	// 创建 Nacos 客户端
	sc := []constant.ServerConfig{
		{
			IpAddr: NacosConfig.Host,
			Port:   NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         NacosConfig.Namespace,
		TimeoutMs:           NacosConfig.Timeout,
		NotLoadCacheAtStart: true,
		LogDir:              NacosConfig.LogDir,
		CacheDir:            NacosConfig.CacheDir,
		LogLevel:            NacosConfig.LogLevel,
	}

	zap.S().Infof("正在创建 Nacos 客户端...")
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		zap.S().Errorf("创建 Nacos 客户端失败: %v", err)
		return err
	}

	zap.S().Infof("正在从 Nacos 获取配置...")
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: NacosConfig.DataId,
		Group:  NacosConfig.Group,
	})
	if err != nil {
		zap.S().Errorf("从 Nacos 获取配置失败: %v", err)
		return err
	}

	zap.S().Infof("成功从 Nacos 获取配置: %s", content)

	// 解析配置内容
	if err := yaml.Unmarshal([]byte(content), &global.ServerConfig); err != nil {
		zap.S().Errorf("解析 Nacos 配置内容失败: %v", err)
		return err
	}

	zap.S().Infof("成功解析 Nacos 配置内容: %+v", global.ServerConfig)

	// 监听配置变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: NacosConfig.DataId,
		Group:  NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Infof("Nacos 配置发生变化: namespace=%s, group=%s, dataId=%s", namespace, group, dataId)
			if err := yaml.Unmarshal([]byte(data), &global.ServerConfig); err != nil {
				zap.S().Errorf("重新解析 Nacos 配置失败: %v", err)
				return
			}
			zap.S().Infof("成功更新配置: %+v", global.ServerConfig)
		},
	})
	if err != nil {
		zap.S().Errorf("设置 Nacos 配置监听失败: %v", err)
		return err
	}

	return nil
}

func GetEnvInfo(env string) bool {
	return viper.GetBool(env)
}
