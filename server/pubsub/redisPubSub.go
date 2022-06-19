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
	logger.Debugf("pub chan name : [%s]", cfg[publishConfigKey])
	return r.client.Publish(r.ctx, cfg[publishConfigKey], bytes).Err()
}

func (r *RedisPubSub) Subscribe(handler SubscribeHandler) {
	logger.Debugf("sub chan name : [%s]", cfg[listeningConfigKey])
	sub := r.client.Subscribe(r.ctx, cfg[listeningConfigKey])
	msgs := sub.Channel()

	go func() {
		defer sub.Close()
		for msg := range msgs {
			logger.Debug("sub got message")
			err := handler([]byte(msg.Payload))
			if err != nil {
				if err != ErrMessageNoNeedToBroadcast {
					logger.Error(err)
					return
				}
				continue
			}
		}
	}()
}
