package pubsub

import (
	"chatting/config"
	"chatting/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var cfg map[string]string

type RedisPubSub struct {
	client *redis.Client
	ctx    context.Context
}

func init() {
	cfg = config.Config().GetStringMapString("redis")
}

func New() *RedisPubSub {
	addr := fmt.Sprintf("%s:%s", cfg["host"], cfg["port"])
	database, err := strconv.ParseInt(cfg["database"], 10, 64)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg["password"], // no password set
		DB:       int(database),   // use default DB
	})

	return &RedisPubSub{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisPubSub) Publish(bytes []byte) error {
	return r.client.Publish(r.ctx, cfg["publishChannelName"], bytes).Err()
}

func (r *RedisPubSub) Subscribe(handler SubscribeHandler) {
	sub := r.client.Subscribe(r.ctx, cfg["listeningChannelName"])
	msgs := sub.Channel()

	go func() {
		defer sub.Close()
		for msg := range msgs {

			err := handler([]byte(msg.Payload))
			if err != nil {
				if err != ErrMessageNoNeedToBroadcast {
					logger.Error(err)
				}
				continue
			}
		}
	}()
}
