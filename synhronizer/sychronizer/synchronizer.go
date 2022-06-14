package sychronizer

import "chatting/model"

type Synchronizer interface {
	Listen(handler ListeningHandler) error
	Synchronize([]byte) error
	SaveToRDB(model.Message) error
	Close()
}

type ListeningHandler func([]byte) error
