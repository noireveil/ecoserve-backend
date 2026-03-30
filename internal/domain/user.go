package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName       string    `gorm:"type:varchar(255);not null"`
	WhatsAppNumber string    `gorm:"type:varchar(20);uniqueIndex;not null"`
	Role           string    `gorm:"type:varchar(50);not null;default:'customer'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
