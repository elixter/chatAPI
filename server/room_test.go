package main

import (
	"chatting/logger"
	"chatting/model"
	pubsub2 "chatting/pubsub"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_room_filterBroadcast(t *testing.T) {
	type fields struct {
		id         int64
		clients    map[*Client]bool
		broadcast  chan []byte
		register   chan *Client
		unregister chan *Client
	}
	type args struct {
		message []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &room{
				id:         tt.fields.id,
				clients:    tt.fields.clients,
				broadcast:  tt.fields.broadcast,
				register:   tt.fields.register,
				unregister: tt.fields.unregister,
			}
			got, err := r.filterBroadcast(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterBroadcast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("filterBroadcast() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_messageListening(t *testing.T) {
	serverId = uuid.New()

	type fields struct {
		room room
	}
	type args struct {
		msg model.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Same servierId Test",
			fields: fields{
				room{
					id:         123,
					clients:    make(map[*Client]bool),
					broadcast:  make(chan []byte),
					register:   make(chan *Client),
					unregister: make(chan *Client),
				},
			},
			args: args{
				model.Message{
					Id:             123,
					OriginServerId: serverId,
					SyncServerId:   uuid.New(),
					MessageType:    model.TypeChatText,
					AuthorId:       123,
					RoomId:         123,
					Content:        "asdf",
					CreateAt:       time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "diff serverId Test",
			fields: fields{
				room{
					id:         123,
					clients:    make(map[*Client]bool),
					broadcast:  make(chan []byte),
					register:   make(chan *Client),
					unregister: make(chan *Client),
				},
			},
			args: args{
				model.Message{
					Id:             123,
					OriginServerId: uuid.New(),
					SyncServerId:   uuid.New(),
					MessageType:    model.TypeChatText,
					AuthorId:       123,
					RoomId:         123,
					Content:        "asdf",
					CreateAt:       time.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			r := tt.fields.room
			newClient := &Client{
				id:   123,
				room: &r,
				conn: nil,
				send: make(chan []byte),
			}
			r.clients[newClient] = true

			bData, err := json.Marshal(tt.args.msg)
			if err != nil {
				t.Error(err)
			}

			go func() {
				defer close(newClient.send)
				err = r.messageListening(bData)
				if !tt.wantErr {
					if err != nil {
						if err != pubsub2.ErrMessageNoNeedToBroadcast {
							t.Error(err)
						}
					}
				}
				return
			}()

			if err == nil {
				message := <-newClient.send
				if !tt.wantErr {
					if !assert.Equal(t, bData, message) {
						t.FailNow()
					}
				}
			}
		})
	}
}

func Test_room_run(t *testing.T) {
	pubsub = pubsub2.New()

	type fields struct {
		room room
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "leak test",
			fields: fields{
				room{
					id:         123,
					clients:    make(map[*Client]bool),
					broadcast:  make(chan []byte),
					register:   make(chan *Client),
					unregister: make(chan *Client),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			r := tt.fields.room
			newClient := &Client{
				id:   123,
				room: &r,
				conn: nil,
				send: make(chan []byte),
			}
			r.clients[newClient] = true
			go r.run()

			r.unregister <- newClient
			<-newClient.send

			pubsub.Close()
		})
	}
}

func Benchmark_room_filterBroadcast(b *testing.B) {
	serverId = uuid.New()

	type fields struct {
		room room
	}
	type args struct {
		msg model.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "bench filterBroadcast",
			fields: fields{
				room: room{
					id:         123,
					clients:    make(map[*Client]bool),
					broadcast:  make(chan []byte),
					register:   make(chan *Client),
					unregister: make(chan *Client),
				},
			},
			args: args{
				msg: model.Message{
					Id:             123,
					OriginServerId: serverId,
					SyncServerId:   uuid.New(),
					MessageType:    model.TypeChatText,
					AuthorId:       123,
					RoomId:         123,
					Content:        "asdf",
					CreateAt:       time.Now(),
				},
			},
		},
	}

	msg := tests[0].args.msg
	data, err := json.Marshal(tests[0].args.msg)
	if err != nil {
		b.Error(err)
	}

	strData := fmt.Sprintf(
		"%d %s %s %s %d %d %s %s",
		msg.Id,
		msg.OriginServerId,
		msg.SyncServerId,
		msg.MessageType,
		msg.AuthorId,
		msg.RoomId,
		msg.Content,
		msg.CreateAt.String(),
	)

	b.Run("benchmark data format", func(b *testing.B) {
		broadcast, err := tests[0].fields.room.filterBroadcast(data)
		if err != nil {
			b.Error(err)
		}

		if broadcast != false {
			b.FailNow()
		}
	})

	b.Run("benchmark string data format", func(b *testing.B) {
		broadcast, err := tests[0].fields.room.stringDataFiltering([]byte(strData))
		if err != nil {
			b.Error(err)
		}

		if broadcast != false {
			b.FailNow()
		}
	})
}

func (r *room) stringDataFiltering(message []byte) (bool, error) {
	msg := string(message[:])
	tokenized := strings.Split(msg, " ")

	if tokenized[1] == serverId.String() && tokenized[2] != "" {
		logger.Debugf("message from same origin : [%s]", tokenized[1])
		return false, nil
	}

	if tokenized[5] != strconv.FormatInt(r.id, 10) {
		return false, nil
	}

	return true, nil
}
