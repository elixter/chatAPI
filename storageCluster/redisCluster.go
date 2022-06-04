package storageCluster

import (
	"chatting/logger"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type RedisCluster struct {
	redisClient *redis.Client
	mqConn      *amqp.Connection
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

func (rc *RedisCluster) Receive() error {
	//TODO implement me
	panic("implement me")
}

func (rc *RedisCluster) Broadcast() error {
	//TODO implement me
	panic("implement me")
}

func (rc *RedisCluster) Close() {
	rc.redisClient.Close()
	rc.mqConn.Close()
}
