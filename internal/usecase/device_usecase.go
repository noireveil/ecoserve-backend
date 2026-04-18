package usecase

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type DeviceUsecase interface {
	CreateDevice(device *domain.DigitalProductPassport) error
	GetUserDevices(ownerID string) ([]domain.DigitalProductPassport, error)
	GetDeviceByID(id string) (*domain.DigitalProductPassport, error)
	DeleteDevice(deviceID string, ownerID string) error
}

type deviceUsecase struct {
	deviceRepo repository.DeviceRepository
}

func NewDeviceUsecase(deviceRepo repository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{deviceRepo}
}

func (u *deviceUsecase) CreateDevice(device *domain.DigitalProductPassport) error {
	return u.deviceRepo.Create(device)
}

func (u *deviceUsecase) GetUserDevices(ownerID string) ([]domain.DigitalProductPassport, error) {
	return u.deviceRepo.FindByOwnerID(ownerID)
}

func (u *deviceUsecase) GetDeviceByID(id string) (*domain.DigitalProductPassport, error) {
	return u.deviceRepo.FindByID(id)
}

func (u *deviceUsecase) DeleteDevice(deviceID string, ownerID string) error {
	return u.deviceRepo.Delete(deviceID, ownerID)
}
