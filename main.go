package main

import (
	"chatting/logger"
	"chatting/synhronizer/redisSynchronizer"
	"chatting/synhronizer/repository/mySqlMeesageRepository"
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

func main() {

	e := echo.New()
	e.Logger = logger.Log
	hub := NewHub()

	cluster := redisSynchronizer.New(mySqlMeesageRepository.New())
	defer cluster.Close()
	go func() {
		err := cluster.Listen()
		if err != nil {
			logger.Log.Panicf("cluster listening failed : [%v]", err)
		}
	}()

	e.GET("/", func(c echo.Context) error {
		logger.Log.Info("test")
		logger.Log.Debug("tetggggggg")
		serveHome(c.Response().Writer, c.Request())
		return nil
	})

	e.GET("/ws/:roomId", hub.WsHandler)
	e.Logger.Fatal(e.Start(":8081"))
}
