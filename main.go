package main

import (
	"github.com/nkmr-jp/go-logger-scaffold/logger"
	"go.uber.org/zap"
)

func main() {
	logger.SetLogFile("./log/hoge/app_%Y-%m-%d.log")
	logger.InitLogger()
	defer logger.Sync()   // flush log buffer
	logger.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// time.Sleep(60*time.Second)
	logger.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
}
