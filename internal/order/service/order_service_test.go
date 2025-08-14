package service_test

import (
	"testing"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/order/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockOrderRepository is a simple mock for testing
type MockOrderRepository struct {
	orders map[uuid.UUID]*domain.Order
	create func(*domain.Order) error
	getByID func(uuid.UUID) (*domain.Order, error)
	getAll func() ([]*domain.Order, error)
	getByUserID func(uuid.UUID) ([]*domain.Order, error)
	updateStatus func(uuid.UUID, domain.OrderStatus) error
	createOrderItem func(*domain.OrderItem) error
	getOrderItems func(uuid.UUID) ([]*domain.OrderItem, error)
}

func (m *MockOrderRepository) Create(order *domain.Order) error {
	if m.create != nil {
		return m.create(order)
	}
	return nil
}

func (m *MockOrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	return m.orders[id], nil
}

func (m *MockOrderRepository) GetAll() ([]*domain.Order, error) {
	if m.getAll != nil {
		return m.getAll()
	}
	orders := make([]*domain.Order, 0, len(m.orders))
	for _, o := range m.orders {
		orders = append(orders, o)
	}
	return orders, nil
}

func (m *MockOrderRepository) GetByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	if m.getByUserID != nil {
		return m.getByUserID(userID)
	}
	return nil, nil
}

func (m *MockOrderRepository) UpdateStatus(id uuid.UUID, status domain.OrderStatus) error {
	if m.updateStatus != nil {
		return m.updateStatus(id, status)
	}
	return nil
}

func (m *MockOrderRepository) CreateOrderItem(item *domain.OrderItem) error {
	if m.createOrderItem != nil {
		return m.createOrderItem(item)
	}
	return nil
}

func (m *MockOrderRepository) GetOrderItems(orderID uuid.UUID) ([]*domain.OrderItem, error) {
	if m.getOrderItems != nil {
		return m.getOrderItems(orderID)
	}
	return nil, nil
}

// MockProductRepository is a simple mock for testing
type MockProductRepository struct {
	products map[uuid.UUID]*domain.Product
	getByID func(uuid.UUID) (*domain.Product, error)
	updateStock func(uuid.UUID, int) error
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	return nil
}

func (m *MockProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	return m.products[id], nil
}

func (m *MockProductRepository) GetAll() ([]*domain.Product, error) {
	return nil, nil
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	return nil
}

func (m *MockProductRepository) Delete(id uuid.UUID) error {
	return nil
}

func (m *MockProductRepository) UpdateStock(id uuid.UUID, quantity int) error {
	if m.updateStock != nil {
		return m.updateStock(id, quantity)
	}
	return nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockProductRepo := &MockProductRepository{}
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	t.Run("empty items list", func(t *testing.T) {
		req := &domain.CreateOrderRequest{
			UserID:          uuid.New(),
			ShippingAddress: "123 Main St, City, Country",
			Items:           []*domain.OrderItemRequest{},
		}

		order, err := orderService.CreateOrder(req)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, "order must contain at least one item", err.Error())
	})

	t.Run("empty shipping address", func(t *testing.T) {
		req := &domain.CreateOrderRequest{
			UserID:          uuid.New(),
			ShippingAddress: "",
			Items: []*domain.OrderItemRequest{
				{
					ProductID: uuid.New(),
					Quantity:  1,
				},
			},
		}

		order, err := orderService.CreateOrder(req)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, "shipping address is required", err.Error())
	})
}

func TestOrderService_GetOrder(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockProductRepo := &MockProductRepository{}
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	t.Run("nil order ID", func(t *testing.T) {
		order, err := orderService.GetOrder(uuid.Nil)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, "invalid order ID", err.Error())
	})
}

func TestOrderService_GetUserOrders(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockProductRepo := &MockProductRepository{}
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	t.Run("nil user ID", func(t *testing.T) {
		orders, err := orderService.GetUserOrders(uuid.Nil)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, "invalid user ID", err.Error())
	})
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockProductRepo := &MockProductRepository{}
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	t.Run("nil order ID", func(t *testing.T) {
		req := &domain.UpdateOrderStatusRequest{
			Status: domain.OrderStatusConfirmed,
		}

		order, err := orderService.UpdateOrderStatus(uuid.Nil, req)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, "invalid order ID", err.Error())
	})
}