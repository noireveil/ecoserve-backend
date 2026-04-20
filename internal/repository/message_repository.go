package repository

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetByOrderID(orderID string) ([]domain.Message, error) {
	var messages []domain.Message
	err := r.db.Where("order_id = ?", orderID).Order("created_at asc").Find(&messages).Error
	return messages, err
}
