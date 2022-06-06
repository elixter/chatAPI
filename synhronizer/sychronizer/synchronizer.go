package synhronizer

import (
	"chatting/model"
)

type Synchronizer interface {
	Listen() error
	Synchronize([]byte) error
	SaveToRDB(model.Message) error
	Close()
}
