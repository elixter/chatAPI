package main

import (
	"chatting/logger"
	"chatting/model"
	pubsub2 "chatting/pubsub"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

func BenchmarkClient_messageProcessing(b *testing.B) {
	pubsub = pubsub2.New()
	message := model.ClientMessage{
		MessageType: model.TypeChatText,
		AuthorId:    123,
		RoomId:      123,
		Content:     json.RawMessage("\"asdf\""),
		CreateAt:    time.Now(),
	}

	testRoom := newRoom(123)

	strData := fmt.Sprintf(
		"%s %d %d %s %s",
		message.MessageType,
		message.AuthorId,
		message.RoomId,
		message.Content,
		message.CreateAt.String(),
	)

	data, err := json.Marshal(message)
	if err != nil {
		b.Error(err)
	}

	type fields struct {
		client Client
	}
	type args struct {
		message []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "benchmark",
			fields: fields{
				client: Client{
					id:   123,
					room: testRoom,
					conn: nil,
				},
			},
			args: args{
				message: data,
			},
		},
	}
	for _, tt := range tests {
		testRoom.clients[&tt.fields.client] = true
		b.Run(tt.name, func(b *testing.B) {
			go func() {
				<-testRoom.broadcast
			}()

			err := tt.fields.client.messageProcessing(tt.args.message)
			if err != nil {
				b.Error(err)
			}
		})
	}

	b.Run("bench string data", func(b *testing.B) {
		go func() {
			<-testRoom.broadcast
		}()

		err := tests[0].fields.client.stringMessageProcessing([]byte(strData))
		if err != nil {
			b.Error(err)
		}
	})
}

func (c *Client) stringMessageProcessing(message []byte) error {
	start := time.Now()
	s := string(message)

	logger.Debug(s)

	strSentData := fmt.Sprintf(
		"%d %s %s %s",
		-1,
		serverId,
		uuid.New(),
		s,
	)
	sentData := []byte(strSentData)

	c.room.broadcast <- sentData
	err := pubsub.Publish(sentData)
	if err != nil {
		logger.Errorf("message publishing failed")
		return err
	}

	logger.Debug(time.Now().Sub(start))

	return nil
}
