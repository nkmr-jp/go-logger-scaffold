# Function

- Output of json structured log (file output)
- Outputting the json structured log to a file
- Only a simple log output to the console
- GoLand can jump to files from the console log when executing 'Run'.
- Log file rotation

# Required
- go 1.14 and up


# Install to your project

You must have `go1.14` or higher, and have `go mod init` running.

```sh
cd [your-project]
curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/nkmr-jp/go-logger-scaffold/master/install.sh | sh
go mod vendor
```

# How to use
Please refer to [main.go](main.go)

```sh
$ go run main.go
2020/09/27 03:29:37 logger.go:26: [INFO] INIT_LOGGER
2020/09/27 03:29:37 main.go:10: [INFO] USER_INFO
```

# jq command sample
```sh
# jq command install
$ brew install jq

# Extract only level INFO
$ tail -fq log/*.log | jq -R 'fromjson? | select(.level=="INFO")'
{
  "level": "INFO",
  "ts": "2020-09-27T03:29:37.063635+09:00",
  "caller": "logger/logger.go:26",
  "function": "github.com/nkmr-jp/go-zap-scaffold/logger.InitLogger.func1",
  "msg": "INIT_LOGGER",
  "version": "ff97115",
  "hostname": "MacBook-Pro-16-inch-2019.local"
}
{
  "level": "INFO",
  "ts": "2020-09-27T03:29:37.064071+09:00",
  "caller": "go-logger-scaffold/main.go:10",
  "function": "main.main",
  "msg": "USER_INFO",
  "version": "ff97115",
  "hostname": "MacBook-Pro-16-inch-2019.local",
  "name": "Alice",
  "age": 20
}
```

# Reference
* [zap package · pkg.go.dev](https://pkg.go.dev/go.uber.org/zap)
* [golangの高速な構造化ログライブラリ「zap」の使い方 - Qiita](https://qiita.com/emonuh/items/28dbee9bf2fe51d28153)
* [Adding Caller Depth support · Issue #715 · uber-go/zap · GitHub](https://github.com/uber-go/zap/issues/715)
* [logging - Is it possible to wrap log.Logger functions without losing the line number prefix? - Stack Overflow](https://stackoverflow.com/questions/42762391/is-it-possible-to-wrap-log-logger-functions-without-losing-the-line-number-prefi)
* [Time based log file rotation with zap — Dhwaneet Bhatt](https://dhwaneetbhatt.com/time-based-log-file-rotation-with-zap)
