/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-09 16:14:19
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-09 16:47:47
 * @FilePath: /joyshop_api/user-web/initialize/consul.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"joyshop_api/user-web/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// InitConsulRegister 初始化 Consul 服务注册
func InitConsulRegister() error {
	// 创建 Consul 客户端配置
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("创建 Consul 客户端失败: %s", err.Error())
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = fmt.Sprintf("%s-%d", global.ServerConfig.Name, global.ServerConfig.Port)
	registration.Port = global.ServerConfig.Port
	registration.Tags = []string{"user", "web", "api"}
	registration.Address = global.ServerConfig.Host // 使用配置中的主机地址

	// 生成对应的检查对象
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/health", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "15s"
	check.DeregisterCriticalServiceAfter = "15s"
	registration.Check = check

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("注册服务到 Consul 失败: %s", err.Error())
	}

	zap.S().Infof("服务 [%s] 已成功注册到 Consul, 地址: %s:%d", registration.Name, registration.Address, registration.Port)
	return nil
}
