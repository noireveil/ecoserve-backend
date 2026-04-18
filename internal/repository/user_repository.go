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
	FindUnscopedByEmail(email string) (*domain.User, error)
	RestoreAndUpdate(email, fullName string) error
	GetConsumerImpact(userID string) (int, float64, error)
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

func (r *userRepository) FindUnscopedByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Unscoped().Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) RestoreAndUpdate(email, fullName string) error {
	return r.db.Unscoped().Model(&domain.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"deleted_at": gorm.Expr("NULL"),
		"full_name":  fullName,
		"role":       "customer",
	}).Error
}

func (r *userRepository) GetConsumerImpact(userID string) (int, float64, error) {
	var result struct {
		TotalRepairs int
		TotalCo2     float64
	}

	err := r.db.Table("orders").
		Select("COUNT(id) as total_repairs, COALESCE(SUM(e_waste_saved_kg), 0) as total_co2").
		Where("customer_id = ? AND status = ?", userID, domain.OrderStatusCompleted).
		Scan(&result).Error

	return result.TotalRepairs, result.TotalCo2, err
}
