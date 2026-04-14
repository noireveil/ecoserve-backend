package repository

import (
	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type TechnicianRepository interface {
	Create(technician *domain.Technician, lon float64, lat float64) error
	FindNearby(lon float64, lat float64, radiusKm int) ([]domain.Technician, error)
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

	query := `
		SELECT id, user_id, specialization, experience_years, rating, ST_AsText(location) as location, created_at, updated_at
		FROM technicians
		WHERE ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
		ORDER BY rating DESC
	`

	err := r.db.Raw(query, lon, lat, radiusMeters).Scan(&technicians).Error
	return technicians, err
}
