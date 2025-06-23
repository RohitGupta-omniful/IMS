package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	HubID     uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"column:sku_id;type:uuid;not null"`
	Quantity  int       `gorm:"default:0"`

	Hub     Hub `gorm:"foreignKey:HubID"`
	Product SKU `gorm:"foreignKey:ProductID;references:ID"` // optional `references`

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

func (Inventory) TableName() string {
	return "inventory"
}
