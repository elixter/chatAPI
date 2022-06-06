package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type MessageType string

const (
	TypeChatText  MessageType = "chat_txt"
	TypeChatImage MessageType = "chat_img"
)

type Message struct {
	Id          int64       `json:"id" db:"id"`
	ServerUUID  uuid.UUID   `json:"server_uuid"`
	MessageType MessageType `json:"message_type" db:"message_type"`
	AuthorId    int64       `json:"author_id" db:"author_id"`
	RoomId      int64       `json:"room_id" db:"room_id"`
	Content     string      `json:"content" db:"content"`
	CreateAt    time.Time   `json:"create_at" db:"create_at"`
}

type ClientMessage struct {
	MessageType MessageType     `json:"message_type"`
	AuthorId    int64           `json:"author_id"`
	RoomId      int64           `json:"room_id"`
	Content     json.RawMessage `json:"content"`
	CreateAt    time.Time       `json:"create_at"`
}
