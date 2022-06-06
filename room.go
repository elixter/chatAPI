package main

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type room struct {
	id         int64
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newRoom(id int64) *room {
	return &room{
		id:         id,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func getMqSource() string {
	mqConfig := config.Config().GetStringMapString("mq")
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		mqConfig["id"],
		mqConfig["password"],
		mqConfig["host"],
		mqConfig["port"],
	)
}

func (r *room) run() {
	go func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		sub := rdb.Subscribe(context.Background(), "chat")
		msgs := sub.Channel()

		go func() {
			defer sub.Close()
			for msg := range msgs {
				// TODO: 발신지와 같은 서버인 경우 continue
				var received model.Message
				err := json.Unmarshal([]byte(msg.Payload), &received)
				if err != nil {
					logger.Log.Error(err)
				}

				if received.ServerUUID.String() == serverId.String() {
					logger.Log.Infof("message is same origin : [%s]", received.ServerUUID.String())
					continue
				}

				r.broadcast <- []byte(msg.Payload)
			}
		}()
	}()

	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			delete(r.clients, client)
			close(client.send)
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					// if client channel has issue, disconnect client
					logger.Log.Debug("client [%d] channel has problem", client.id)
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
