package repository

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device *domain.DigitalProductPassport) error
	FindByOwnerID(ownerID string) ([]domain.DigitalProductPassport, error)
	FindByID(id string) (*domain.DigitalProductPassport, error)
	Delete(id string, ownerID string) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db}
}

func (r *deviceRepository) Create(device *domain.DigitalProductPassport) error {
	return r.db.Create(device).Error
}

func (r *deviceRepository) FindByOwnerID(ownerID string) ([]domain.DigitalProductPassport, error) {
	var devices []domain.DigitalProductPassport
	err := r.db.Where("owner_id = ?", ownerID).Order("created_at desc").Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) FindByID(id string) (*domain.DigitalProductPassport, error) {
	var device domain.DigitalProductPassport
	err := r.db.Where("id = ?", id).First(&device).Error
	return &device, err
}

func (r *deviceRepository) Delete(id string, ownerID string) error {
	result := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&domain.DigitalProductPassport{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
