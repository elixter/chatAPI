package pubsub

import (
	"chatting/model"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisPubSub_PublishSubscribe(t *testing.T) {
	var ch chan model.Message

	type args struct {
		handler SubscribeHandler
	}
	tests := []struct {
		name    string
		args    args
		message model.Message
	}{
		// TODO: Add test cases.
		{
			name: "pubsub test",
			args: args{
				handler: func(bytes []byte) error {
					var msg model.Message
					err := json.Unmarshal(bytes, &msg)
					if err != nil {
						return err
					}
					ch <- msg
					close(ch)

					return nil
				},
			},
			message: model.Message{
				Id:             1,
				OriginServerId: uuid.New(),
				SyncServerId:   uuid.New(),
				MessageType:    model.TypeChatText,
				AuthorId:       1,
				RoomId:         1,
				Content:        "test",
				CreateAt:       time.Now().UTC(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			ch = make(chan model.Message)

			msg, err := json.Marshal(tt.message)
			if err != nil {
				t.Error(err)
			}

			r.Subscribe(tt.args.handler)

			err = r.Publish(msg)
			if err != nil {
				t.Error(err)
			}

			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						return
					}

					if !assert.Equal(t, tt.message, msg) {
						t.FailNow()
					}
				}
			}
		})
	}
}
