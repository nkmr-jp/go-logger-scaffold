// Created from https://github.com/nkmr-jp/go-logger-scaffold
package logger

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	logFileDefault      = "./log/app_%Y-%m-%d.log"
	rotationTimeDefault = 24 * time.Hour
	purgeTimeDefault    = 7 * 24 * time.Hour
)

var (
	logFile      string
	rotationTime time.Duration
	purgeTime    time.Duration
)

// See https://github.com/lestrrat-go/file-rotatelogs
func newRotateLogs() *rotatelogs.RotateLogs {
	setRotateDefault()
	res, err := rotatelogs.New(
		logFile,
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithMaxAge(purgeTime),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("log file path: %v", logFile)
	return res
}

// SetLogFile set log file path ex. "./log/app_%Y-%m-%d.log"
func SetLogFile(file string) {
	logFile = file
}

func SetRotationTime(duration time.Duration) {
	rotationTime = duration
}

func SetPurgeTime(duration time.Duration) {
	purgeTime = duration
}

func setRotateDefault() {
	if logFile == "" {
		logFile = logFileDefault
	}
	if rotationTime == 0 {
		rotationTime = rotationTimeDefault
	}
	if purgeTime == 0 {
		purgeTime = purgeTimeDefault
	}
}
