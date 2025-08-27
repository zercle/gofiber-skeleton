package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

// OrderRepository defines the interface for order-related database operations.
type OrderRepository interface {
	CreateOrder(ctx context.Context, order domain.Order, orderItems []domain.OrderItem) (domain.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (domain.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	GetAllOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (domain.Order, error)
}

// orderRepository implements OrderRepository.
type orderRepository struct {
	q     *sqlc.Queries // The generated Queries struct (holds methods)
	rawDB *sql.DB       // The underlying DB connection (passed as DBTX)
}

// NewOrderRepository creates a new instance of OrderRepository.
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		q:     sqlc.New(), // Call the parameterless New()
		rawDB: db,         // Store the actual DB connection
	}
}

// CreateOrder creates a new order and its items in a transaction.
func (r *orderRepository) CreateOrder(ctx context.Context, order domain.Order, orderItems []domain.OrderItem) (domain.Order, error) {
	tx, err := r.rawDB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after Rollback
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil { // err is not nil, so Rollback
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			err = tx.Commit() // err is nil, so Commit
		}
	}()

	// Use r.q (the Queries struct) and pass tx as the DBTX argument
	createdOrder, err := r.q.CreateOrder(ctx, tx, sqlc.CreateOrderParams{
		UserID: uuid.MustParse(order.UserID),
		Total:  fmt.Sprintf("%.2f", order.Total),
		Status: string(order.Status),
	})
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range orderItems {
		// itemID, err := uuid.NewV7() // No longer needed as ID is not passed to CreateOrderItemParams
		// if err != nil {
		// 	return domain.Order{}, fmt.Errorf("failed to generate UUID for order item: %w", err)
		// }
		_, err = r.q.CreateOrderItem(ctx, tx, sqlc.CreateOrderItemParams{
			OrderID:   createdOrder.ID,
			ProductID: uuid.MustParse(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     fmt.Sprintf("%.2f", item.Price),
		})
		if err != nil {
			return domain.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return mapSQLCToDomainOrder(createdOrder), nil
}

// GetOrderByID retrieves an order by its ID.
func (r *orderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	order, err := r.q.GetOrderByID(ctx, r.rawDB, id) // Pass r.rawDB as DBTX
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to get order by ID: %w", err)
	}

	orderItems, err := r.q.GetOrderItemsByOrderID(ctx, r.rawDB, order.ID) // Pass r.rawDB as DBTX
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to get order items for order ID %s: %w", order.ID, err)
	}

	domainOrder := mapSQLCToDomainOrder(order)
	domainOrder.Items = mapSQLCToDomainOrderItems(orderItems)

	return domainOrder, nil
}

// GetOrdersByUserID retrieves orders by user ID.
func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	orders, err := r.q.GetOrdersByUserID(ctx, r.rawDB, userID) // Pass r.rawDB as DBTX
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}

	var domainOrders []domain.Order
	for _, order := range orders {
		orderItems, err := r.q.GetOrderItemsByOrderID(ctx, r.rawDB, order.ID) // Pass r.rawDB as DBTX
		if err != nil {
			return nil, fmt.Errorf("failed to get order items for order ID %s: %w", order.ID, err)
		}

		domainOrder := mapSQLCToDomainOrder(order)
		domainOrder.Items = mapSQLCToDomainOrderItems(orderItems)
		domainOrders = append(domainOrders, domainOrder)
	}

	return domainOrders, nil
}

// GetAllOrders retrieves all orders.
func (r *orderRepository) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := r.q.GetAllOrders(ctx, r.rawDB) // Pass r.rawDB as DBTX
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	var domainOrders []domain.Order
	for _, order := range orders {
		orderItems, err := r.q.GetOrderItemsByOrderID(ctx, r.rawDB, order.ID) // Pass r.rawDB as DBTX
		if err != nil {
			return nil, fmt.Errorf("failed to get order items for order ID %s: %w", order.ID, err)
		}

		domainOrder := mapSQLCToDomainOrder(order)
		domainOrder.Items = mapSQLCToDomainOrderItems(orderItems)
		domainOrders = append(domainOrders, domainOrder)
	}

	return domainOrders, nil
}

// UpdateOrder updates an existing order.
func (r *orderRepository) UpdateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	updatedOrder, err := r.q.UpdateOrder(ctx, r.rawDB, sqlc.UpdateOrderParams{ // Pass r.rawDB as DBTX
		ID:     uuid.MustParse(order.ID),
		UserID: uuid.MustParse(order.UserID),
		Total:  fmt.Sprintf("%.2f", order.Total),
		Status: string(order.Status),
	})
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to update order: %w", err)
	}
	return mapSQLCToDomainOrder(updatedOrder), nil
}

// UpdateOrderStatus updates the status of an order.
func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (domain.Order, error) {
	updatedOrder, err := r.q.UpdateOrderStatus(ctx, r.rawDB, sqlc.UpdateOrderStatusParams{ // Pass r.rawDB as DBTX
		ID:     id,
		Status: status,
	})
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to update order status: %w", err)
	}
	return mapSQLCToDomainOrder(updatedOrder), nil
}

func mapSQLCToDomainOrder(order sqlc.Order) domain.Order {
	total, _ := strconv.ParseFloat(order.Total, 64) // Convert string to float64
	return domain.Order{
		ID:        order.ID.String(),
		UserID:    order.UserID.String(),
		Total:     total,
		Status:    domain.OrderStatus(order.Status),
		CreatedAt: order.CreatedAt.Time,
	}
}

func mapSQLCToDomainOrderItems(sqlcOrderItems []sqlc.OrderItem) []domain.OrderItem {
	var domainOrderItems []domain.OrderItem
	for _, item := range sqlcOrderItems {
		price, _ := strconv.ParseFloat(item.Price, 64) // Convert string price to float64
		domainOrderItems = append(domainOrderItems, domain.OrderItem{
			ID:        item.ID.String(),
			OrderID:   item.OrderID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  int(item.Quantity),
			Price:     price,
		})
	}
	return domainOrderItems
}
