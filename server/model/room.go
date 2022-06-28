package model

import "time"

type Room struct {
	Id       int64     `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Private  bool      `json:"private" db:"private"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}
