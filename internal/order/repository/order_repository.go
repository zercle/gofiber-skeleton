package repository

import (
	"context"
	"database/sql"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/platform/db"

	"github.com/google/uuid"
)

type OrderRepository struct {
	queries *db.Queries
}

func NewOrderRepository(queries *db.Queries) *OrderRepository {
	return &OrderRepository{
		queries: queries,
	}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	ctx := context.Background()
	
	dbOrder, err := r.queries.CreateOrder(ctx, db.CreateOrderParams{
		UserID:          order.UserID,
		TotalAmount:     order.TotalAmount,
		ShippingAddress: order.ShippingAddress,
	})
	if err != nil {
		return err
	}

	// Update the order with the generated ID and timestamps
	order.ID = dbOrder.ID
	order.Status = domain.OrderStatus(dbOrder.Status)
	order.CreatedAt = dbOrder.CreatedAt
	order.UpdatedAt = dbOrder.UpdatedAt
	return nil
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	ctx := context.Background()
	
	dbOrder, err := r.queries.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		ID:              dbOrder.ID,
		UserID:          dbOrder.UserID,
		Status:          domain.OrderStatus(dbOrder.Status),
		TotalAmount:     float64(dbOrder.TotalAmount),
		ShippingAddress: dbOrder.ShippingAddress,
		CreatedAt:       dbOrder.CreatedAt,
		UpdatedAt:       dbOrder.UpdatedAt,
	}

	// Get order items
	items, err := r.GetOrderItems(id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r *OrderRepository) GetAll() ([]*domain.Order, error) {
	ctx := context.Background()
	
	dbOrders, err := r.queries.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		orders[i] = &domain.Order{
			ID:              dbOrder.ID,
			UserID:          dbOrder.UserID,
			Status:          domain.OrderStatus(dbOrder.Status),
			TotalAmount:     float64(dbOrder.TotalAmount),
			ShippingAddress: dbOrder.ShippingAddress,
			CreatedAt:       dbOrder.CreatedAt,
			UpdatedAt:       dbOrder.UpdatedAt,
		}
	}

	return orders, nil
}

func (r *OrderRepository) GetByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	ctx := context.Background()
	
	dbOrders, err := r.queries.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		orders[i] = &domain.Order{
			ID:              dbOrder.ID,
			UserID:          dbOrder.UserID,
			Status:          domain.OrderStatus(dbOrder.Status),
			TotalAmount:     float64(dbOrder.TotalAmount),
			ShippingAddress: dbOrder.ShippingAddress,
			CreatedAt:       dbOrder.CreatedAt,
			UpdatedAt:       dbOrder.UpdatedAt,
		}
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(id uuid.UUID, status domain.OrderStatus) error {
	ctx := context.Background()
	
	_, err := r.queries.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     id,
		Status: string(status),
	})
	return err
}

func (r *OrderRepository) CreateOrderItem(item *domain.OrderItem) error {
	ctx := context.Background()
	
	dbItem, err := r.queries.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   item.OrderID,
		ProductID: item.ProductID,
		Quantity:  int32(item.Quantity),
		UnitPrice: item.UnitPrice,
		Subtotal:  item.Subtotal,
	})
	if err != nil {
		return err
	}

	// Update the item with the generated ID and timestamp
	item.ID = dbItem.ID
	item.CreatedAt = dbItem.CreatedAt
	return nil
}

func (r *OrderRepository) GetOrderItems(orderID uuid.UUID) ([]*domain.OrderItem, error) {
	ctx := context.Background()
	
	dbItems, err := r.queries.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.OrderItem, len(dbItems))
	for i, dbItem := range dbItems {
		items[i] = &domain.OrderItem{
			ID:        dbItem.ID,
			OrderID:   dbItem.OrderID,
			ProductID: dbItem.ProductID,
			Quantity:  int(dbItem.Quantity),
			UnitPrice: float64(dbItem.UnitPrice),
			Subtotal:  float64(dbItem.Subtotal),
			CreatedAt: dbItem.CreatedAt,
		}
	}

	return items, nil
}