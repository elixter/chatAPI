package main

import (
	"chatting/logger"
	"chatting/model"
	pubsub2 "chatting/pubsub"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

type ctxKey int

const (
	ctxRoomId ctxKey = iota
)

type Room struct {
	Id         int64
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	ctx        context.Context
}

func newRoom(id int64) *Room {
	return &Room{
		Id:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		ctx:        context.Background(),
	}
}

func (r *Room) run() {
	ctx, cancel := context.WithCancel(r.ctx)
	ctx = context.WithValue(ctx, ctxRoomId, r.Id)
	go pubsub.Subscribe(r.messageListening, ctx)

	for {
		select {
		case client := <-r.Register:
			logger.Infof("Client [%d] entered Room", client.id)
			r.Clients[client] = true
		case client := <-r.Unregister:
			logger.Infof("Client [%d] leaved Room", client.id)
			close(client.send)
			delete(r.Clients, client)

			if len(r.Clients) == 0 {
				cancel()
				logger.Infof("Room[%d] socket closed", ctx.Value(ctxRoomId))
				return
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.send <- message:
				default:
					// if client channel has issue, disconnect client
					logger.Debugf("client [%d] channel has problem", client.id)
					delete(r.Clients, client)
					_, ok := <-client.send
					if !ok {
						close(client.send)
					}
				}
			}
		}
	}

}

func (r *Room) messageListening(msg []byte) error {
	valid, err := r.filterBroadcast(msg)
	if err != nil {
		return errors.Errorf("failed to valid message : [%v]", err)
	}

	if !valid {
		return pubsub2.ErrMessageNoNeedToBroadcast
	}

	for client := range r.Clients {
		client.send <- msg
	}

	return nil
}

func (r *Room) filterBroadcast(message []byte) (bool, error) {
	var received model.Message
	err := json.Unmarshal(message, &received)
	if err != nil {
		return false, err
	}

	if received.OriginServerId == serverId && received.SyncServerId.String() != "" {
		logger.Debugf("message from same origin : [%s]", received.OriginServerId.String())
		return false, nil
	}

	if received.RoomId != r.Id {
		return false, nil
	}

	return true, nil
}
