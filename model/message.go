package model

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	TypeChatText  MessageType = "chat_txt"
	TypeChatImage MessageType = "chat_img"
)

type Message struct {
	Id          int64       `json:"id"`
	MessageType MessageType `json:"message_type"`
	AuthorId    int64       `json:"author_id"`
	RoomId      int64       `json:"room_id"`
	Content     string      `json:"content"`
	CreateAt    time.Time   `json:"create_at"`
}

type ClientMessage struct {
	MessageType MessageType     `json:"message_type"`
	AuthorId    int64           `json:"author_id"`
	RoomId      int64           `json:"room_id"`
	Content     json.RawMessage `json:"content"`
	CreateAt    int64           `json:"create_at"`
}
