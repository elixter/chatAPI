package logger

import (
	"chatting/config"
	"github.com/labstack/echo/v4"
	log2 "github.com/labstack/gommon/log"
)

var Log echo.Logger

const (
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
)

func init() {
	Log = log2.New("echo")
	logConfig := config.Config().GetStringMapString("logger")
	switch logConfig["level"] {
	case logLevelDebug:
		Log.SetLevel(log2.DEBUG)
		break
	case logLevelInfo:
		Log.SetLevel(log2.INFO)
		break
	case logLevelWarn:
		Log.SetLevel(log2.WARN)
		break
	case logLevelError:
		Log.SetLevel(log2.ERROR)
		break
	}

}
