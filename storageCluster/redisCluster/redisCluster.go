package redisCluster

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"chatting/storageCluster/repository"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type RedisCluster struct {
	redisClient *redis.Client
	mqConn      *amqp.Connection
	mqChan      *amqp.Channel
	repository  repository.MessageRepository
}

func New(repository repository.MessageRepository) RedisCluster {
	rdb := redis.NewClient(&redis.Options{
		Addr:     getRedisSource(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn, err := amqp.Dial(getMqSource())
	if err != nil {
		logger.Log.Panicf("connect with message queue failed : [%v]", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Panicf("open message queue channel failed : [%v]", err)
	}

	return RedisCluster{
		redisClient: rdb,
		mqConn:      conn,
		mqChan:      ch,
		repository:  repository,
	}
}

func getRedisSource() string {
	redisConfig := config.Config().GetStringMapString("redis")
	return fmt.Sprintf(
		"%s:%s",
		redisConfig["host"],
		redisConfig["port"],
	)
}

func getMqSource() string {
	mqConfig := config.Config().GetStringMapString("mq")
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		mqConfig["id"],
		mqConfig["password"],
		mqConfig["host"],
		mqConfig["port"],
	)
}

func (rc *RedisCluster) Listen() error {
	queueName := config.Config().GetString("mq.listeningQueueName")

	msgs, err := rc.mqChan.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Errorf("receiving message failed : [%v]", err)
		return err
	}

	go func() {
		for msg := range msgs {
			go func() {
				err := rc.Synchronize()
				if err != nil {
					logger.Log.Errorf("message synchronize failed : [%v]", err)
				}
			}()

			var body model.Message
			err := json.Unmarshal(msg.Body, &body)
			if err != nil {
				logger.Log.Errorf("binding message body failed : [%v]", err)
			}
			go func(message model.Message) {
				err := rc.SaveToRDB(message)
				if err != nil {
					logger.Log.Errorf("saving message to RDB failed : [%v]", err)
				}
			}(body)
		}
	}()

	return nil
}

func (rc *RedisCluster) Synchronize() error {
	//TODO : Synchronizing with other server
	logger.Log.Info("sync!")
	return nil
}

func (rc *RedisCluster) SaveToRDB(message model.Message) error {
	//TODO : save Message to rdb
	return rc.repository.Save(message)
}

func (rc *RedisCluster) Close() {
	rc.redisClient.Close()
	rc.mqChan.Close()
	rc.mqConn.Close()
	rc.repository.Close()
}
