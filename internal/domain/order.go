package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusAccepted  OrderStatus = "ACCEPTED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
)

type Order struct {
	ID                 uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID         uuid.UUID   `gorm:"type:uuid;not null"`
	Customer           User        `gorm:"foreignKey:CustomerID"`
	TechnicianID       *uuid.UUID  `gorm:"type:uuid"`
	Technician         *Technician `gorm:"foreignKey:TechnicianID"`
	DeviceCategory     string      `gorm:"type:varchar(100);not null"`
	ProblemDescription string      `gorm:"type:text;not null"`
	CustomerLatitude   float64     `gorm:"type:decimal(9,6)"`
	CustomerLongitude  float64     `gorm:"type:decimal(9,6)"`
	Status             OrderStatus `gorm:"type:varchar(50);not null;default:'PENDING'"`
	TotalFee           float64     `gorm:"type:decimal(12,2);default:0.00"`
	PlatformFee        float64     `gorm:"type:decimal(12,2);default:0.00"`
	NetTechnicianFee   float64     `gorm:"type:decimal(12,2);default:0.00"`
	EWasteSavedKg      float64     `gorm:"type:decimal(10,2);default:0.00"`
	PhotoProofURL      *string     `gorm:"type:text"`
	GPSLockCoord       *string     `gorm:"type:geometry(Point,4326)"`
	IsDualConfirmed    bool        `gorm:"default:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
