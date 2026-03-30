package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                 uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID         uuid.UUID   `gorm:"type:uuid;not null"`
	Customer           User        `gorm:"foreignKey:CustomerID"`
	TechnicianID       *uuid.UUID  `gorm:"type:uuid"`
	Technician         *Technician `gorm:"foreignKey:TechnicianID"`
	DeviceCategory     string      `gorm:"type:varchar(100);not null"`
	ProblemDescription string      `gorm:"type:text;not null"`
	Status             string      `gorm:"type:varchar(50);not null;default:'PENDING'"`
	EWasteSavedKg      float64     `gorm:"type:decimal(10,2);default:0.00"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
