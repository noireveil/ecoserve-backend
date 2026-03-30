package usecase

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type TechnicianUsecase interface {
	GetNearbyTechnicians(lon, lat float64, radiusKm int) ([]domain.Technician, error)
	RegisterTechnician(technician *domain.Technician, lon, lat float64) error
}

type technicianUsecase struct {
	techRepo repository.TechnicianRepository
}

func NewTechnicianUsecase(techRepo repository.TechnicianRepository) TechnicianUsecase {
	return &technicianUsecase{techRepo}
}

func (u *technicianUsecase) GetNearbyTechnicians(lon, lat float64, radiusKm int) ([]domain.Technician, error) {
	return u.techRepo.FindNearby(lon, lat, radiusKm)
}

func (u *technicianUsecase) RegisterTechnician(technician *domain.Technician, lon, lat float64) error {
	return u.techRepo.Create(technician, lon, lat)
}
