package main

import (
	pubsub2 "chatting/pubsub"
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	pubsub = pubsub2.New()
	defer pubsub.Close()

	e.GET("/", func(c echo.Context) error {
		serveHome(c.Response().Writer, c.Request())
		return nil
	})

	e.GET("/ws/:roomId", hub.WsHandler)

	// Graceful ShutDown
	// https://echo.labstack.com/cookbook/graceful-shutdown/
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
