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
	userRepo repository.UserRepository
}

func NewTechnicianUsecase(techRepo repository.TechnicianRepository, userRepo repository.UserRepository) TechnicianUsecase {
	return &technicianUsecase{techRepo, userRepo}
}

func (u *technicianUsecase) GetNearbyTechnicians(lon, lat float64, radiusKm int) ([]domain.Technician, error) {
	return u.techRepo.FindNearby(lon, lat, radiusKm)
}

func (u *technicianUsecase) RegisterTechnician(technician *domain.Technician, lon, lat float64) error {
	err := u.techRepo.Create(technician, lon, lat)
	if err != nil {
		return err
	}

	return u.userRepo.UpdateRole(technician.UserID.String(), "technician")
}
