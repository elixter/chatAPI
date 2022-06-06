package mqSynchronizer

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"chatting/synhronizer/repository"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/streadway/amqp"
)

type MqSynchronizer struct {
	mqConn     *amqp.Connection
	mqChan     *amqp.Channel
	repository repository.MessageRepository
}

func New(repository repository.MessageRepository) MqSynchronizer {
	conn, err := amqp.Dial(getMqSource())
	if err != nil {
		logger.Log.Panicf("connect with message queue failed : [%v]", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Panicf("open message queue channel failed : [%v]", err)
	}

	return MqSynchronizer{
		mqConn:     conn,
		mqChan:     ch,
		repository: repository,
	}
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

func (rc *MqSynchronizer) Listen() error {
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
			go func(message []byte) {
				err := rc.Synchronize(message)
				if err != nil {
					logger.Log.Errorf("message synchronize failed : [%v]", err)
				}
			}(msg.Body)

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

func (rc *MqSynchronizer) Synchronize(message []byte) error {
	requestId := random.String(32)
	payload := amqp.Publishing{
		DeliveryMode:  amqp.Persistent,
		ContentType:   "application/json",
		CorrelationId: requestId,
		Body:          message,
	}

	queueName := config.Config().GetString("mq.listeningQueueName")

	return rc.mqChan.Publish(
		"",
		queueName,
		false,
		false,
		payload,
	)
}

func (rc *MqSynchronizer) SaveToRDB(message model.Message) error {
	//TODO : save Message to rdb
	return rc.repository.Save(message)
}

func (rc *MqSynchronizer) Close() {
	rc.mqChan.Close()
	rc.mqConn.Close()
	rc.repository.Close()
}
