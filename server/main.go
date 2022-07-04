package main

import (
	pubsub2 "chatting/pubsub"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

var serverId uuid.UUID
var pubsub pubsub2.PubSub

var e *echo.Echo

func main() {
	e = echo.New()
	hub := NewHub()
	serverId = uuid.New()
	e.Logger.Infof("server Id : [%s]", serverId.String())

	e.Use(middleware.Logger())

	pubsub = pubsub2.New()
	defer pubsub.Close()

	e.GET("/", func(c echo.Context) error {
		serveHome(c.Response().Writer, c.Request())
		return nil
	})

	e.GET("/room/:id", EnterRoom)
	e.GET("/ws/:roomId", hub.WsHandler)

	e.POST("/room", CreateRoom)

	e.Logger.Fatal(e.Start(":8080"))
}
