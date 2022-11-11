package logger

import (
	"net/http"
	"os"
	"runtime"

	"github.com/aaronzjc/mu/pkg/helper"
	"github.com/sirupsen/logrus"
)

const (
	TYPE_COMMON = "common"
	TYPE_REQ    = "request"
)

type AppLogger struct {
	logger *logrus.Entry
}

var appLogger *AppLogger

func SetLevel(l string) {
	level, _ := logrus.ParseLevel(l)
	appLogger.logger.Logger.SetLevel(level)
}

func Fatal(args ...interface{}) {
	appLogger.logger.Fatal(args...)
}

func Info(args ...interface{}) {
	appLogger.logger.Info(args...)
}

func Debug(args ...interface{}) {
	appLogger.logger.Debug(args...)
}

func Error(args ...interface{}) {
	appLogger.logger.Error(args...)
}

func ErrorWithStack(args ...interface{}) {
	stack := make([]byte, 2048)
	stack = stack[:runtime.Stack(stack, true)]
	appLogger.logger.WithFields(logrus.Fields{
		"stack": string(stack),
	}).Error(args...)
}

func init() {
	appLogger = &AppLogger{
		logger: logrus.NewEntry(logrus.New()),
	}
}

func Request(req *http.Request, resp *http.Response, ts float64, err error) {
	fields := logrus.Fields{
		"consume": ts,
	}
	if req != nil {
		fields["req_host"] = req.URL.Host
		fields["req_path"] = req.URL.Path
		fields["req_params"] = req.URL.RawQuery
	}
	if resp != nil {
		fields["resp_code"] = resp.StatusCode
	}
	if err != nil {
		fields["err"] = err.Error()
		appLogger.logger.WithFields(fields).Error()
	} else {
		appLogger.logger.WithFields(fields).Info()
	}
}

func Setup(appName string, path string) error {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if path != "" {
		if f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
			return err
		} else {
			log.SetOutput(f)
		}
	}
	appLogger.logger = log.WithFields(logrus.Fields{
		"app_name": appName,
		"hostname": helper.LocalHostname(),
		"ip":       helper.LocalAddr(),
	})
	return nil
}
