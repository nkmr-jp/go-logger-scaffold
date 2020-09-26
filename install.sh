RAW_FILE_PATH=https://raw.githubusercontent.com/nkmr-jp/go-logger-scaffold/master/logger
mkdir ./logger
curl -L "$RAW_FILE_PATH/logger.go" > ./logger/logger.go
curl -L "$RAW_FILE_PATH/wrapper.go" > ./logger/wrapper.go
