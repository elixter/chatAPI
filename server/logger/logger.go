package logger

import (
	"github.com/labstack/echo/v4"
)

var Log echo.Logger

func SetLogger(logger echo.Logger) {
	Log = logger
}
