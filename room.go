package main

import (
	"chatting/config"
	"chatting/logger"
	"fmt"
	"github.com/streadway/amqp"
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
		conn, err := amqp.Dial(getMqSource())
		if err != nil {
			logger.Log.Panicf("connect with message queue failed : [%v]", err)
		}

		channel, err := conn.Channel()
		msgs, err := channel.Consume(
			config.Config().GetString("mq.listeningQueueName"),
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logger.Log.Errorf("receiving message failed : [%v]", err)
			return
		}
		go func() {
			for msg := range msgs {
				r.broadcast <- msg.Body
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
