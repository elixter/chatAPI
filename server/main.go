package main

import (
	"chatting/config"
	"chatting/logger"
	pubsub2 "chatting/pubsub"
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log2 "github.com/labstack/gommon/log"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

const (
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
)

var Ctx context.Context
var serverId uuid.UUID
var pubsub pubsub2.PubSub

var e *echo.Echo

func init() {
	e = echo.New()

	logConfig := config.Config().GetStringMapString("logger")
	switch logConfig["level"] {
	case logLevelDebug:
		e.Logger.SetLevel(log2.DEBUG)
		break
	case logLevelInfo:
		e.Logger.SetLevel(log2.INFO)
		break
	case logLevelWarn:
		e.Logger.SetLevel(log2.WARN)
		break
	case logLevelError:
		e.Logger.SetLevel(log2.ERROR)
		break
	}

	logger.SetLogger(e.Logger)
}

func main() {
	hub := NewHub()
	serverId = uuid.New()
	e.Logger.Infof("server id : [%s]", serverId.String())

	pubsub = pubsub2.New()

	e.GET("/", func(c echo.Context) error {
		logger.Log.Info("test")
		logger.Log.Debug("tetggggggg")
		serveHome(c.Response().Writer, c.Request())
		return nil
	})

	e.GET("/ws/:roomId", hub.WsHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
