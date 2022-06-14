package redisSynchronizer

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"chatting/repository"
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
	logger.Log.Info("synchronizer starts listening")

	channelName := config.Config().GetString("redis.publishChannelName")

	sub := r.redis.Subscribe(r.ctx, channelName)

	msgs := sub.Channel()

	go func() {
		defer sub.Close()
		for msg := range msgs {
			logger.Log.Infof("received message")
			payload := []byte(msg.Payload)

			var received model.Message
			err := json.Unmarshal([]byte(msg.Payload), &received)
			if err != nil {
				logger.Log.Errorf("received message unmarshal error: [%v]", err)
			}
			if received.SyncServerId == config.ServerId {
				logger.Log.Infof("same sync server id")
				continue
			}

			go func(message []byte) {
				err := r.Synchronize(message)
				if err != nil {
					logger.Log.Errorf("message synchronize failed: [%v]", err)
				}
			}(payload)

			go func() {
				var message model.Message
				err = json.Unmarshal(payload, &message)
				if err != nil {
					logger.Log.Errorf("binding message body failed : [%v]", err)
				}

				err := r.SaveToRDB(message)
				if err != nil {
					logger.Log.Errorf("saving message to RDB failed : [%v]", err)
				}
			}()
		}
		logger.Log.Info("goroutine in listening is ended")
	}()

	return nil
}

func (r *RedisSynchronizer) Synchronize(message []byte) error {
	channelName := config.Config().GetString("redis.listeningChannelName")

	var msg model.Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return err
	}
	msg.SyncServerId = config.ServerId
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	message = marshal

	return r.redis.Publish(r.ctx, channelName, message).Err()
}

func (r *RedisSynchronizer) SaveToRDB(message model.Message) error {
	return r.repository.Save(message)
}

func (r *RedisSynchronizer) Close() {
	r.redis.Close()
	r.repository.Close()
}
