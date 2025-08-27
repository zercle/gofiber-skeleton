//go:generate mockgen -source=order.go -destination=./mock/mock_order.go -package=mock
package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Status    OrderStatus `json:"status"`
	Total     float64     `json:"total"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order, orderItems []OrderItem) (Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order Order) (Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (Order, error)
}

// OrderUseCase defines the interface for order business logic
type OrderUseCase interface {
	CreateOrder(userID string, items []OrderItem) (*Order, error)
	GetOrder(id string) (*Order, error)
	GetUserOrders(userID string) ([]*Order, error)
	GetAllOrders() ([]*Order, error)
	UpdateOrderStatus(id string, status OrderStatus) error
}

var (
	ErrOrderNotFound = errors.New("order not found")
)
