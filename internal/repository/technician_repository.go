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
	FindByID(id string) (*domain.Technician, error)
	GetPerformanceByUserID(userID string) (float32, int, float64, error)
	GetEarningsData(userID string) (float64, float64, int, error)
	UpdateAvailability(userID string, isAvailable bool) error
	GetAvailabilityByUserID(userID string) (bool, error)
}

type technicianRepository struct {
	db *gorm.DB
}

func NewTechnicianRepository(db *gorm.DB) TechnicianRepository {
	return &technicianRepository{db}
}

func (r *technicianRepository) Create(technician *domain.Technician, lon float64, lat float64) error {
	var existing domain.Technician
	err := r.db.Unscoped().Where("user_id = ?", technician.UserID).First(&existing).Error
	if err == nil {
		r.db.Unscoped().Delete(&existing)
	}

	if technician.ID == uuid.Nil {
		technician.ID = uuid.New()
	}

	technician.Latitude = lat
	technician.Longitude = lon
	technician.IsAvailable = true

	query := `
		INSERT INTO technicians (id, user_id, specialization, experience_years, rating, latitude, longitude, is_available, location, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ST_SetSRID(ST_MakePoint(?, ?), 4326), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	return r.db.Exec(query, technician.ID, technician.UserID, technician.Specialization, technician.ExperienceYears, technician.Rating, lat, lon, technician.IsAvailable, lon, lat).Error
}

func (r *technicianRepository) FindNearby(lon float64, lat float64, radiusKm int) ([]domain.Technician, error) {
	var technicians []domain.Technician
	radiusMeters := radiusKm * 1000

	err := r.db.Preload("User").
		Select("id, user_id, specialization, experience_years, rating, latitude, longitude, is_available, created_at, updated_at").
		Where("is_available = ?", true).
		Where("ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)", lon, lat, radiusMeters).
		Order("rating DESC").
		Find(&technicians).Error

	return technicians, err
}

func (r *technicianRepository) FindByID(id string) (*domain.Technician, error) {
	var tech domain.Technician
	err := r.db.Preload("User").Where("id = ?", id).First(&tech).Error
	return &tech, err
}

func (r *technicianRepository) UpdateAvailability(userID string, isAvailable bool) error {
	return r.db.Model(&domain.Technician{}).Where("user_id = ?", userID).Update("is_available", isAvailable).Error
}

func (r *technicianRepository) GetAvailabilityByUserID(userID string) (bool, error) {
	var isAvailable bool
	err := r.db.Model(&domain.Technician{}).
		Where("user_id = ?", userID).
		Pluck("is_available", &isAvailable).Error
	return isAvailable, err
}

func (r *technicianRepository) GetPerformanceByUserID(userID string) (float32, int, float64, error) {
	var result struct {
		Rating       float32
		TotalRepairs int
		TotalCo2     float64
	}

	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	err := r.db.Table("technicians").
		Select("technicians.rating, COUNT(orders.id) as total_repairs, COALESCE(SUM(orders.e_waste_saved_kg), 0) as total_co2").
		Joins("LEFT JOIN orders ON orders.technician_id = technicians.id AND orders.status = ? AND orders.created_at >= ?", domain.OrderStatusCompleted, threeMonthsAgo).
		Where("technicians.user_id = ?", userID).
		Group("technicians.id").
		Scan(&result).Error

	return result.Rating, result.TotalRepairs, result.TotalCo2, err
}

func (r *technicianRepository) GetEarningsData(userID string) (float64, float64, int, error) {
	var orders []domain.Order
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	err := r.db.Joins("JOIN technicians ON technicians.id = orders.technician_id").
		Where("technicians.user_id = ? AND orders.status = ? AND orders.created_at >= ?", userID, domain.OrderStatusCompleted, threeMonthsAgo).
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
