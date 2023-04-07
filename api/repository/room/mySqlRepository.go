package room

import (
	"chatting/config"
	"chatting/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
)

const (
	queryInsert     = "INSERT INTO room(name, private, create_at) VALUES(:name, :private, :create_at)"
	queryDelete     = "DELETE FROM room WHERE id = ?"
	querySelectId   = "SELECT * FROM room WHERE id = ?"
	querySelectName = "SELECT * FROM room WHERE name = ?"
)

type MySqlRepository struct {
	db *sqlx.DB
}

var instance *MySqlRepository
var once sync.Once

func GetMySqlRepository() *MySqlRepository {
	once.Do(
		func() {
			instance = newMySqlRepository()
		},
	)

	return instance
}

func newMySqlRepository() *MySqlRepository {
	datasource := getDatasource()
	db, err := sqlx.Connect("mysql", datasource)
	if err != nil {
		panic(err)
	}

	return &MySqlRepository{
		db: db,
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

func (m *MySqlRepository) Save(room model.Room) (model.Room, error) {
	exec, err := m.db.NamedExec(queryInsert, room)
	if err != nil {
		return room, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return room, err
	}
	room.Id = id

	return room, nil
}

func (m *MySqlRepository) FindById(id int64) (model.Room, error) {
	var result model.Room
	err := m.db.Get(&result, querySelectId, id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *MySqlRepository) FindAllByName(name string) ([]model.Room, error) {
	var result []model.Room
	err := m.db.Select(&result, querySelectName, name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *MySqlRepository) DeleteById(id int64) error {
	_, err := m.db.Exec(queryDelete, id)

	return err
}
