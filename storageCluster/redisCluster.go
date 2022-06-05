package storageCluster

import (
	"chatting/logger"
	"chatting/model"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
)

type RedisCluster struct {
	redisClient *redis.Client
	mqConn      *amqp.Connection
	rdbConn     *sqlx.DB
}

func New() RedisCluster {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Log.Panicf("connect with message queue failed : [%v]", err)
	}

	return RedisCluster{
		redisClient: rdb,
		mqConn:      conn,
	}
}

func (rc *RedisCluster) Listen() error {
	ch, err := rc.mqConn.Channel()
	if err != nil {
		logger.Log.Errorf("open message queue channel failed : [%v]", err)
		return err
	}

	msgs, err := ch.Consume(
		"queue name",
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
	panic("implement me")
}

func (rc *RedisCluster) SaveToRDB(message model.Message) error {
	//TODO : save Message to rdb
	panic("implement me")
}

func (rc *RedisCluster) Close() {
	rc.redisClient.Close()
	rc.mqConn.Close()
}
