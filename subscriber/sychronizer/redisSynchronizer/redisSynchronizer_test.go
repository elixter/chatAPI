package redisSynchronizer

import (
	"chatting/config"
	"chatting/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		ctx := context.Background()

		message := model.Message{
			MessageType: model.TypeChatText,
			AuthorId:    int64(9999),
			RoomId:      int64(1234),
			Content:     "testing",
			CreateAt:    time.Now(),
		}
		msg, err := json.Marshal(message)
		if err != nil {
			return
		}

		go func() {
			err = rdb.Publish(ctx, "mychannel1", msg).Err()
			if err != nil {
				panic(err)
			}

		}()

		pubsub := rdb.Subscribe(ctx, "mychannel1")

		// Close the subscription when we are done.
		defer pubsub.Close()

		ch := pubsub.Channel()

		for msg := range ch {
			var target model.Message

			err := json.Unmarshal([]byte(msg.Payload), &target)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(target)
			break
		}
	})
}

type testRepository struct{}

func newTestRepository() testRepository {
	return testRepository{}
}

func (testRepository) Save(message model.Message) error {
	//TODO implement me
	panic("implement me")
}

func (testRepository) Close() {
	//TODO implement me
	panic("implement me")
}

func TestRedisSynchronizer(t *testing.T) {
	tests := []struct {
		name    string
		message model.Message
	}{
		{
			name: "redis synchronizer test",
			message: model.Message{
				Id:             1,
				OriginServerId: uuid.New(),
				SyncServerId:   config.ServerId,
				MessageType:    model.TypeChatText,
				AuthorId:       1,
				RoomId:         1,
				Content:        "test",
				CreateAt:       time.Now().UTC(),
			},
		},
	}
	for _, tt := range tests {
		synchronizer := New(newTestRepository())
		var ch chan model.Message

		t.Run(tt.name, func(t *testing.T) {
			ch = make(chan model.Message)

			synchronizer.Listen(func(bytes []byte) error {
				var msg model.Message
				err := json.Unmarshal(bytes, &msg)
				if err != nil {
					return err
				}

				ch <- msg
				close(ch)

				return nil
			})

			msg, err := json.Marshal(tt.message)
			if err != nil {
				t.Error(err)
			}
			synchronizer.Synchronize(msg)

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
