package repository

import (
	"chatting/model"
	"github.com/jmoiron/sqlx"
)

const (
	insertMessage = "INSERT INTO messages(message_type, author_id, room_id, content, create_at) VALUES(:message_type, :author_id, :room_id, :content, :create_at)"
)

type MySqlMessageRepository struct {
	db *sqlx.DB
}

func (m *MySqlMessageRepository) Save(message model.Message) error {
	_, err := m.db.Exec(insertMessage, message)
	return err
}
