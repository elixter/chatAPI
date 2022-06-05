package main

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/bytes"
	"github.com/labstack/gommon/random"
	"github.com/streadway/amqp"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = bytes.MB
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the room.
type Client struct {
	id int64

	room *room

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the room.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	pongHandler := func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	}

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(pongHandler)

	for {
		_, message, err := c.conn.ReadMessage()
		logger.Log.Info("reading")
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Log.Errorf("error: %v", err)
			}
			break
		}

		broadcastMsg, err := messageProcessing(message)
		if err != nil {
			logger.Log.Errorf("message processing error : [%v]", err)
			continue
		}

		c.room.broadcast <- broadcastMsg
	}
}

// writePump pumps messages from the room to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func messageProcessing(message []byte) ([]byte, error) {
	readMessage := model.ClientMessage{}
	err := json.Unmarshal(message, &readMessage)
	if err != nil {
		logger.Log.Errorf("message unmarshalling error : [%v]", err)
		return nil, err
	}
	logger.Log.Debug(readMessage)

	// TODO : send message to Redis cluster
	conn, err := amqp.Dial(getMqSource())
	if err != nil {
		logger.Log.Panicf("connect with message queue failed : [%v]", err)
	}

	channel, err := conn.Channel()

	requestId := random.String(32)
	payload := amqp.Publishing{
		DeliveryMode:  amqp.Persistent,
		ContentType:   "application/json",
		CorrelationId: requestId,
		Body:          message,
	}

	queueName := config.Config().GetString("mq.listeningQueueName")

	err = channel.Publish("", queueName, false, false, payload)
	if err != nil {
		logger.Log.Errorf("publish failed : [%v]", err)
	}

	sentData, err := json.Marshal(readMessage)
	if err != nil {
		logger.Log.Errorf("message marshalling error : [%v]", err)
		return nil, err
	}

	return sentData, nil
}
