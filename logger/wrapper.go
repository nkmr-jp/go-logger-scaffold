// Create a scaffold by https://github.com/nkmr-jp/go-logger-scaffold
package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

// Wrapper of Zap's Sync.
func Sync() {
	Info("FLUSH_LOG_BUFFER")
	if err := zapLogger.Sync(); err != nil {
		log.Fatal(err)
	}
}

// flush log buffer. when interrupt or terminated.
func SyncWhenStop() {
	c := make(chan os.Signal, 1)

	go func() {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		s := <-c

		sigCode := 0
		switch s.String() {
		case "interrupt":
			sigCode = 2
		case "terminated":
			sigCode = 15
		}

		Info(fmt.Sprintf("GOT_SIGNAL_%v", strings.ToUpper(s.String())))
		Sync() // flush log buffer
		os.Exit(128 + sigCode)
	}()
}

// Wrapper of Zap's Info.
// Outputs a short log to the console. Detailed json log output to log file.
func Info(msg string, fields ...zap.Field) {
	checkInit()
	shortLog(msg, "INFO")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Wrapper of Zap's Debug.
func Debug(msg string, fields ...zap.Field) {
	checkInit()
	shortLog(msg, "DEBUG")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Wrapper of Zap's Warn.
func Warn(msg string, fields ...zap.Field) {
	checkInit()
	shortLog(msg, "WARN")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Wrapper of Zap's Error.
func Error(msg string, fields ...zap.Field) {
	checkInit()
	shortLog(msg, "ERROR")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Wrapper of Zap's Fatal.
func Fatal(msg string, fields ...zap.Field) {
	checkInit()
	shortLog(msg, "FATAL")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Outputs a Error log with formatted error.
func Errorf(msg string, err error, fields ...zap.Field) {
	checkInit()
	shortLogWithError(msg, "ERROR", err)
	fields = append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Outputs a Fatal log with formatted error.
func Fatalf(msg string, err error, fields ...zap.Field) {
	checkInit()
	shortLogWithError(msg, "FATAL", err)
	fields = append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Short log to output to the console.
func shortLog(msg string, level string) {
	err := log.Output(3, fmt.Sprintf("[%v] %v", level, msg))
	if err != nil {
		log.Fatal(err)
	}
}

// Short log to output to the console with error.
func shortLogWithError(msg string, level string, err error) {
	err2 := log.Output(3, fmt.Sprintf("[%v] %v error: %+v", level, msg, err))
	if err2 != nil {
		log.Fatal(err2)
	}
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. InitLogger() must be called.")
	}
}
