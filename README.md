# go-logger-scaffold

Easily install and Edit rich logging capabilities into go projects.

## Features
- Output json structured log to file
  - You can also output information about the number of lines in the file, the method, the host, and the git version
  - These can be useful for debugging
  - Easily change the settings. ([logger/logger.go](https://github.com/nkmr-jp/go-logger-scaffold/blob/master/logger/logger.go#L32)).
  - Since it is a json structure, you can use the `jq` command to extract only the data you need.
  - [zap](https://github.com/uber-go/zap) use.
- Only a simple log output to console
  - It's hard to know what's going on when the console is flooded with detailed logs.
  - The output to the console is a minimal information.
  - [log](https://pkg.go.dev/log) use.
- It can jump to code from the console log when you 'Run' in GoLand.
  - This is why we use the standard log for console logs, not zap.
- Log file rotation
  - [file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs) use.

## Install

must have `go 1.14` or higher, and have `go mod init` running.

```sh
cd [your-project-path]
curl -sSf https://raw.githubusercontent.com/nkmr-jp/go-logger-scaffold/master/install.sh | sh
go mod vendor
```

## How to use
Please refer to [main.go](main.go)

```go
package main

import (
	"[your-project]/logger"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Sync() // flush log buffer

	logger.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
}
```

## Output

### Console
```sh
$ go run main.go
2020/09/27 03:29:37 logger.go:26: [INFO] INIT_LOGGER
2020/09/27 03:29:37 main.go:10: [INFO] USER_INFO
```

### File

```sh
ls log/
app-2020-09-27_03.log
```

#### jq example

```sh
# jq command install
$ brew install jq

# Extract only [INFO]
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

## Reference
* [zap package · pkg.go.dev](https://pkg.go.dev/go.uber.org/zap)
* [golangの高速な構造化ログライブラリ「zap」の使い方 - Qiita](https://qiita.com/emonuh/items/28dbee9bf2fe51d28153)
* [Adding Caller Depth support · Issue #715 · uber-go/zap · GitHub](https://github.com/uber-go/zap/issues/715)
* [logging - Is it possible to wrap log.Logger functions without losing the line number prefix? - Stack Overflow](https://stackoverflow.com/questions/42762391/is-it-possible-to-wrap-log-logger-functions-without-losing-the-line-number-prefi)
* [Time based log file rotation with zap — Dhwaneet Bhatt](https://dhwaneetbhatt.com/time-based-log-file-rotation-with-zap)
