package service

import (
	"errors"

	"gofiber-skeleton/internal/core/domain"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderService(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(req *domain.CreateOrderRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	if req.ShippingAddress == "" {
		return nil, errors.New("shipping address is required")
	}

	// Calculate total amount and validate stock
	var totalAmount float64
	var orderItems []*domain.OrderItem

	for _, itemReq := range req.Items {
		// Get product to check stock and price
		product, err := s.productRepo.GetByID(itemReq.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + itemReq.ProductID.String())
		}

		if product.Stock < itemReq.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Calculate subtotal
		subtotal := product.Price * float64(itemReq.Quantity)
		totalAmount += subtotal

		// Create order item
		orderItem := &domain.OrderItem{
			ID:        uuid.New(),
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
		}
		orderItems = append(orderItems, orderItem)

		// Update product stock
		if err := s.productRepo.UpdateStock(itemReq.ProductID, itemReq.Quantity); err != nil {
			return nil, errors.New("failed to update product stock: " + err.Error())
		}
	}

	// Create order
	order := &domain.Order{
		ID:              uuid.New(),
		UserID:          req.UserID,
		Status:          domain.OrderStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: req.ShippingAddress,
		Items:           orderItems,
	}

	// Save order to database
	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Save order items
	for _, item := range orderItems {
		item.OrderID = order.ID
		if err := s.orderRepo.CreateOrderItem(item); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (s *OrderService) GetOrder(id uuid.UUID) (*domain.Order, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid order ID")
	}

	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetAllOrders() ([]*domain.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *OrderService) GetUserOrders(userID uuid.UUID) ([]*domain.Order, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.orderRepo.GetByUserID(userID)
}

func (s *OrderService) UpdateOrderStatus(id uuid.UUID, req *domain.UpdateOrderStatusRequest) (*domain.Order, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid order ID")
	}

	// Validate status transition
	if err := s.validateStatusTransition(id, req.Status); err != nil {
		return nil, err
	}

	// Update status
	if err := s.orderRepo.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}

	// Return updated order
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) validateStatusTransition(orderID uuid.UUID, newStatus domain.OrderStatus) error {
	// Get current order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	// Simple status validation - in a real system, you might want more complex rules
	switch order.Status {
	case domain.OrderStatusPending:
		if newStatus != domain.OrderStatusConfirmed && newStatus != domain.OrderStatusCancelled {
			return errors.New("pending orders can only be confirmed or cancelled")
		}
	case domain.OrderStatusConfirmed:
		if newStatus != domain.OrderStatusShipped {
			return errors.New("confirmed orders can only be shipped")
		}
	case domain.OrderStatusShipped:
		if newStatus != domain.OrderStatusDelivered {
			return errors.New("shipped orders can only be delivered")
		}
	case domain.OrderStatusDelivered, domain.OrderStatusCancelled:
		return errors.New("delivered or cancelled orders cannot be updated")
	}

	return nil
}