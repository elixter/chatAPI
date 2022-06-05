package synhronizer

import (
	"chatting/model"
)

type Synchronizer interface {
	Receive() error
	Synchronize([]byte) error
	SaveToRDB(model.Message) error
	Close()
}
