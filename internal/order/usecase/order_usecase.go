package orderusecase

import (
	"errors"

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
		product, err := uc.productRepo.GetByID(item.ProductID)
		if err != nil {
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
	order := &domain.Order{
		UserID:    userID,
		Status:    domain.OrderStatusPending,
		Total:     total,
		Items:     processedItems,
	}

	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *orderUseCase) GetOrder(id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	return uc.orderRepo.GetByID(id)
}

func (uc *orderUseCase) GetUserOrders(userID string) ([]*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return uc.orderRepo.GetByUserID(userID)
}

func (uc *orderUseCase) GetAllOrders() ([]*domain.Order, error) {
	return uc.orderRepo.GetAll()
}

func (uc *orderUseCase) UpdateOrderStatus(id string, status domain.OrderStatus) error {
	if id == "" {
		return errors.New("order ID is required")
	}

	// Validate status
	validStatuses := []domain.OrderStatus{
		domain.OrderStatusPending,
		domain.OrderStatusConfirmed,
		domain.OrderStatusShipped,
		domain.OrderStatusDelivered,
		domain.OrderStatusCancelled,
	}

	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("invalid order status")
	}

	return uc.orderRepo.UpdateStatus(id, status)
}
