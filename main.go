package main

import (
	"fmt"
	"time"

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

	// NewWrapper example.
	// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
	w := logger.NewWrapper(
		zap.Int("user_id", 1),
		zap.Int64("trace_id", time.Now().UnixNano()),
	)
	w.Info("CONTEXT_SCOPE_INFO")
	w.Errorf("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
}
