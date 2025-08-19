package domain

import (
	"context"
	"time"
)

//go:generate /home/kawin-vir/go/bin/mockgen -source=order.go -destination=mock_order_repository.go -package=domain OrderRepository

type Order struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// OrderRepository defines the interface for order data operations.
type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) (Order, error)
	GetOrderByID(ctx context.Context, id int64) (Order, error)
	ListOrders(ctx context.Context) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status string) (Order, error)
	CreateOrderItem(ctx context.Context, item OrderItem) (OrderItem, error)
	ListOrderItemsByOrderID(ctx context.Context, orderID int64) ([]OrderItem, error)
}