package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/repository/db"
)

type orderRepository struct {
	db *db.Queries
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *db.Queries) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	ctx := context.Background()
	
	userUUID, err := parseUUID(order.UserID)
	if err != nil {
		return err
	}

	dbOrder, err := r.db.CreateOrder(ctx, db.CreateOrderParams{
		UserID: userUUID,
		Status: string(order.Status),
		Total:  fmt.Sprintf("%.2f", order.Total),
	})
	if err != nil {
		return err
	}

	order.ID = dbOrder.ID.String()
	order.CreatedAt = dbOrder.CreatedAt.Time
	order.UpdatedAt = dbOrder.UpdatedAt.Time
	return nil
}

func (r *orderRepository) GetByID(id string) (*domain.Order, error) {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	dbOrder, err := r.db.GetOrderByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	total, _ := strconv.ParseFloat(dbOrder.Total, 64)
	return &domain.Order{
		ID:        dbOrder.ID.String(),
		UserID:    dbOrder.UserID.String(),
		Status:    domain.OrderStatus(dbOrder.Status),
		Total:     total,
		CreatedAt: dbOrder.CreatedAt.Time,
		UpdatedAt: dbOrder.UpdatedAt.Time,
	}, nil
}

func (r *orderRepository) GetByUserID(userID string) ([]*domain.Order, error) {
	ctx := context.Background()
	
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	dbOrders, err := r.db.GetOrdersByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		total, _ := strconv.ParseFloat(dbOrder.Total, 64)
		orders[i] = &domain.Order{
			ID:        dbOrder.ID.String(),
			UserID:    dbOrder.UserID.String(),
			Status:    domain.OrderStatus(dbOrder.Status),
			Total:     total,
			CreatedAt: dbOrder.CreatedAt.Time,
			UpdatedAt: dbOrder.UpdatedAt.Time,
		}
	}

	return orders, nil
}

func (r *orderRepository) GetAll() ([]*domain.Order, error) {
	ctx := context.Background()
	
	dbOrders, err := r.db.GetAllOrders(ctx)
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		total, _ := strconv.ParseFloat(dbOrder.Total, 64)
		orders[i] = &domain.Order{
			ID:        dbOrder.ID.String(),
			UserID:    dbOrder.UserID.String(),
			Status:    domain.OrderStatus(dbOrder.Status),
			Total:     total,
			CreatedAt: dbOrder.CreatedAt.Time,
			UpdatedAt: dbOrder.UpdatedAt.Time,
		}
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     uuid,
		Status: string(status),
	})
	return err
}

func (r *orderRepository) Update(order *domain.Order) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(order.ID)
	if err != nil {
		return err
	}

	userUUID, err := parseUUID(order.UserID)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOrder(ctx, db.UpdateOrderParams{
		ID:     uuid,
		UserID: userUUID,
		Status: string(order.Status),
		Total:  fmt.Sprintf("%.2f", order.Total),
	})
	return err
}