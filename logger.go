package testclient

import (
	"fmt"
)

// Logger is from https://github.com/go-log/log
// It is not imported to avoid dependency

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type defaultLogger struct{}

func (logger *defaultLogger) Log(v ...interface{}) {
	logger.log(fmt.Sprint(v...))
}

func (logger *defaultLogger) Logf(format string, v ...interface{}) {
	logger.log(fmt.Sprintf(format, v...))
}

func (logger *defaultLogger) log(entry string) {
	fmt.Println(entry)
}

func DefaultLogger() Logger {
	return &defaultLogger{}
}
