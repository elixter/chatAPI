package redisSynchronizer

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"chatting/repository"
	"chatting/sychronizer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
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

func (r *RedisSynchronizer) Listen(handler sychronizer.ListeningHandler) error {
	logger.Log.Info("synchronizer starts listening")

	channelName := config.Config().GetString("redis.publishChannelName")

	sub := r.redis.Subscribe(r.ctx, channelName)
	msgs := sub.Channel()

	go func() {
		defer sub.Close()
		for msg := range msgs {
			logger.Log.Debug("received message")

			err := handler([]byte(msg.Payload))
			if err != nil {
				logger.Log.Error(err)
			}
		}
		logger.Log.Debug("goroutine in listening is ended")
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

func (r *RedisSynchronizer) ListeningHandler(payload []byte) error {
	message, err := binding(payload)
	if err != nil {
		return err
	}
	if checkSameOrigin(message) {
		logger.Log.Debugf("same origin : [%v]", config.ServerId)
		return nil
	}

	err = r.Synchronize(payload)
	if err != nil {
		return errors.Errorf("message synchronize failed: [%v]", err)
	}

	err = r.SaveToRDB(message)
	if err != nil {
		return errors.Errorf("saving message to RDB failed : [%v]", err)
	}

	return nil
}

func binding(payload []byte) (model.Message, error) {
	var result model.Message
	err := json.Unmarshal(payload, &result)
	if err != nil {
		return model.Message{}, err
	}

	return result, nil
}

func checkSameOrigin(message model.Message) bool {
	return message.SyncServerId == config.ServerId
}
