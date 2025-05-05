package initialize

import (
	"go.uber.org/zap"
	"os"
)

func InitLogger() {
	// 根据环境变量设置日志模式
	var logger *zap.Logger
	if os.Getenv("APP_ENV") == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	zap.ReplaceGlobals(logger)
}
