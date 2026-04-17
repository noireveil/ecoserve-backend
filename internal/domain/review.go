package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID      uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null"`
	Order        Order      `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID   uuid.UUID  `gorm:"type:uuid;not null"`
	User         User       `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TechnicianID uuid.UUID  `gorm:"type:uuid;not null"`
	Technician   Technician `gorm:"foreignKey:TechnicianID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rating       int        `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment      string     `gorm:"type:text"`
	CreatedAt    time.Time
}
