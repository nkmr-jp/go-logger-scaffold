package main

import (
	"github.com/nkmr-jp/go-zap-scaffold/logger"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
	logger.Sync()
}
