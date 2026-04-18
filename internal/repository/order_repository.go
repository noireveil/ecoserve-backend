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
	FindByUserID(userID string) ([]domain.Order, error)
	CompleteWithAntiFraud(id string, photoURL string, lon float64, lat float64, eWasteSaved, totalFee, platformFee, netFee float64) error
	FindIncomingOrders() ([]domain.Order, error)
	AcceptOrder(orderID string, userID string) error
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
	err := r.db.Preload("Customer").Preload("Technician.User").First(&order, "id = ?", id).Error
	return &order, err
}

func (r *orderRepository) FindByUserID(userID string) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Customer").
		Preload("Technician.User").
		Joins("LEFT JOIN technicians ON technicians.id = orders.technician_id").
		Where("orders.customer_id = ? OR technicians.user_id = ?", userID, userID).
		Order("orders.created_at desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindIncomingOrders() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Customer").
		Where("status = ?", domain.OrderStatusPending).
		Where("technician_id IS NULL").
		Order("created_at desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) AcceptOrder(orderID string, userID string) error {
	var tech domain.Technician
	if err := r.db.Where("user_id = ?", userID).First(&tech).Error; err != nil {
		return fmt.Errorf("akses ditolak: profil teknisi tidak ditemukan")
	}

	result := r.db.Model(&domain.Order{}).Where("id = ? AND status = ?", orderID, domain.OrderStatusPending).Updates(map[string]interface{}{
		"technician_id": tech.ID,
		"status":        domain.OrderStatusAccepted,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("pesanan tidak tersedia atau sudah diambil oleh teknisi lain")
	}

	return nil
}

func (r *orderRepository) CompleteWithAntiFraud(id string, photoURL string, lon float64, lat float64, eWasteSaved, totalFee, platformFee, netFee float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		point := fmt.Sprintf("SRID=4326;POINT(%f %f)", lon, lat)

		result := tx.Model(&domain.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":             domain.OrderStatusCompleted,
			"total_fee":          totalFee,
			"platform_fee":       platformFee,
			"net_technician_fee": netFee,
			"e_waste_saved_kg":   eWasteSaved,
			"photo_proof_url":    photoURL,
			"gps_lock_coord":     gorm.Expr("ST_GeomFromEWKT(?)", point),
			"is_dual_confirmed":  true,
		})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("pesanan dengan ID %s tidak ditemukan atau sudah diselesaikan", id)
		}

		orderUUID, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("format UUID pesanan tidak valid: %v", err)
		}

		impact := domain.ImpactTracker{
			OrderID:        orderUUID,
			CO2AvoidedKg:   eWasteSaved,
			EwasteDiverted: eWasteSaved / 10,
		}

		if err := tx.Create(&impact).Error; err != nil {
			return fmt.Errorf("gagal mencatat metrik lingkungan: %v", err)
		}

		return nil
	})
}
