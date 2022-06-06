package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type Hub struct {
	rooms map[int64]*room
}

func NewHub() *Hub {
	return &Hub{
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
	if _, ok := h.rooms[roomId]; !ok {
		h.rooms[roomId] = newRoom(roomId)
		go h.rooms[roomId].run()
	}
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

	client := &Client{room: room, conn: conn, send: make(chan []byte, 256)}
	client.room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}