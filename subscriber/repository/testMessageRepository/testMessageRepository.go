package testMessageRepository

import (
	"chatting/logger"
	"chatting/model"
)

type TestMessageRepository struct{}

func New() *TestMessageRepository {
	return &TestMessageRepository{}
}

func (t *TestMessageRepository) Save(message model.Message) error {
	logger.Log.Debug(message)

	return nil
}

func (t *TestMessageRepository) Close() {
	logger.Log.Info("database closed")
}
