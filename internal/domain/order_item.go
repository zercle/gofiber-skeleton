package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"` // Changed to uuid.UUID
	Quantity  int64     `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}