package pubsub

import "errors"

type SubscribeHandler func([]byte) error

var (
	ErrMessageNoNeedToBroadcast = errors.New("message no need to broadcast")
)

type PubSub interface {
	Publish([]byte) error
	Subscribe(handler SubscribeHandler, destruct chan bool)
}
