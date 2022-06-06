package main

import (
	"chatting/logger"
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

var Ctx context.Context
var serverId uuid.UUID

func main() {
	e := echo.New()
	e.Logger = logger.Log
	hub := NewHub()
	serverId = uuid.New()

	e.Logger.Info(serverId.String())

	e.GET("/", func(c echo.Context) error {
		logger.Log.Info("test")
		logger.Log.Debug("tetggggggg")
		serveHome(c.Response().Writer, c.Request())
		return nil
	})

	e.GET("/ws/:roomId", hub.WsHandler)
	e.Logger.Fatal(e.Start(":8080"))
}