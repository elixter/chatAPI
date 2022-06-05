package storageCluster

import (
	"chatting/model"
)

type StorageCluster interface {
	Receive() error
	Synchronize([]byte) error
	SaveToRDB(model.Message) error
	Close()
}
