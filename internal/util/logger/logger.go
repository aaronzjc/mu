package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	InfoLevel  = "info"
	WarnLevel  = "warning"
	ErrLevel   = "error"
	DebugLevel = "debug"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	file, _ := os.OpenFile("/var/log/mu.log", os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)
}

func Logger() *logrus.Logger {
	return log
}

func Write(level, format string, v ...interface{}) {
	log.Printf("["+level+"] "+format, v...)
}
func Info(format string, v ...interface{}) {
	Write(InfoLevel, format, v...)
}

func Warning(format string, v ...interface{}) {
	Write(WarnLevel, format, v...)
}

func Error(format string, v ...interface{}) {
	Write(ErrLevel, format, v...)
}

func Debug(format string, v ...interface{}) {
	if gin.Mode() == gin.DebugMode {
		Write(DebugLevel, format, v...)
	}
}

func Fatal(v interface{}) {
	log.Fatal(v)
}
