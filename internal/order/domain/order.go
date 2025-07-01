package domain

import (
	"time"

	"gofiber-skeleton/internal/infra/types"

	"gorm.io/gorm"
)

// Order represents the order model.
type Order struct {
	ID        types.UUIDv7   `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    types.UUIDv7   `json:"user_id" gorm:"type:uuid;not null"`
	ProductID types.UUIDv7   `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	TotalPrice float64        `json:"total_price" gorm:"not null"`
	Status    string         `json:"status" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}