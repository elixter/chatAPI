package main

import (
	"chatting/config"
	"chatting/repository/testMessageRepository"
	"chatting/sychronizer/redisSynchronizer"
	"github.com/google/uuid"
)

func main() {
	config.ServerId = uuid.New()

	synchronizer := redisSynchronizer.New(testMessageRepository.New())

	forever := make(chan bool)
	synchronizer.Listen(synchronizer.ListeningHandler)

	<-forever
}
