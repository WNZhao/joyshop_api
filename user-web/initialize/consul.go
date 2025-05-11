/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-09 16:14:19
 * @LastEditors: Will zw37520@gmail.com
 * @LastEditTime: 2025-05-10 16:07:18
 * @FilePath: /joyshop_api/user-web/initialize/consul.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"joyshop_api/user-web/global"
	"joyshop_api/user-web/utils"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// InitConsulRegister 初始化 Consul 服务注册
func InitConsulRegister() error {
	// 创建 Consul 客户端配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.7:8500" // 修改为正确的 Consul 地址

	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("创建 Consul 客户端失败: %s", err.Error())
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	zap.S().Infof("服务名称: %s", global.ServerConfig.Name)
	zap.S().Infof("服务端口: %d", global.ServerConfig.Port)

	// 获取本机IP地址
	localIP, err := utils.GetLocalIP()
	if err != nil {
		zap.S().Errorw("获取本机IP地址失败", "msg", err.Error())
		return err
	}
	zap.S().Infof("本机IP地址: %s", localIP)

	registration.ID = fmt.Sprintf("%s-%d", global.ServerConfig.Name, global.ServerConfig.Port)
	registration.Port = global.ServerConfig.Port
	registration.Tags = []string{"user-web", "web", "api"}
	registration.Address = localIP // 使用自动获取的IP地址

	// 生成对应的检查对象
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/health", registration.Address, registration.Port)
	check.Timeout = "10s"
	check.Interval = "15s"
	check.DeregisterCriticalServiceAfter = "30s"
	registration.Check = check

	// 打印健康检查URL
	zap.S().Infof("健康检查URL: %s", check.HTTP)

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorw("注册服务到 Consul 失败", "msg", err.Error())
		return err
	}

	zap.S().Infof("服务 [%s] 已成功注册到 Consul, 地址: %s:%d", registration.Name, registration.Address, registration.Port)
	return nil
}
