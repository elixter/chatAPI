package storageCluster

import (
	"chatting/model"
)

type StorageCluster interface {
	Receive() error
	Synchronize() error
	SaveToRDB(model.Message) error
}
