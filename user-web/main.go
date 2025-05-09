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
	myvalidate "joyshop_api/user-web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
)

func main() {
	// 1.初始化日志
	initialize.InitLogger()

	// 2.初始化配置
	initialize.InitConfig()

	// 3.初始化路由
	Router := initialize.Routers()

	// 4.初始化翻译器
	if err := initialize.InitTrans(global.ServerConfig.Lang); err != nil {
		zap.S().Panicf("翻译器初始化失败: %v", err)
	}

	// 5.注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidate.ValidateMobile)
		// 自定义翻译器
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}不是一个有效的手机号", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 6.初始化 gRPC 客户端
	if err := initialize.InitUserGrpcClient(); err != nil {
		zap.S().Panic("初始化 gRPC 客户端失败", err.Error())
	}
	defer global.UserConn.Close()

	// 7.注册服务到 Consul
	if err := initialize.InitConsulRegister(); err != nil {
		zap.S().Panicf("服务注册失败: %v", err)
	}

	// 8.启动服务
	port := fmt.Sprintf(":%d", global.ServerConfig.Port)
	zap.S().Infof("服务启动 端口: %s", port)

	if err := Router.Run(port); err != nil {
		zap.S().Panic("启动失败", err)
	}
}
