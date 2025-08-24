package orderrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

type orderRepository struct {
	db sqlc.Querier
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db sqlc.Querier) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	ctx := context.Background()

	// Convert float64 total to string for database storage
	totalStr := fmt.Sprintf("%.2f", order.Total)

	dbOrder, err := r.db.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID: uuid.MustParse(order.UserID),
		Status: string(order.Status),
		Total:  totalStr,
	})
	if err != nil {
		return err
	}

	order.ID = dbOrder.ID.String()
	order.CreatedAt = dbOrder.CreatedAt.Time
	order.UpdatedAt = dbOrder.UpdatedAt.Time

	for _, item := range order.Items {
		priceStr := fmt.Sprintf("%.2f", item.Price)
		_, err := r.db.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID:   dbOrder.ID,
			ProductID: uuid.MustParse(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     priceStr,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *orderRepository) GetByID(id string) (*domain.Order, error) {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbOrder, err := r.db.GetOrderByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	total, err := strconv.ParseFloat(dbOrder.Total, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order total: %w", err)
	}

	dbOrderItems, err := r.db.GetOrderItemsByOrderID(ctx, dbOrder.ID)
	if err != nil {
		return nil, err
	}

	orderItems := make([]domain.OrderItem, len(dbOrderItems))
	for i, dbOrderItem := range dbOrderItems {
		itemPrice, err := strconv.ParseFloat(dbOrderItem.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order item price for item ID %s: %w", dbOrderItem.ID.String(), err)
		}
		orderItems[i] = domain.OrderItem{
			ID:        dbOrderItem.ID.String(),
			OrderID:   dbOrderItem.OrderID.String(),
			ProductID: dbOrderItem.ProductID.String(),
			Quantity:  int(dbOrderItem.Quantity),
			Price:     itemPrice,
		}
	}

	return &domain.Order{
		ID:        dbOrder.ID.String(),
		UserID:    dbOrder.UserID.String(),
		Status:    domain.OrderStatus(dbOrder.Status),
		Total:     total,
		Items:     orderItems,
		CreatedAt: dbOrder.CreatedAt.Time,
		UpdatedAt: dbOrder.UpdatedAt.Time,
	}, nil
}

func (r *orderRepository) GetByUserID(userID string) ([]*domain.Order, error) {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	dbOrders, err := r.db.GetOrdersByUserID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		total, err := strconv.ParseFloat(dbOrder.Total, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order total for order ID %s: %w", dbOrder.ID.String(), err)
		}

		dbOrderItems, err := r.db.GetOrderItemsByOrderID(ctx, dbOrder.ID)
		if err != nil {
			return nil, err
		}

		orderItems := make([]domain.OrderItem, len(dbOrderItems))
		for j, dbOrderItem := range dbOrderItems {
			itemPrice, err := strconv.ParseFloat(dbOrderItem.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse order item price for item ID %s: %w", dbOrderItem.ID.String(), err)
			}
			orderItems[j] = domain.OrderItem{
				ID:        dbOrderItem.ID.String(),
				OrderID:   dbOrderItem.OrderID.String(),
				ProductID: dbOrderItem.ProductID.String(),
				Quantity:  int(dbOrderItem.Quantity),
				Price:     itemPrice,
			}
		}

		orders[i] = &domain.Order{
			ID:        dbOrder.ID.String(),
			UserID:    dbOrder.UserID.String(),
			Status:    domain.OrderStatus(dbOrder.Status),
			Total:     total,
			Items:     orderItems,
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
		total, err := strconv.ParseFloat(dbOrder.Total, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order total for order ID %s: %w", dbOrder.ID.String(), err)
		}

		dbOrderItems, err := r.db.GetOrderItemsByOrderID(ctx, dbOrder.ID)
		if err != nil {
			return nil, err
		}

		orderItems := make([]domain.OrderItem, len(dbOrderItems))
		for j, dbOrderItem := range dbOrderItems {
			itemPrice, err := strconv.ParseFloat(dbOrderItem.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse order item price for item ID %s: %w", dbOrderItem.ID.String(), err)
			}
			orderItems[j] = domain.OrderItem{
				ID:        dbOrderItem.ID.String(),
				OrderID:   dbOrderItem.OrderID.String(),
				ProductID: dbOrderItem.ProductID.String(),
				Quantity:  int(dbOrderItem.Quantity),
				Price:     itemPrice,
			}
		}

		orders[i] = &domain.Order{
			ID:        dbOrder.ID.String(),
			UserID:    dbOrder.UserID.String(),
			Status:    domain.OrderStatus(dbOrder.Status),
			Total:     total,
			Items:     orderItems,
			CreatedAt: dbOrder.CreatedAt.Time,
			UpdatedAt: dbOrder.UpdatedAt.Time,
		}
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     parsedUUID,
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrOrderNotFound
		}
		return err
	}
	return nil
}

func (r *orderRepository) Update(order *domain.Order) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(order.ID)
	if err != nil {
		return err
	}

	totalStr := fmt.Sprintf("%.2f", order.Total)
	_, err = r.db.UpdateOrder(ctx, sqlc.UpdateOrderParams{
		ID:     parsedUUID,
		UserID: uuid.MustParse(order.UserID),
		Status: string(order.Status),
		Total:  totalStr,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrOrderNotFound
		}
		return err
	}

	// For order items, we'll delete existing and re-create.
	// In a real application, you might want to do a diff and only update/insert/delete as needed.
	// For simplicity, we'll delete all and re-insert.
	existingOrderItems, err := r.db.GetOrderItemsByOrderID(ctx, parsedUUID)
	if err != nil {
		return err
	}
	for _, item := range existingOrderItems {
		err = r.db.DeleteOrderItem(ctx, item.ID)
		if err != nil {
			return err
		}
	}

	for _, item := range order.Items {
		priceStr := fmt.Sprintf("%.2f", item.Price)
		_, err := r.db.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID:   parsedUUID,
			ProductID: uuid.MustParse(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     priceStr,
		})
		if err != nil {
			return err
		}
	}

	return nil
}