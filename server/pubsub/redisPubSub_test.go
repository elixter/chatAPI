package pubsub

import (
	"chatting/logger"
	"chatting/model"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
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
			defer goleak.VerifyNone(t)

			r := New()
			ch = make(chan model.Message)
			destruct := make(chan struct{})

			msg, err := json.Marshal(tt.message)
			if err != nil {
				t.Error(err)
			}

			r.Subscribe(tt.args.handler, destruct)

			err = r.Publish(msg)
			if err != nil {
				t.Error(err)
			}

			msg2 := <-ch
			if !assert.Equal(t, tt.message, msg2) {
				t.FailNow()
			}

			destruct <- struct{}{}
			r.client.Close()

			err = goleak.Find()
			if err != nil {
				logger.Error(err)
			}
		})
	}
}
