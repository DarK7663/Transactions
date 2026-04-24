package main

import (
	"time"

	"github.com/google/uuid"
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

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SenderID    uint      `gorm:"not null;index"`
	RecipientID uint      `gorm:"not null;index"`
	Amount      int64     `gorm:"not null;check:amount > 0"`
	Status      string    `gorm:"type:varchar(20);default:'completed'"` // completed, pending, failed
	Reference   string    `gorm:"uniqueIndex;not null"`                 // идемпотентность
	Description string
	CreatedAt   time.Time
}
