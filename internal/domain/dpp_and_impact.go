package domain

import (
	"time"

	"github.com/google/uuid"
)

type DigitalProductPassport struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID      uuid.UUID `gorm:"type:uuid;not null"`
	User         User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category     string    `gorm:"type:varchar(100);not null"`
	BrandName    string    `gorm:"type:varchar(100);not null"`
	WeightInKg   float64   `gorm:"type:decimal(5,2);not null"`
	PurchaseDate time.Time `gorm:"type:date"`
	CreatedAt    time.Time
}

type AIDiagnosisLog struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DPPID           *uuid.UUID `gorm:"type:uuid"`
	ReportedSymptom string     `gorm:"type:text;not null"`
	ConfidenceScore float64    `gorm:"type:decimal(3,2);not null"`
	IsDIYEligible   bool       `gorm:"default:false"`
	CreatedAt       time.Time
}

type ImpactTracker struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null"`
	Order          Order     `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CO2AvoidedKg   float64   `gorm:"type:decimal(10,2);not null"`
	EwasteDiverted float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt      time.Time
}
