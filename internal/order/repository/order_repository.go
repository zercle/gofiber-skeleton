//go:generate mockgen -source=order_repository.go -destination=mocks/mock_order_repository.go -package=mocks

package repository

import (
	"context"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	
	// Convert UserID to pgtype.UUID
	userID := pgtype.UUID{}
	userID.Scan(order.UserID)

	// Convert TotalAmount to pgtype.Numeric
	totalAmount := pgtype.Numeric{}
	if err := totalAmount.Scan(order.TotalAmount); err != nil {
		return err
	}
	
	dbOrder, err := r.queries.CreateOrder(ctx, db.CreateOrderParams{
		UserID:          userID,
		TotalAmount:     totalAmount,
		ShippingAddress: order.ShippingAddress,
	})
	if err != nil {
		return err
	}

	// Update the order with the generated ID and timestamps
	order.ID = dbOrder.ID.Bytes
	order.Status = domain.OrderStatus(dbOrder.Status)
	order.CreatedAt = dbOrder.CreatedAt.Time
	order.UpdatedAt = dbOrder.UpdatedAt.Time
	return nil
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)
	
	dbOrder, err := r.queries.GetOrder(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	// Convert TotalAmount from pgtype.Numeric to float64
	var totalAmount float64
	if dbOrder.TotalAmount.Valid {
		if err := dbOrder.TotalAmount.Scan(&totalAmount); err != nil {
			return nil, err
		}
	}

	order := &domain.Order{
		ID:              dbOrder.ID.Bytes,
		UserID:          dbOrder.UserID.Bytes,
		Status:          domain.OrderStatus(dbOrder.Status),
		TotalAmount:     totalAmount,
		ShippingAddress: dbOrder.ShippingAddress,
		CreatedAt:       dbOrder.CreatedAt.Time,
		UpdatedAt:       dbOrder.UpdatedAt.Time,
	}

	// Get order items
	items, err := r.GetOrderItems(order.ID)
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
		// Convert TotalAmount from pgtype.Numeric to float64
		var totalAmount float64
		if dbOrder.TotalAmount.Valid {
			if err := dbOrder.TotalAmount.Scan(&totalAmount); err != nil {
				return nil, err
			}
		}

		orders[i] = &domain.Order{
			ID:              dbOrder.ID.Bytes,
			UserID:          dbOrder.UserID.Bytes,
			Status:          domain.OrderStatus(dbOrder.Status),
			TotalAmount:     totalAmount,
			ShippingAddress: dbOrder.ShippingAddress,
			CreatedAt:       dbOrder.CreatedAt.Time,
			UpdatedAt:       dbOrder.UpdatedAt.Time,
		}
	}

	return orders, nil
}

func (r *OrderRepository) GetByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(userID)
	
	dbOrders, err := r.queries.GetUserOrders(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		// Convert TotalAmount from pgtype.Numeric to float64
		var totalAmount float64
		if dbOrder.TotalAmount.Valid {
			if err := dbOrder.TotalAmount.Scan(&totalAmount); err != nil {
				return nil, err
			}
		}

		orders[i] = &domain.Order{
			ID:              dbOrder.ID.Bytes,
			UserID:          dbOrder.UserID.Bytes,
			Status:          domain.OrderStatus(dbOrder.Status),
			TotalAmount:     totalAmount,
			ShippingAddress: dbOrder.ShippingAddress,
			CreatedAt:       dbOrder.CreatedAt.Time,
			UpdatedAt:       dbOrder.UpdatedAt.Time,
		}
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(id uuid.UUID, status domain.OrderStatus) error {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)
	
	_, err := r.queries.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     pgUUID,
		Status: string(status),
	})
	return err
}

func (r *OrderRepository) CreateOrderItem(item *domain.OrderItem) error {
	ctx := context.Background()
	
	// Convert OrderID to pgtype.UUID
	orderID := pgtype.UUID{}
	orderID.Scan(item.OrderID)

	// Convert ProductID to pgtype.UUID
	productID := pgtype.UUID{}
	productID.Scan(item.ProductID)

	// Convert UnitPrice to pgtype.Numeric
	unitPrice := pgtype.Numeric{}
	if err := unitPrice.Scan(item.UnitPrice); err != nil {
		return err
	}

	// Convert Subtotal to pgtype.Numeric
	subtotal := pgtype.Numeric{}
	if err := subtotal.Scan(item.Subtotal); err != nil {
		return err
	}
	
	dbItem, err := r.queries.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  int32(item.Quantity),
		UnitPrice: unitPrice,
		Subtotal:  subtotal,
	})
	if err != nil {
		return err
	}

	// Update the item with the generated ID and timestamp
	item.ID = dbItem.ID.Bytes
	item.CreatedAt = dbItem.CreatedAt.Time
	return nil
}

func (r *OrderRepository) GetOrderItems(orderID uuid.UUID) ([]*domain.OrderItem, error) {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(orderID)
	
	dbItems, err := r.queries.GetOrderItems(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.OrderItem, len(dbItems))
	for i, dbItem := range dbItems {
		// Convert UnitPrice from pgtype.Numeric to float64
		var unitPrice float64
		if dbItem.UnitPrice.Valid {
			if err := dbItem.UnitPrice.Scan(&unitPrice); err != nil {
				return nil, err
			}
		}

		// Convert Subtotal from pgtype.Numeric to float64
		var subtotal float64
		if dbItem.Subtotal.Valid {
			if err := dbItem.Subtotal.Scan(&subtotal); err != nil {
				return nil, err
			}
		}

		items[i] = &domain.OrderItem{
			ID:        dbItem.ID.Bytes,
			OrderID:   dbItem.OrderID.Bytes,
			ProductID: dbItem.ProductID.Bytes,
			Quantity:  int(dbItem.Quantity),
			UnitPrice: unitPrice,
			Subtotal:  subtotal,
			CreatedAt: dbItem.CreatedAt.Time,
		}
	}

	return items, nil
}