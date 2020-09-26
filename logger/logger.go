package logger

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once
var zapLogger *zap.Logger

// Initialize the Logger.
// Outputs short logs to the console and Write structured and detailed json logs to the log file
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
		zap.DebugLevel,
	)
	zapLogger = zap.New(core, zap.AddCaller()).With(
		zap.String("version", *getVersion(Revision)),
		zap.String("hostname", *getHost()),
	)
}

// See https://pkg.go.dev/github.com/lestrrat-go/file-rotatelogs
func newRotateLogs() *rotatelogs.RotateLogs {
	logFile := "./log/app-%Y-%m-%d_%H.log"
	res, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

type VersionType int

const (
	_ VersionType = iota
	Revision
	Tag
)

// You can use the git revision or tag as a version.
// When using tag, recommend semantic versioning.
// See https://semver.org/
func getVersion(versionType VersionType) *string {
	var out []byte
	var err error

	switch versionType {
	case Revision:
		out, err = exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	case Tag:
		out, err = exec.Command("git", "tag").Output()
	}
	if err != nil {
		log.Print(err)
		return nil
	}

	ret := strings.TrimRight(string(out), "\n")
	return &ret
}

func getHost() *string {
	ret, err := os.Hostname()
	if err != nil {
		log.Print(err)
		return nil
	}
	return &ret
}
