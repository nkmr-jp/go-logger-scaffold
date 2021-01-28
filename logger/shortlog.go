// Created from https://github.com/nkmr-jp/go-logger-scaffold
package logger

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	. "github.com/logrusorgru/aurora"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Short log to output to the console.
func shortLog(msg, levelStr string, fields []zap.Field) {
	if consoleType != ConsoleTypeAll {
		return
	}
	if outputType != OutputTypeShortConsoleAndFile {
		return
	}
	if !checkLevel(levelStr) {
		return
	}

	err := log.Output(4, fmt.Sprintf("%v %v%v", color(levelStr), msg, getConsoleMsg(fields)))
	if err != nil {
		log.Fatal(err)
	}
}

// Short log to output to the console with error.
func shortLogWithError(msg string, levelStr string, err error, fields []zap.Field) {
	if consoleType == ConsoleTypeNone {
		return
	}
	if outputType != OutputTypeShortConsoleAndFile {
		return
	}
	if !checkLevel(levelStr) {
		return
	}
	err2 := log.Output(
		4,
		fmt.Sprintf("%v %v: %v %v", color(levelStr), msg, Magenta(err.Error()), getConsoleMsg(fields)),
	)
	if err2 != nil {
		log.Fatal(err2)
	}
}

func checkLevel(levelStr string) bool {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return false
	}
	if logLevel > level {
		return false
	}
	return true
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
