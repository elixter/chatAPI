package repository

import (
	"chatting/logger"
	"chatting/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	insertMessage = "INSERT INTO messages(message_type, author_id, room_id, content, create_at) VALUES(:message_type, :author_id, :room_id, :content, :create_at)"
)

type MySqlMessageRepository struct {
	db *sqlx.DB
}

func New() *MySqlMessageRepository {
	conn, err := sqlx.Connect("mysql", "")
	if err != nil {
		logger.Log.Panicf("open mysql failed: [%v]", err)
	}

	return &MySqlMessageRepository{
		db: conn,
	}
}

func (m *MySqlMessageRepository) Save(message model.Message) error {
	_, err := m.db.Exec(insertMessage, message)
	return err
}
