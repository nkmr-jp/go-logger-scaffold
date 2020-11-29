module github.com/nkmr-jp/go-logger-scaffold

go 1.14

replace github.com/nkmr-jp/go-logger-scaffold/logger => ./logger

require (
	github.com/nkmr-jp/go-logger-scaffold/logger v0.0.0-00010101000000-000000000000
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
)
