package main

import (
	"chatting/logger"
	"chatting/model"
	pubsub2 "chatting/pubsub"
	"encoding/json"
	"github.com/pkg/errors"
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

func (r *room) run() {
	destruct := make(chan struct{})
	go pubsub.Subscribe(r.messageListening, destruct)

	for {
		select {
		case client := <-r.register:
			logger.Infof("Client [%d] entered room", client.id)
			r.clients[client] = true
		case client := <-r.unregister:
			logger.Infof("Client [%d] leaved room", client.id)
			close(client.send)
			delete(r.clients, client)

			if len(r.clients) == 0 {
				destruct <- struct{}{}
				close(destruct)
				logger.Info("Room socket destructed")
				return
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					// if client channel has issue, disconnect client
					logger.Debugf("client [%d] channel has problem", client.id)
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

func (r *room) messageListening(msg []byte) error {
	valid, err := r.filterBroadcast(msg)
	if err != nil {
		return errors.Errorf("failed to valid message : [%v]", err)
	}

	if !valid {
		return pubsub2.ErrMessageNoNeedToBroadcast
	}

	for client := range r.clients {
		client.send <- msg
	}

	return nil
}

func (r *room) filterBroadcast(message []byte) (bool, error) {
	var received model.Message
	err := json.Unmarshal(message, &received)
	if err != nil {
		return false, err
	}

	if received.OriginServerId == serverId && received.SyncServerId.String() != "" {
		logger.Debugf("message from same origin : [%s]", received.OriginServerId.String())
		return false, nil
	}

	if received.RoomId != r.id {
		return false, nil
	}

	return true, nil
}
