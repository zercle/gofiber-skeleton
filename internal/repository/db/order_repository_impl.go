package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/zercle/gofiber-skeleton/internal/domain"
)

// orderRepo implements the OrderRepository interface using SQLC generated queries.
type orderRepo struct {
	queries *Queries
}

// NewOrderRepository creates a new OrderRepository implementation.
func NewOrderRepository(db *sql.DB) *orderRepo {
	return &orderRepo{
		queries: New(db),
	}
}

// CreateOrder creates a new order in the database.
func (r *orderRepo) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	newOrder, err := r.queries.CreateOrder(ctx, CreateOrderParams{
		UserID:     int32(parseUserID(order.UserID)), // Parse string to int32
		TotalPrice: fmt.Sprintf("%.2f", order.TotalPrice), // Convert float64 to string
		Status:     order.Status,
	})
	if err != nil {
		return domain.Order{}, err
	}
	return domain.Order{
		ID:         int64(newOrder.ID),
		UserID:     fmt.Sprintf("%d", newOrder.UserID), // Convert int32 to string
		TotalPrice: parseFloat(newOrder.TotalPrice), // Convert string to float64
		Status:     newOrder.Status,
		CreatedAt:  newOrder.CreatedAt,
		UpdatedAt:  newOrder.UpdatedAt,
	}, nil
}

// GetOrderByID retrieves an order by its ID.
func (r *orderRepo) GetOrderByID(ctx context.Context, id int64) (domain.Order, error) {
	fetchedOrder, err := r.queries.GetOrderByID(ctx, int32(id))
	if err != nil {
		return domain.Order{}, err
	}
	return domain.Order{
		ID:         int64(fetchedOrder.ID),
		UserID:     fmt.Sprintf("%d", fetchedOrder.UserID), // Convert int32 to string
		TotalPrice: parseFloat(fetchedOrder.TotalPrice), // Convert string to float64
		Status:     fetchedOrder.Status,
		CreatedAt:  fetchedOrder.CreatedAt,
		UpdatedAt:  fetchedOrder.UpdatedAt,
	}, nil
}

// ListOrders retrieves all orders.
func (r *orderRepo) ListOrders(ctx context.Context) ([]domain.Order, error) {
	fetchedOrders, err := r.queries.ListOrders(ctx)
	if err != nil {
		return nil, err
	}
	orders := make([]domain.Order, len(fetchedOrders))
	for i, o := range fetchedOrders {
		orders[i] = domain.Order{
			ID:         int64(o.ID),
			UserID:     fmt.Sprintf("%d", o.UserID), // Convert int32 to string
			TotalPrice: parseFloat(o.TotalPrice), // Convert string to float64
			Status:     o.Status,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		}
	}
	return orders, nil
}

// UpdateOrderStatus updates the status of an order.
func (r *orderRepo) UpdateOrderStatus(ctx context.Context, id int64, status string) (domain.Order, error) {
	updatedOrder, err := r.queries.UpdateOrderStatus(ctx, UpdateOrderStatusParams{
		ID:     int32(id), // Convert int64 to int32
		Status: status,
	})
	if err != nil {
		return domain.Order{}, err
	}
	return domain.Order{
		ID:         int64(updatedOrder.ID),
		UserID:     fmt.Sprintf("%d", updatedOrder.UserID), // Convert int32 to string
		TotalPrice: parseFloat(updatedOrder.TotalPrice), // Convert string to float64
		Status:     updatedOrder.Status,
		CreatedAt:  updatedOrder.CreatedAt,
		UpdatedAt:  updatedOrder.UpdatedAt,
	}, nil
}

// CreateOrderItem creates a new order item in the database.
func (r *orderRepo) CreateOrderItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	newOrderItem, err := r.queries.CreateOrderItem(ctx, CreateOrderItemParams{
		OrderID:   int32(item.OrderID),
		ProductID: item.ProductID, // Now uuid.UUID
		Quantity:  int32(item.Quantity),
		Price:     fmt.Sprintf("%.2f", item.Price), // Convert float64 to string
	})
	if err != nil {
		return domain.OrderItem{}, err
	}
	return domain.OrderItem{
		ID:        int64(newOrderItem.ID),
		OrderID:   int64(newOrderItem.OrderID),
		ProductID: newOrderItem.ProductID, // Now uuid.UUID
		Quantity:  int64(newOrderItem.Quantity),
		Price:     parseFloat(newOrderItem.Price), // Convert string to float64
		CreatedAt: newOrderItem.CreatedAt,
		UpdatedAt: newOrderItem.UpdatedAt,
	}, nil
}

// ListOrderItemsByOrderID retrieves order items by order ID.
func (r *orderRepo) ListOrderItemsByOrderID(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	fetchedItems, err := r.queries.ListOrderItemsByOrderID(ctx, int32(orderID))
	if err != nil {
		return nil, err
	}
	items := make([]domain.OrderItem, len(fetchedItems))
	for i, item := range fetchedItems {
		items[i] = domain.OrderItem{
			ID:        int64(item.ID),
			OrderID:   int64(item.OrderID),
			ProductID: item.ProductID, // Now uuid.UUID
			Quantity:  int64(item.Quantity),
			Price:     parseFloat(item.Price), // Convert string to float64
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}
	return items, nil
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
		return 0.0
	}
	return f
}

func parseUserID(s string) int64 {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		return 0
	}
	return id
}