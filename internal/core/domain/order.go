package domain

import (
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

// Order represents an order in the e-commerce system
type Order struct {
	ID              uuid.UUID   `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	Status          OrderStatus `json:"status"`
	TotalAmount     float64     `json:"total_amount"`
	ShippingAddress string      `json:"shipping_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Items           []*OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int       `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	Subtotal   float64   `json:"subtotal"`
	CreatedAt  time.Time `json:"created_at"`
	Product    *Product  `json:"product,omitempty"`
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	UserID          uuid.UUID           `json:"user_id" validate:"required"`
	Items           []*OrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingAddress string              `json:"shipping_address" validate:"required"`
}

// OrderItemRequest represents the request for an order item
type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// UpdateOrderStatusRequest represents the request to update an order status
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(order *Order) error
	GetByID(id uuid.UUID) (*Order, error)
	GetAll() ([]*Order, error)
	GetByUserID(userID uuid.UUID) ([]*Order, error)
	UpdateStatus(id uuid.UUID, status OrderStatus) error
	CreateOrderItem(item *OrderItem) error
	GetOrderItems(orderID uuid.UUID) ([]*OrderItem, error)
}

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(req *CreateOrderRequest) (*Order, error)
	GetOrder(id uuid.UUID) (*Order, error)
	GetAllOrders() ([]*Order, error)
	GetUserOrders(userID uuid.UUID) ([]*Order, error)
	UpdateOrderStatus(id uuid.UUID, req *UpdateOrderStatusRequest) (*Order, error)
}