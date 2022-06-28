package main

import (
	"chatting/model"
	"encoding/json"
	json2 "github.com/goccy/go-json"
	"testing"
	"time"
)

func BenchmarkJsonEncoding(b *testing.B) {
	message := model.ClientMessage{
		MessageType: model.TypeChatText,
		AuthorId:    123,
		RoomId:      123,
		Content:     json.RawMessage("\"asdf\""),
		CreateAt:    time.Now(),
	}

	byteMsg, err := json.Marshal(message)
	if err != nil {
		b.Error(err)
	}

	b.Run("encoding/json", func(b *testing.B) {
		var msg model.Message

		b.StartTimer()
		err = json.Unmarshal(byteMsg, &msg)
		if err != nil {
			b.Error(err)
		}

		_, err := json.Marshal(msg)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
	})

	b.ResetTimer()
	b.Run("goccy/json", func(b *testing.B) {
		var msg model.Message

		b.StartTimer()
		err = json2.Unmarshal(byteMsg, &msg)
		if err != nil {
			b.Error(err)
		}

		_, err := json2.Marshal(msg)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
	})
}
