package usecase

import (
	"context"

	"github.com/zercle/gofiber-skeleton/internal/domain"
)

// OrderUsecase defines the interface for order-related business logic.
type OrderUsecase interface {
	CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrderByID(ctx context.Context, id int64) (domain.Order, error)
	ListOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status domain.OrderStatus) (domain.Order, error)
	CreateOrderItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error)
	ListOrderItemsByOrderID(ctx context.Context, orderID int64) ([]domain.OrderItem, error)
	ProcessOrder(ctx context.Context, userID string, items []domain.OrderItem) (domain.Order, error)
}

type orderUsecase struct {
	repo      domain.OrderRepository
	productUC domain.ProductUseCase
}

// NewOrderUsecase creates a new OrderUsecase implementation.
func NewOrderUsecase(repo domain.OrderRepository, productUC domain.ProductUseCase) OrderUsecase {
	return &orderUsecase{repo: repo, productUC: productUC}
}

// CreateOrder creates a new order.
func (ou *orderUsecase) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	return ou.repo.CreateOrder(ctx, order)
}

// GetOrderByID retrieves an order by its ID.
func (ou *orderUsecase) GetOrderByID(ctx context.Context, id int64) (domain.Order, error) {
	return ou.repo.GetOrderByID(ctx, id)
}

// ListOrders retrieves all orders.
func (ou *orderUsecase) ListOrders(ctx context.Context) ([]domain.Order, error) {
	return ou.repo.ListOrders(ctx)
}

// UpdateOrderStatus updates the status of an order.
func (ou *orderUsecase) UpdateOrderStatus(ctx context.Context, id int64, status domain.OrderStatus) (domain.Order, error) {
	return ou.repo.UpdateOrderStatus(ctx, id, string(status))
}

// CreateOrderItem creates a new order item.
func (ou *orderUsecase) CreateOrderItem(ctx context.Context, item domain.OrderItem) (domain.OrderItem, error) {
	return ou.repo.CreateOrderItem(ctx, item)
}

// ListOrderItemsByOrderID retrieves order items by order ID.
func (ou *orderUsecase) ListOrderItemsByOrderID(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	return ou.repo.ListOrderItemsByOrderID(ctx, orderID)
}

// ProcessOrder processes a new order with items and adjusts product stock.
func (ou *orderUsecase) ProcessOrder(ctx context.Context, userID string, items []domain.OrderItem) (domain.Order, error) {
	order := domain.Order{
		UserID: userID,
		Status: string(domain.OrderStatusPending), // Default status
	}

	createdOrder, err := ou.repo.CreateOrder(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}

	for _, item := range items {
		item.OrderID = createdOrder.ID
		_, err := ou.repo.CreateOrderItem(ctx, item)
		if err != nil {
			// TODO: Handle rollback or compensation for already created order/items
			return domain.Order{}, err
		}

		// Assuming ProductUsecase has a method to reduce stock
		err = ou.productUC.ReduceStock(ctx, item.ProductID.String(), int(item.Quantity)) // Convert UUID to string for ProductID
		if err != nil {
			// TODO: Handle rollback or compensation for already created order/items
			return domain.Order{}, err
		}
	}

	return createdOrder, nil
}