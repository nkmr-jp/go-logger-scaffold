// Created from https://github.com/nkmr-jp/go-logger-scaffold
package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/pp"
	. "github.com/logrusorgru/aurora"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
)

type Wrapper struct {
	Fields []zap.Field
}

// NewWrapper can additional fields.
// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
func NewWrapper(fields ...zap.Field) *Wrapper {
	return &Wrapper{Fields: fields}
}

func (w *Wrapper) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "DEBUG", fields).Debug(msg, fields...)
}

func (w *Wrapper) Info(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "INFO", fields).Info(msg, fields...)
}

func (w *Wrapper) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "WARN", fields).Warn(msg, fields...)
}

func (w *Wrapper) Error(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "ERROR", fields).Error(msg, fields...)
}

func (w *Wrapper) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "FATAL", fields).Fatal(msg, fields...)
}

func (w *Wrapper) Debugf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "DEBUG", err, fields).Debug(msg, fields...)
}

func (w *Wrapper) Infof(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "INFO", err, fields).Info(msg, fields...)
}

func (w *Wrapper) Warnf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "WARN", err, fields).Warn(msg, fields...)
}

func (w *Wrapper) Errorf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "ERROR", err, fields).Error(msg, fields...)
}

func (w *Wrapper) Fatalf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "FATAL", err, fields).Fatal(msg, fields...)
}

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

// Debug is Wrapper of Zap's Debug.
// Outputs a short log to the console. Detailed json log output to log file.
func Debug(msg string, fields ...zap.Field) {
	wrapper(msg, "DEBUG", fields).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	wrapper(msg, "INFO", fields).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	wrapper(msg, "WARN", fields).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	wrapper(msg, "ERROR", fields).Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	wrapper(msg, "FATAL", fields).Fatal(msg, fields...)
}

// Debugf is Outputs a Debug log with formatted error.
func Debugf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "DEBUG", err, fields).Debug(msg, addErrorField(fields, err)...)
}

func Infof(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "INFO", err, fields).Info(msg, addErrorField(fields, err)...)
}

func Warnf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "WARN", err, fields).Warn(msg, addErrorField(fields, err)...)
}

func Errorf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "ERROR", err, fields).Error(msg, addErrorField(fields, err)...)
}

func Fatalf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "FATAL", err, fields).Fatal(msg, addErrorField(fields, err)...)
}

func wrapper(msg, level string, fields []zap.Field) *zap.Logger {
	checkInit()
	shortLog(msg, level, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func wrapperf(msg, level string, err error, fields []zap.Field) *zap.Logger {
	checkInit()
	shortLogWithError(msg, level, err, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func addErrorField(fields []zap.Field, err error) []zap.Field {
	return append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
}

// Short log to output to the console.
func shortLog(msg, level string, fields []zap.Field) {
	if consoleType != ConsoleTypeAll {
		return
	}
	if outputType != OutputTypeSimpleConsoleAndFile {
		return
	}

	err := log.Output(4, fmt.Sprintf("%v %v%v", color(level), msg, getConsoleMsg(fields)))
	if err != nil {
		log.Fatal(err)
	}
}

func getConsoleMsg(fields []zap.Field) string {
	var ret string
	var consoles []string
	for _, v := range fields {
		if funk.ContainsString(consoleFields, v.Key) {
			var val string
			if v.String != "" {
				val = v.String
			} else {
				val = strconv.Itoa(int(v.Integer))
			}
			// consoles = append(consoles, fmt.Sprintf("%s=%s", v.Key, val))
			consoles = append(consoles, val)
		}
	}
	if consoles != nil {
		ret = ": " + fmt.Sprintf("%v", Cyan(strings.Join(consoles, ", ")))
	}
	return ret
}

// Short log to output to the console with error.
func shortLogWithError(msg string, level string, err error, fields []zap.Field) {
	if consoleType == ConsoleTypeNone {
		return
	}
	if outputType != OutputTypeSimpleConsoleAndFile {
		return
	}
	err2 := log.Output(
		4,
		fmt.Sprintf("%v %v: %v %v", color(level), msg, Magenta(err.Error()), getConsoleMsg(fields)),
	)
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
	switch level {
	case "FATAL":
		level = Red(level).String()
	case "ERROR":
		level = Red(level).String()
	case "WARN":
		level = Yellow(level).String()
	case "INFO":
		level = Green(level).String()
	case "DEBUG":
		level = Green(level).String()
	}
	return level
}

// Wrapper of pp.Print()
func Print(i interface{}) (n int, err error) {
	shortLog("pp.Print (console only)", "DEBUG", []zap.Field{})
	return pp.Print(i)
}

// Wrapper of pp.Println()
func Println(i interface{}) (n int, err error) {
	shortLog("pp.Println (console only)", "DEBUG", []zap.Field{})
	return pp.Println(i)
}

// Wrapper of spew.Dump()
func Dump(i interface{}) {
	shortLog("spew.Dump (console only)", "DEBUG", []zap.Field{})
	spew.Dump(i)
}
