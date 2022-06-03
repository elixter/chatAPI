package logger

import (
	"github.com/labstack/echo/v4"
	"sync"
)

var Log echo.Logger
var once sync.Once

func New(logger echo.Logger) {
	once.Do(func() {
		Log = logger
	})
}
