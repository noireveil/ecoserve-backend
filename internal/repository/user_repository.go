package repository

import (
	"time"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	UpdateOTP(email, code string, expiresAt time.Time) error
	Delete(id string) error
	UpdateRole(id string, role string) error
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

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateOTP(email, code string, expiresAt time.Time) error {
	return r.db.Model(&domain.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"otp_code":       code,
		"otp_expires_at": expiresAt,
	}).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *userRepository) UpdateRole(id string, role string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("role", role).Error
}
