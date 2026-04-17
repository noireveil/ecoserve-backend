package usecase

import (
	"errors"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type TechnicianPerformanceDTO struct {
	Rating       float32 `json:"rating" example:"4.85"`
	TotalRepairs int     `json:"total_repairs" example:"24"`
	CO2SavedKg   float64 `json:"co2_saved_kg" example:"150.5"`
}

type TechnicianUsecase interface {
	GetNearbyTechnicians(lon, lat float64, radiusKm int) ([]domain.Technician, error)
	RegisterTechnician(technician *domain.Technician, lon, lat float64) error
	GetPerformance(userID string) (TechnicianPerformanceDTO, error)
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

func (u *technicianUsecase) GetPerformance(userID string) (TechnicianPerformanceDTO, error) {
	rating, repairs, co2, err := u.techRepo.GetPerformanceByUserID(userID)
	if err != nil {
		return TechnicianPerformanceDTO{}, errors.New("gagal mengkalkulasi metrik performa teknisi")
	}

	return TechnicianPerformanceDTO{
		Rating:       rating,
		TotalRepairs: repairs,
		CO2SavedKg:   co2,
	}, nil
}
