package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	CompleteWithAntiFraud(id string, photoURL string, lon float64, lat float64, eWasteSaved float64) error
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

func (r *orderRepository) CompleteWithAntiFraud(id string, photoURL string, lon float64, lat float64, eWasteSaved float64) error {
	point := fmt.Sprintf("SRID=4326;POINT(%f %f)", lon, lat)

	result := r.db.Model(&domain.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":            "COMPLETED",
		"e_waste_saved_kg":  eWasteSaved,
		"photo_proof_url":   photoURL,
		"gps_lock_coord":    gorm.Expr("ST_GeomFromEWKT(?)", point),
		"is_dual_confirmed": true,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("pesanan dengan ID %s tidak ditemukan atau sudah diselesaikan", id)
	}

	orderUUID, parseErr := uuid.Parse(id)
	if parseErr == nil {
		impact := domain.ImpactTracker{
			OrderID:        orderUUID,
			CO2AvoidedKg:   eWasteSaved,
			EwasteDiverted: eWasteSaved / 10,
		}
		r.db.Create(&impact)
	}

	return nil
}
