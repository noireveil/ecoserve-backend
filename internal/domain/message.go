package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	SenderID   uuid.UUID      `gorm:"type:uuid;not null" json:"sender_id"`
	SenderRole string         `gorm:"type:varchar(20);not null" json:"sender_role"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type MessageRepository interface {
	Create(message *Message) error
	GetByOrderID(orderID string) ([]Message, error)
}

type MessageUsecase interface {
	SendMessage(orderID, senderID, senderRole, content string) (*Message, error)
	GetOrderMessages(orderID string) ([]Message, error)
}
