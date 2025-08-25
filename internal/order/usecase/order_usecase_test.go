package orderusecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/domain/mock"
)

func TestOrderUseCase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	userID := uuid.New().String()
	productID := uuid.New().String()
	orderItems := []domain.OrderItem{
		{ProductID: productID, Quantity: 2},
	}
	productPrice := 10.0

	t.Run("successful order creation", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&domain.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 10,
		}, nil)
		mockProductRepo.EXPECT().UpdateStock(productID, -2).Return(nil)
		mockOrderRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(order *domain.Order) error {
			order.ID = uuid.New().String() // Simulate repository setting the ID
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()
			assert.Equal(t, userID, order.UserID)
			assert.Equal(t, domain.OrderStatusPending, order.Status)
			assert.Equal(t, productPrice*2, order.Total)
			assert.Len(t, order.Items, 1)
			assert.Equal(t, productID, order.Items[0].ProductID)
			assert.Equal(t, 2, order.Items[0].Quantity)
			assert.Equal(t, productPrice, order.Items[0].Price)
			return nil
		})

		order, err := usecase.CreateOrder(userID, orderItems)
		require.NoError(t, err)
		assert.NotNil(t, order)
		assert.NotEmpty(t, order.ID)
	})

	t.Run("empty user ID", func(t *testing.T) {
		order, err := usecase.CreateOrder("", orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "user ID is required")
	})

	t.Run("empty order items", func(t *testing.T) {
		order, err := usecase.CreateOrder(userID, []domain.OrderItem{})
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "order must contain at least one item")
	})

	t.Run("zero quantity item", func(t *testing.T) {
		items := []domain.OrderItem{{ProductID: productID, Quantity: 0}}
		order, err := usecase.CreateOrder(userID, items)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "item quantity must be positive")
	})

	t.Run("product not found", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(nil, errors.New("product not found"))

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "product not found")
	})

	t.Run("insufficient stock", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&domain.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 1,
		}, nil)

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "insufficient stock for product: Test Product")
	})

	t.Run("product stock update error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&domain.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 10,
		}, nil)
		mockProductRepo.EXPECT().UpdateStock(productID, -2).Return(errors.New("stock update failed"))

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "stock update failed")
	})

	t.Run("order repository creation error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&domain.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 10,
		}, nil)
		mockProductRepo.EXPECT().UpdateStock(productID, -2).Return(nil)
		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db create error"))

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "db create error")
	})
}

func TestOrderUseCase_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	orderID := uuid.New().String()
	expectedOrder := &domain.Order{
		ID: orderID, UserID: uuid.New().String(), Status: domain.OrderStatusPending, Total: 100.0,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	t.Run("successful order retrieval", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetByID(orderID).Return(expectedOrder, nil)

		order, err := usecase.GetOrder(orderID)
		require.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, expectedOrder.ID, order.ID)
	})

	t.Run("empty order ID", func(t *testing.T) {
		order, err := usecase.GetOrder("")
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "order ID is required")
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetByID(orderID).Return(nil, errors.New("db error"))

		order, err := usecase.GetOrder(orderID)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "db error")
	})
}

func TestOrderUseCase_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	userID := uuid.New().String()
	expectedOrders := []*domain.Order{
		{ID: uuid.New().String(), UserID: userID, Status: domain.OrderStatusPending, Total: 50.0},
		{ID: uuid.New().String(), UserID: userID, Status: domain.OrderStatusConfirmed, Total: 75.0},
	}

	t.Run("successful retrieval of user orders", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetByUserID(userID).Return(expectedOrders, nil)

		orders, err := usecase.GetUserOrders(userID)
		require.NoError(t, err)
		assert.NotNil(t, orders)
		assert.Len(t, orders, 2)
		assert.Equal(t, expectedOrders[0].ID, orders[0].ID)
	})

	t.Run("empty user ID", func(t *testing.T) {
		orders, err := usecase.GetUserOrders("")
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "user ID is required")
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetByUserID(userID).Return(nil, errors.New("db error"))

		orders, err := usecase.GetUserOrders(userID)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "db error")
	})
}

func TestOrderUseCase_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	expectedOrders := []*domain.Order{
		{ID: uuid.New().String(), UserID: uuid.New().String(), Status: domain.OrderStatusPending, Total: 50.0},
		{ID: uuid.New().String(), UserID: uuid.New().String(), Status: domain.OrderStatusConfirmed, Total: 75.0},
	}

	t.Run("successful retrieval of all orders", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetAll().Return(expectedOrders, nil)

		orders, err := usecase.GetAllOrders()
		require.NoError(t, err)
		assert.NotNil(t, orders)
		assert.Len(t, orders, 2)
		assert.Equal(t, expectedOrders[0].ID, orders[0].ID)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetAll().Return(nil, errors.New("db error"))

		orders, err := usecase.GetAllOrders()
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "db error")
	})
}

func TestOrderUseCase_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	orderID := uuid.New().String()

	t.Run("successful status update", func(t *testing.T) {
		mockOrderRepo.EXPECT().UpdateStatus(orderID, domain.OrderStatusShipped).Return(nil)

		err := usecase.UpdateOrderStatus(orderID, domain.OrderStatusShipped)
		require.NoError(t, err)
	})

	t.Run("empty order ID", func(t *testing.T) {
		err := usecase.UpdateOrderStatus("", domain.OrderStatusShipped)
		assert.Error(t, err)
		assert.EqualError(t, err, "order ID is required")
	})

	t.Run("invalid status", func(t *testing.T) {
		err := usecase.UpdateOrderStatus(orderID, "invalid_status")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid order status")
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockOrderRepo.EXPECT().UpdateStatus(orderID, domain.OrderStatusShipped).Return(errors.New("db error"))

		err := usecase.UpdateOrderStatus(orderID, domain.OrderStatusShipped)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}
