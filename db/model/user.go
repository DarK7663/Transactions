package model

import (
	"time"
)

type User struct {
	ID        int    `gorm:"primaryKey;not null;autoIncrement;unique" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Balance   int64  `json:"balance"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
