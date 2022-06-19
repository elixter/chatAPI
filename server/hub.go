package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

type Hub struct {
	mutex *sync.RWMutex
	rooms map[int64]*room
}

func NewHub() *Hub {
	return &Hub{
		mutex: &sync.RWMutex{},
		rooms: make(map[int64]*room),
	}
}

func (h *Hub) WsHandler(c echo.Context) error {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err,
		})
	}

	h.mutex.Lock()
	if _, ok := h.rooms[roomId]; !ok {
		h.rooms[roomId] = newRoom(roomId)
		go func() {
			h.rooms[roomId].run()
			delete(h.rooms, roomId)
		}()
	}
	h.mutex.Unlock()
	serveWs(h.rooms[roomId], c.Response().Writer, c.Request())

	return c.NoContent(http.StatusOK)
}

// serveWs handles websocket requests from the peer.
func serveWs(room *room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{id: rand.Int63(), room: room, conn: conn, send: make(chan []byte, 256)}
	client.room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
