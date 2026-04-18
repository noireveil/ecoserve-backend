package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type TechnicianRepository interface {
	Create(technician *domain.Technician, lon float64, lat float64) error
	FindNearby(lon float64, lat float64, radiusKm int) ([]domain.Technician, error)
	GetPerformanceByUserID(userID string) (float32, int, float64, error)
	GetEarningsData(userID string) (float64, float64, int, error)
}

type technicianRepository struct {
	db *gorm.DB
}

func NewTechnicianRepository(db *gorm.DB) TechnicianRepository {
	return &technicianRepository{db}
}

func (r *technicianRepository) Create(technician *domain.Technician, lon float64, lat float64) error {
	if technician.ID == uuid.Nil {
		technician.ID = uuid.New()
	}

	query := `
		INSERT INTO technicians (id, user_id, specialization, experience_years, rating, location, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ST_SetSRID(ST_MakePoint(?, ?), 4326), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	return r.db.Exec(query, technician.ID, technician.UserID, technician.Specialization, technician.ExperienceYears, technician.Rating, lon, lat).Error
}

func (r *technicianRepository) FindNearby(lon float64, lat float64, radiusKm int) ([]domain.Technician, error) {
	var technicians []domain.Technician

	radiusMeters := radiusKm * 1000

	err := r.db.Preload("User").
		Select("id, user_id, specialization, experience_years, rating, ST_AsText(location) as location, created_at, updated_at").
		Where("ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)", lon, lat, radiusMeters).
		Order("rating DESC").
		Find(&technicians).Error

	return technicians, err
}

func (r *technicianRepository) GetPerformanceByUserID(userID string) (float32, int, float64, error) {
	var result struct {
		Rating       float32
		TotalRepairs int
		TotalCo2     float64
	}

	err := r.db.Table("technicians").
		Select("technicians.rating, COUNT(orders.id) as total_repairs, COALESCE(SUM(orders.e_waste_saved_kg), 0) as total_co2").
		Joins("LEFT JOIN orders ON orders.technician_id = technicians.id AND orders.status = ?", domain.OrderStatusCompleted).
		Where("technicians.user_id = ?", userID).
		Group("technicians.id").
		Scan(&result).Error

	if err != nil {
		return 0, 0, 0, err
	}

	return result.Rating, result.TotalRepairs, result.TotalCo2, nil
}

func (r *technicianRepository) GetEarningsData(userID string) (float64, float64, int, error) {
	var orders []domain.Order
	err := r.db.Joins("JOIN technicians ON technicians.id = orders.technician_id").
		Where("technicians.user_id = ? AND orders.status = ?", userID, domain.OrderStatusCompleted).
		Find(&orders).Error

	if err != nil {
		return 0, 0, 0, err
	}

	var total, thisMonth float64
	completed := len(orders)
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	for _, o := range orders {
		total += o.NetTechnicianFee
		if o.UpdatedAt.Month() == currentMonth && o.UpdatedAt.Year() == currentYear {
			thisMonth += o.NetTechnicianFee
		}
	}

	return total, thisMonth, completed, nil
}
