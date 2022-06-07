package main

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"context"
	"encoding/json"
	"fmt"
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
		channelName := config.Config().GetString("redis.listeningChannelName")
		sub := rdb.Subscribe(context.Background(), channelName)
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

				if received.OriginServerId == serverId && received.SyncServerId.String() != "" {
					logger.Log.Infof("message is same origin : [%s]", received.OriginServerId.String())
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
			close(client.send)
			delete(r.clients, client)
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					// if client channel has issue, disconnect client
					logger.Log.Debugf("client [%d] channel has problem", client.id)
					delete(r.clients, client)
					_, ok := <-client.send
					if !ok {
						close(client.send)
					}
				}
			}
		}
	}
}
