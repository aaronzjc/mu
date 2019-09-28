package logger

import (
	"github.com/gin-gonic/gin"
	"log"
)

const (
	InfoLevel = "info"
	WarnLevel = "warning"
	ErrLevel = "error"
	DebugLevel = "debug"
)

func Write(level, format string, v ...interface{}) {
	log.Printf("[" + level + "] " + format + "\n", v...)
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