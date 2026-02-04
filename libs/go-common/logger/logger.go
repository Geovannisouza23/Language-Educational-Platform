package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	return log
}

func Info(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Info(msg)
}

func Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	log.WithFields(fields).Error(msg)
}

func Debug(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Debug(msg)
}

func Warn(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Warn(msg)
}
