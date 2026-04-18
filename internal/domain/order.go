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
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID                 uuid.UUID               `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID         uuid.UUID               `gorm:"type:uuid;not null"`
	Customer           User                    `gorm:"foreignKey:CustomerID"`
	TechnicianID       *uuid.UUID              `gorm:"type:uuid"`
	Technician         *Technician             `gorm:"foreignKey:TechnicianID"`
	DeviceID           *uuid.UUID              `gorm:"type:uuid"`
	Device             *DigitalProductPassport `gorm:"foreignKey:DeviceID"`
	DeviceCategory     string                  `gorm:"type:varchar(100);not null"`
	ProblemDescription string                  `gorm:"type:text;not null"`
	CustomerLatitude   float64                 `gorm:"type:decimal(9,6)"`
	CustomerLongitude  float64                 `gorm:"type:decimal(9,6)"`
	Status             OrderStatus             `gorm:"type:varchar(50);not null;default:'PENDING'"`
	TotalFee           float64                 `gorm:"type:decimal(12,2);default:0.00"`
	PlatformFee        float64                 `gorm:"type:decimal(12,2);default:0.00"`
	NetTechnicianFee   float64                 `gorm:"type:decimal(12,2);default:0.00"`
	EWasteSavedKg      float64                 `gorm:"type:decimal(10,2);default:0.00"`
	PhotoProofURL      *string                 `gorm:"type:text"`
	GPSLockCoord       *string                 `gorm:"type:geometry(Point,4326)"`
	IsDualConfirmed    bool                    `gorm:"default:false"`
	IsReviewed         bool                    `gorm:"default:false" json:"is_reviewed"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
