package storageCluster

type StorageCluster interface {
	Receive() error
	Broadcast() error
}
