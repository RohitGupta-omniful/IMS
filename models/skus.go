package models

import (
	"time"

	"github.com/google/uuid"
)

type SKU struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	SKU       string    `gorm:"type:varchar(50);unique;not null"`
	Price     float64   `gorm:"type:numeric(10,2);not null"`
	Quantity  int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}
