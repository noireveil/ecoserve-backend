package repository

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status string, eWasteSaved float64) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Customer").Preload("Technician").First(&order, "id = ?", id).Error
	return &order, err
}

func (r *orderRepository) UpdateStatus(id string, status string, eWasteSaved float64) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           status,
		"e_waste_saved_kg": eWasteSaved,
	}).Error
}
