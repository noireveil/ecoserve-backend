package repository

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByWhatsApp(whatsapp string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByWhatsApp(whatsapp string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("whats_app_number = ?", whatsapp).First(&user).Error
	return &user, err
}
