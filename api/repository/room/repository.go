package room

import "chatting/model"

type Repository interface {
	Save(room model.Room) (model.Room, error)
	FindById(id int64) (model.Room, error)
	FindAllByName(name string) ([]model.Room, error)
	DeleteById(id int64) error
}
