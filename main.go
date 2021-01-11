package main

import (
	"fmt"

	"github.com/nkmr-jp/go-logger-scaffold/logger"
	"go.uber.org/zap"
)

func main() {
	logger.SetLogFile("./log/app_%Y-%m-%d.log")
	logger.InitLogger()
	defer logger.Sync()   // flush log buffer
	logger.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// time.Sleep(60*time.Second)

	// example
	logger.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))

	// console log example
	logger.Info("OUT_PUT_TO_CONSOLE", logger.ConsoleField("messages to be displayed on the console"))

	// error log example
	var err error
	err = fmt.Errorf("error message")
	logger.Errorf("SOME_ERROR", err)
}
