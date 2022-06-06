package redisSynchronizer

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"chatting/synhronizer/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisSynchronizer struct {
	ctx        context.Context
	redis      *redis.Client
	repository repository.MessageRepository
}

func New(repository repository.MessageRepository) RedisSynchronizer {
	redisConfig := config.Config().GetStringMapString("redis")

	db, err := strconv.ParseInt(redisConfig["database"], 10, 32)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     getRedisSource(),
		Password: redisConfig["password"],
		DB:       int(db),
	})

	return RedisSynchronizer{
		ctx:        context.Background(),
		redis:      rdb,
		repository: repository,
	}
}

func getRedisSource() string {
	redisConfig := config.Config().GetStringMapString("redis")

	return fmt.Sprintf("%s:%s", redisConfig["host"], redisConfig["port"])
}

func (r *RedisSynchronizer) Listen() error {
	channelName := config.Config().GetString("redis.listeningChannelName")

	sub := r.redis.Subscribe(r.ctx, channelName)

	msgs := sub.Channel()

	go func() {
		defer sub.Close()
		for msg := range msgs {
			payload := []byte(msg.Payload)

			go func(message []byte) {
				err := r.Synchronize(message)
				if err != nil {
					logger.Log.Errorf("message synchronize failed: [%v]", err)
				}
			}(payload)

			var body model.Message
			err := json.Unmarshal(payload, &body)
			if err != nil {
				logger.Log.Errorf("binding message body failed : [%v]", err)
			}

			go func(message model.Message) {
				err := r.SaveToRDB(message)
				if err != nil {
					logger.Log.Errorf("saving message to RDB failed : [%v]", err)
				}
			}(body)
		}
		logger.Log.Info("goroutine in listening is ended")
	}()

	return nil
}

func (r *RedisSynchronizer) Synchronize(message []byte) error {
	channelName := config.Config().GetString("redis.listeningChannelName")

	return r.redis.Publish(r.ctx, channelName, message).Err()
}

func (r *RedisSynchronizer) SaveToRDB(message model.Message) error {
	return r.repository.Save(message)
}

func (r *RedisSynchronizer) Close() {
	r.redis.Close()
	r.repository.Close()
}
