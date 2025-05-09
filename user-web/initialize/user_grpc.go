/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-09 16:31:26
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-09 16:37:37
 * @FilePath: /joyshop_api/user-web/initialize/user_grpc.go
 * @Description: 用户服务 gRPC 客户端初始化
 */
package initialize

import (
	"fmt"
	"joyshop_api/user-web/global"
	"joyshop_api/user-web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitUserGrpcClient 初始化用户服务 gRPC 客户端
func InitUserGrpcClient() error {
	// 从consul中获取服务信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	consulClient, err := api.NewClient(cfg)

	if err != nil {
		zap.S().Errorw("[InitUserGrpcClient] 创建consul客户端失败", "msg", err.Error())
		return err
	}

	// 使用服务名称进行服务发现
	serviceName := global.ServerConfig.UserSrvInfo.Name
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)

	if err != nil {
		zap.S().Errorw("[InitUserGrpcClient] 获取服务列表失败", "msg", err.Error())
		return err
	}

	if len(services) == 0 {
		zap.S().Errorw("[InitUserGrpcClient] 未找到可用服务", "service", serviceName)
		return fmt.Errorf("未找到可用服务: %s", serviceName)
	}

	// 获取第一个健康的服务实例
	service := services[0].Service
	userSrvHost := service.Address
	userSrvPort := service.Port

	// 建立grpc连接
	userConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", userSrvHost, userSrvPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.S().Errorw("[InitUserGrpcClient] 连接用户服务失败", "msg", err.Error())
		return err
	}

	// 生成grpc客户端
	global.UserConn = userConn
	global.UserClient = proto.NewUserClient(userConn)
	zap.S().Infof("成功连接到用户服务: %s:%d", userSrvHost, userSrvPort)
	return nil
}
