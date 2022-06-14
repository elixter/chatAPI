package logger

import (
	"github.com/labstack/echo/v4"
)

var log echo.Logger

func SetLogger(logger echo.Logger) {
	log = logger
}

func Debug(i ...interface{}) {
	log.Debug(i)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Error(i ...interface{}) {
	log.Error(i)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Info(i ...interface{}) {
	log.Info(i)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Print(i ...interface{}) {
	log.Print(i)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args)
}
