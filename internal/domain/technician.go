package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Technician struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Specialization  string    `gorm:"type:varchar(100);not null"`
	ExperienceYears int       `gorm:"not null;default:0"`
	Rating          float32   `gorm:"type:decimal(3,2);default:0.00"`
	Location        string    `gorm:"type:geometry(Point,4326);not null" json:"-"`
	Latitude        float64   `gorm:"type:decimal(9,6);not null;default:0" json:"latitude"`
	Longitude       float64   `gorm:"type:decimal(9,6);not null;default:0" json:"longitude"`
	IsAvailable     bool      `gorm:"not null;default:true" json:"is_available"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
