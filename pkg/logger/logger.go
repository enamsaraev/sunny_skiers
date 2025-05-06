package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

var (
	logger *Logger
	once   sync.Once
)

func (l *Logger) Info(s string) {
	l.info.Println(s)
}

func (l *Logger) Infof(format string, args ...any) {
	l.info.Printf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.info.Printf(format, args...)
}

func CreateLogger() *Logger {
	once.Do(func() {
		logger = &Logger{
			info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})

	return logger
}

func GetLogger() *Logger {
	return logger
}
