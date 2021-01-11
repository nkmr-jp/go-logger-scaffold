// Created from https://github.com/nkmr-jp/go-logger-scaffold
package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

var consoleField string

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

// Wrapper of Zap's Debug.
func Debug(msg string, fields ...zap.Field) {
	wrapper(msg, "DEBUG").Debug(msg, fields...)
}

// Wrapper of Zap's Info.
// Outputs a short log to the console. Detailed json log output to log file.
func Info(msg string, fields ...zap.Field) {
	wrapper(msg, "INFO").Info(msg, fields...)
}

// Wrapper of Zap's Warn.
func Warn(msg string, fields ...zap.Field) {
	wrapper(msg, "WARN").Warn(msg, fields...)
}

// Wrapper of Zap's Error.
func Error(msg string, fields ...zap.Field) {
	wrapper(msg, "ERROR").Error(msg, fields...)
}

// Wrapper of Zap's Fatal.
func Fatal(msg string, fields ...zap.Field) {
	wrapper(msg, "FATAL").Fatal(msg, fields...)
}

// Outputs a Info log with formatted error.
func Infof(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "FATAL", err).Info(msg, addErrorField(fields, err)...)
}

// Outputs a Warn log with formatted error.
func Warnf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "FATAL", err).Warn(msg, addErrorField(fields, err)...)
}

// Outputs a Error log with formatted error.
func Errorf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "FATAL", err).Error(msg, addErrorField(fields, err)...)
}

// Outputs a Fatal log with formatted error.
func Fatalf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "FATAL", err).Fatal(msg, addErrorField(fields, err)...)
}

func wrapper(msg, level string) *zap.Logger {
	checkInit()
	shortLog(msg, level)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func wrapperf(msg, level string, err error) *zap.Logger {
	checkInit()
	shortLogWithError(msg, level, err)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func addErrorField(fields []zap.Field, err error) []zap.Field {
	return append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
}

// Short log to output to the console.
func shortLog(msg, level string) {
	if consoleType != ConsoleTypeAll {
		return
	}
	var str string
	if consoleField != "" {
		str = fmt.Sprintf("%v %v: %v", color(level), msg, Cyan.Add(consoleField))
		consoleField = ""
	} else {
		str = fmt.Sprintf("%v %v", color(level), msg)
	}
	err := log.Output(4, str)
	if err != nil {
		log.Fatal(err)
	}
}

// ConsoleField messages to be displayed on the console.
// It is recommended to use it for the minimum necessary short message.
func ConsoleField(str string) zap.Field {
	consoleField = str
	return zap.String("console", str)
}

// Short log to output to the console with error.
func shortLogWithError(msg string, level string, err error) {
	if consoleType == ConsoleTypeNone {
		return
	}
	err2 := log.Output(4, fmt.Sprintf("%v %v: %v", color(level), msg, Magenta.Add(err.Error())))
	if err2 != nil {
		log.Fatal(err2)
	}
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. InitLogger() must be called.")
	}
}

func color(level string) string {
	var color Color
	switch level {
	case "FATAL":
		color = Red
	case "ERROR":
		color = Red
	case "WARN":
		color = Yellow
	case "INFO":
		color = Green
	case "DEBUG":
		color = Green
	}
	return color.Add(level)
}

// Wrapper of pp.Print()
func Print(i interface{}) (n int, err error) {
	shortLog("pp.Print (console only)", "DEBUG")
	return pp.Print(i)
}

// Wrapper of pp.Println()
func Println(i interface{}) (n int, err error) {
	shortLog("pp.Println (console only)", "DEBUG")
	return pp.Println(i)
}

// Wrapper of spew.Dump()
func Dump(i interface{}) {
	shortLog("spew.Dump (console only)", "DEBUG")
	spew.Dump(i)
}

// See: https://github.com/uber-go/zap/blob/404189cf44aea95b0cd9bddcb0242dd4cf88c510/internal/color/color.go
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
