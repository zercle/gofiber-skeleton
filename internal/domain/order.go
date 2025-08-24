//go:generate mockgen -source=order.go -destination=./mock/mock_order.go -package=mock
package domain

import (
	"errors"
	"time"
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
	Create(order *Order) error
	GetByID(id string) (*Order, error)
	GetByUserID(userID string) ([]*Order, error)
	GetAll() ([]*Order, error)
	UpdateStatus(id string, status OrderStatus) error
	Update(order *Order) error
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
