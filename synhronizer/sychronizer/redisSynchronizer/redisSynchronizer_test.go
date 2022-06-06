package redisSynchronizer

import (
	"chatting/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
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
