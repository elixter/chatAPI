package redisCluster

import (
	"chatting/model"
	"chatting/storageCluster/repository/mySqlMeesageRepository"
	"testing"
	"time"
)

func TestRedisCluster_Listen(t *testing.T) {
	rc := New(mySqlMeesageRepository.New())
	defer rc.Close()

	t.Run("Listening test", func(t *testing.T) {
		if err := rc.Listen(); err != nil {
			t.Errorf("Listen() error = %v", err)
		}
	})
}

func TestRedisCluster_SaveToRDB(t *testing.T) {
	message := model.Message{
		MessageType: model.TypeChatText,
		AuthorId:    int64(9999),
		RoomId:      int64(1234),
		Content:     "testing",
		CreateAt:    time.Now(),
	}

	rc := New(mySqlMeesageRepository.New())
	defer rc.Close()

	t.Run("SaveToRDB test", func(t *testing.T) {
		err := rc.SaveToRDB(message)
		if err != nil {
			t.Errorf("SaveToRDB() error = %v", err)
		}
	})
}

func TestRedisCluster_Synchronize(t *testing.T) {
	rc := New(mySqlMeesageRepository.New())
	defer rc.Close()

	t.Run("Synchronize test", func(t *testing.T) {
		err := rc.Synchronize()
		if err != nil {
			t.Errorf("Synchronize() error = %v", err)
		}
	})
}
