package main

import (
	"github.com/nkmr-jp/go-zap-scaffold/logger"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()   // flush log buffer
	logger.SyncWhenStop() // flush log buffer. when Interrupt or kill.

	logger.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
}
