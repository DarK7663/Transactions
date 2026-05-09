package model

import (
	"time"

	"github.com/google/uuid"
)

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
