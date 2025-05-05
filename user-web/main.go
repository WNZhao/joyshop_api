/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-04 11:40:59
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-05 16:31:09
 * @FilePath: /joyshop_api/user-web/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"joyshop_api/user-web/global"
	"joyshop_api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	// 1.初始化日志
	initialize.InitLogger()

	// 2.初始化配置
	initialize.InitConfig()

	// 3.初始化路由
	Router := initialize.Routers()

	// 4.启动服务
	port := fmt.Sprintf(":%d", global.ServerConfig.Port)
	zap.S().Infof("服务启动 端口: %s", port)

	if err := Router.Run(port); err != nil {
		zap.S().Panic("启动失败", err)
	}
}
