package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

// Wrapper of Zap's Sync.
func Sync() {
	if err := zapLogger.Sync(); err != nil {
		log.Fatal(err)
	}
}

// Wrapper of Zap's Info.
// Outputs a short log to the console. Detailed json log output to log file.
func Info(msg string, fields ...zap.Field) {
	shortLog(msg, "INFO")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Wrapper of Zap's Debug.
func Debug(msg string, fields ...zap.Field) {
	shortLog(msg, "DEBUG")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Wrapper of Zap's Warn.
func Warn(msg string, fields ...zap.Field) {
	shortLog(msg, "WARN")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Wrapper of Zap's Error.
func Error(msg string, fields ...zap.Field) {
	shortLog(msg, "ERROR")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Wrapper of Zap's Fatal.
func Fatal(msg string, fields ...zap.Field) {
	shortLog(msg, "FATAL")
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Outputs a Error log with formatted error.
func Errorf(msg string, err error, fields ...zap.Field) {
	shortLogWithError(msg, "ERROR", err)
	fields = append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Outputs a Fatal log with formatted error.
func Fatalf(msg string, err error, fields ...zap.Field) {
	shortLogWithError(msg, "FATAL", err)
	fields = append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
	zapLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// Short log to output to the console
func shortLog(msg string, level string) {
	err := log.Output(3, fmt.Sprintf("[%v] %v", level, msg))
	if err != nil {
		log.Fatal(err)
	}
}

// Short log to output to the console with error
func shortLogWithError(msg string, level string, err error) {
	err2 := log.Output(3, fmt.Sprintf("[%v] %v error: %+v", level, msg, err))
	if err2 != nil {
		log.Fatal(err2)
	}
}
