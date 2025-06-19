package models

import (
	"time"

	"github.com/google/uuid"
)

type InventoryTransaction struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID       uuid.UUID `gorm:"type:uuid;not null"`
	HubID           uuid.UUID `gorm:"type:uuid;not null"`
	QuantityChange  int       `gorm:"not null"`
	TransactionType string    `gorm:"type:varchar(50);not null"` // e.g., restock, deduct, adjustment

	Product SKU `gorm:"foreignKey:ProductID"`
	Hub     Hub `gorm:"foreignKey:HubID"`

	CreatedAt time.Time `gorm:"default:now()"`
}
