package repository

import (
	"chatting/config"
	"chatting/logger"
	"chatting/model"
	"fmt"
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

	conn, err := sqlx.Connect("mysql", getDatasource())
	if err != nil {
		logger.Log.Panicf("open mysql failed: [%v]", err)
	}

	return &MySqlMessageRepository{
		db: conn,
	}
}

func getDatasource() string {
	dbConfig := config.Config().GetStringMapString("db")

	return fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbConfig["id"],
		dbConfig["password"],
		dbConfig["host"],
		dbConfig["port"],
		dbConfig["database"],
	)
}

func (m *MySqlMessageRepository) Save(message model.Message) error {

	_, err := m.db.NamedExec(insertMessage, message)
	return err
}
