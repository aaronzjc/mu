package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrLevel   = logrus.ErrorLevel
	DebugLevel = logrus.DebugLevel
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})

	file, _ := os.OpenFile("/tmp/fuck.log", os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)
}

func Logger() *logrus.Logger {
	return log
}

func Write(level logrus.Level, format string, v ...interface{}) {
	host, _ := os.Hostname()
	ctxLog := log.WithFields(logrus.Fields{
		"host": host,
	})
	ctxLog.Logf(level, format, v...)
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
