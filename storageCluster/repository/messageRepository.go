package repository

import "chatting/model"

type MessageRepository interface {
	Save(model.Message) error
}
