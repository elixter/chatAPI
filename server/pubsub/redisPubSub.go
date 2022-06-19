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

const (
	passwordConfigKey  = "password"
	hostConfigKey      = "host"
	portConfigKey      = "port"
	databaseConfigKey  = "database"
	listeningConfigKey = "listeningchannelname"
	publishConfigKey   = "publishchannelname"
)

func init() {
	cfg = config.Config().GetStringMapString("redis")
}

func New() *RedisPubSub {
	addr := fmt.Sprintf("%s:%s", cfg["host"], cfg["port"])
	database, err := strconv.ParseInt(cfg[databaseConfigKey], 10, 64)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg[passwordConfigKey], // no password set
		DB:       int(database),          // use default DB
	})

	return &RedisPubSub{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisPubSub) Publish(bytes []byte) error {
	return r.client.Publish(r.ctx, cfg[publishConfigKey], bytes).Err()
}

func (r *RedisPubSub) Subscribe(handler SubscribeHandler, destruct chan bool) {
	sub := r.client.Subscribe(r.ctx, cfg[listeningConfigKey])
	msgs := sub.Channel()

	go func() {
		defer sub.Close()
		for {
			select {
			case <-destruct:
				return
			case msg := <-msgs:
				err := handler([]byte(msg.Payload))
				if err != nil {
					if err != ErrMessageNoNeedToBroadcast {
						logger.Error(err)
					}
					continue
				}
			}
		}
	}()
}
