package main

import (
	"chatting/config"
	"chatting/repository/mySqlMeesageRepository"
	"chatting/sychronizer/redisSynchronizer"
	"github.com/google/uuid"
)

func main() {
	config.ServerId = uuid.New()

	synchronizer := redisSynchronizer.New(mySqlMeesageRepository.New())

	forever := make(chan bool)
	synchronizer.Listen(synchronizer.ListeningHandler)

	<-forever
}
