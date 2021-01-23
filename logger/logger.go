// Created from https://github.com/nkmr-jp/go-logger-scaffold
package logger

import (
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultVersion = "1.0.0"
)

type ConsoleType int

const (
	ConsoleTypeAll ConsoleType = iota
	ConsoleTypeError
	ConsoleTypeNone
)

var (
	once        sync.Once
	zapLogger   *zap.Logger
	consoleType ConsoleType
	version     string
)

// Initialize the Logger.
// Outputs short logs to the console and Write structured and detailed json logs to the log file.
func InitLogger() *zap.Logger {
	once.Do(func() {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		initZapLogger()
		Info("INIT_LOGGER")
	})
	return zapLogger
}

// See https://pkg.go.dev/go.uber.org/zap
func initZapLogger() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "function",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(newRotateLogs())),
		zapcore.DebugLevel,
	)
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).With(
		zap.String("version", getVersion()),
		zap.String("hostname", *getHost()),
	)
}

func getVersion() string {
	if version == "" {
		version = defaultVersion
	}

	return version
}

func SetVersion(version string) {
	version = version
}

func getHost() *string {
	ret, err := os.Hostname()
	if err != nil {
		log.Print(err)
		return nil
	}
	return &ret
}

func SetConsoleType(option ConsoleType) {
	consoleType = option
}
