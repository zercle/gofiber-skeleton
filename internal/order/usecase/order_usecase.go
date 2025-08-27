package orderusecase

import (
	"context"
	"errors"
	"fmt" // Added fmt import

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type orderUseCase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) domain.OrderUseCase {
	return &orderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (uc *orderUseCase) CreateOrder(userID string, items []domain.OrderItem) (*domain.Order, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Calculate total and validate stock
	var total float64
	// Create a new slice to store processed order items with updated prices
	processedItems := make([]domain.OrderItem, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be positive")
		}

		// Get product to check stock and price
		fmt.Printf("Looking up product with ID: %s\n", item.ProductID) // Debug print
		product, err := uc.productRepo.GetByID(item.ProductID)
		if err != nil {
			fmt.Printf("Error looking up product %s: %v\n", item.ProductID, err) // Debug print
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Use product price and update the item in the new slice
		processedItems[i] = item
		processedItems[i].Price = product.Price
		total += product.Price * float64(item.Quantity)

		// Reduce stock
		if err := uc.productRepo.UpdateStock(item.ProductID, -item.Quantity); err != nil {
			return nil, err
		}
	}

	// Create order
	order := domain.Order{
		UserID: userID,
		Status: domain.OrderStatusPending,
		Total:  total,
		Items:  processedItems,
	}

	createdOrder, err := uc.orderRepo.CreateOrder(context.Background(), order, processedItems)
	if err != nil {
		return nil, err
	}

	return &createdOrder, nil
}

func (uc *orderUseCase) GetOrder(id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid order ID format")
	}

	order, err := uc.orderRepo.GetOrderByID(context.Background(), parsedID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) { // Assuming domain.ErrOrderNotFound is defined
			return nil, err
		}
		return nil, errors.New("failed to fetch order")
	}

	return &order, nil
}

func (uc *orderUseCase) GetUserOrders(userID string) ([]*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	orders, err := uc.orderRepo.GetOrdersByUserID(context.Background(), parsedUserID)
	if err != nil {
		return nil, errors.New("failed to fetch user orders")
	}

	// Convert []domain.Order to []*domain.Order
	result := make([]*domain.Order, len(orders))
	for i := range orders {
		result[i] = &orders[i]
	}
	return result, nil
}

func (uc *orderUseCase) GetAllOrders() ([]*domain.Order, error) {
	orders, err := uc.orderRepo.GetAllOrders(context.Background())
	if err != nil {
		return nil, errors.New("failed to fetch all orders")
	}

	// Convert []domain.Order to []*domain.Order
	result := make([]*domain.Order, len(orders))
	for i := range orders {
		result[i] = &orders[i]
	}
	return result, nil
}

func (uc *orderUseCase) UpdateOrderStatus(id string, status domain.OrderStatus) error {
	if id == "" {
		return errors.New("order ID is required")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid order ID format")
	}

	// Validate status
	validStatuses := map[domain.OrderStatus]struct{}{
		domain.OrderStatusPending:   {},
		domain.OrderStatusConfirmed: {},
		domain.OrderStatusShipped:   {},
		domain.OrderStatusDelivered: {},
		domain.OrderStatusCancelled: {},
	}

	if _, isValid := validStatuses[status]; !isValid {
		return errors.New("invalid order status")
	}

	_, err = uc.orderRepo.UpdateOrderStatus(context.Background(), parsedID, string(status))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return err
		}
		return errors.New("failed to update order status")
	}
	return nil
}
