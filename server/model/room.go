package model

import "time"

type Room struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Private  bool      `json:"private"`
	CreateAt time.Time `json:"create_at"`
}
