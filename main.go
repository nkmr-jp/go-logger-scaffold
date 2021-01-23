package main

import (
	"fmt"
	"time"

	"github.com/nkmr-jp/go-logger-scaffold/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// set values from cli.
// ex.`go run -ldflags "-X main.version=v1.0.0 -X main.srcRootDir=$PWD" main.go`

// version git revision or tag. set from go cli.
var version string

// srcRootDir set from cli.
var srcRootDir string

const (
	consoleField = "console"
	traceIDField = "trace_id"
	urlFormat    = "https://github.com/nkmr-jp/go-logger-scaffold/blob/%s"
)

func main() {
	// Set options
	logger.SetLogFile("./log/app_%Y-%m-%d.log")
	logger.SetVersion(version)
	logger.SetRepositoryCallerEncoder(urlFormat, version, srcRootDir)
	logger.SetConsoleField(consoleField, traceIDField)
	logger.SetLogLevel(zapcore.InfoLevel)
	logger.SetOutputType(logger.OutputTypeShortConsoleAndFile)

	// Initialize
	logger.InitLogger()
	defer logger.Sync()   // flush log buffer
	logger.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Run examples
	examples()
	newWrapperExample()
}

func examples() {
	// basic
	logger.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
	// error log
	err := fmt.Errorf("error message")
	logger.Errorf("SOME_ERROR", err)
	// debug log
	logger.Debug("DEBUG_MESSAGE")
	logger.Println(err)
	// warn log
	logger.Warn("WARN_MESSAGE")
	// display to console log
	logger.Info("DISPLAY_TO_CONSOLE", zap.String(consoleField, "display to console"))
}

func newWrapperExample() {
	// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
	w := logger.NewWrapper(
		zap.Int("user_id", 1),
		zap.Int64(traceIDField, time.Now().UnixNano()),
	)
	w.Info("CONTEXT_SCOPE_INFO")
	w.Errorf("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
}
