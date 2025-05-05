package main

import (
	"go.uber.org/zap"
	"time"
)

// 自定义生产环境 Logger 配置
func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log", // 输出日志到当前目录下的 myproject.log 文件
	}
	return cfg.Build()
}

func main() {
	// 初始化 logger
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 获取 SugarLogger（提供更简洁的格式化输出）
	su := logger.Sugar()
	defer su.Sync()

	url := "https://imooc.com"
	su.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
