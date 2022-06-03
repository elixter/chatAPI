package main

import (
	"chatting/logger"
	"flag"
	"github.com/labstack/echo/v4"
	log2 "github.com/labstack/gommon/log"
	"log"
	"net/http"
)

// TODO: 서버가 재시작 되어도 채팅방은 어떻게 유지시킬것인가
// TODO: 메세지는 어떻게 저장할것인가

var addr = flag.String("addr", ":8080", "http service address")

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

func main() {

	e := echo.New()
	e.Logger.SetLevel(log2.DEBUG)
	logger.New(e.Logger)
	hub := newHub()

	go hub.run()

	e.GET("/", func(c echo.Context) error {
		logger.Log.Info("test")
		logger.Log.Debug("tetggggggg")
		serveHome(c.Response().Writer, c.Request())
		return nil
	})

	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response().Writer, c.Request())

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
