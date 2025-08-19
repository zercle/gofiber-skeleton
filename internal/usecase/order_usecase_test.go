package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	order := domain.Order{
		ID:        1,
		UserID:    1,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), order).Return(order, nil).Times(1)

	createdOrder, err := orderUseCase.CreateOrder(context.Background(), order)
	assert.NoError(t, err)
	assert.Equal(t, order, createdOrder)
}

func TestCreateOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	order := domain.Order{
		ID:        1,
		UserID:    1,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
	expectedError := errors.New("database error")

	mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), order).Return(domain.Order{}, expectedError).Times(1)

	createdOrder, err := orderUseCase.CreateOrder(context.Background(), order)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.Order{}, createdOrder)
}

func TestGetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	expectedOrder := domain.Order{
		ID:        1,
		UserID:    1,
		Status:    "completed",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockOrderRepo.EXPECT().GetOrderByID(gomock.Any(), int64(1)).Return(expectedOrder, nil).Times(1)

	retrievedOrder, err := orderUseCase.GetOrderByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, retrievedOrder)
}

func TestGetOrderByIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	expectedError := errors.New("order not found")

	mockOrderRepo.EXPECT().GetOrderByID(gomock.Any(), int64(1)).Return(domain.Order{}, expectedError).Times(1)

	retrievedOrder, err := orderUseCase.GetOrderByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.Order{}, retrievedOrder)
}

func TestListOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	expectedOrders := []domain.Order{
		{ID: 1, UserID: 1, Status: "pending", CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 2, Status: "completed", CreatedAt: now, UpdatedAt: now},
	}

	mockOrderRepo.EXPECT().ListOrders(gomock.Any()).Return(expectedOrders, nil).Times(1)

	retrievedOrders, err := orderUseCase.ListOrders(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, retrievedOrders)
}

func TestListOrdersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	expectedError := errors.New("database error")

	mockOrderRepo.EXPECT().ListOrders(gomock.Any()).Return(nil, expectedError).Times(1)

	retrievedOrders, err := orderUseCase.ListOrders(context.Background())
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, retrievedOrders)
}

func TestUpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	updatedOrder := domain.Order{
		ID:        1,
		UserID:    1,
		Status:    "shipped",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), int64(1), "shipped").Return(updatedOrder, nil).Times(1)

	retrievedOrder, err := orderUseCase.UpdateOrderStatus(context.Background(), 1, "shipped")
	assert.NoError(t, err)
	assert.Equal(t, updatedOrder, retrievedOrder)
}

func TestUpdateOrderStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	expectedError := errors.New("update error")

	mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), int64(1), "shipped").Return(domain.Order{}, expectedError).Times(1)

	retrievedOrder, err := orderUseCase.UpdateOrderStatus(context.Background(), 1, "shipped")
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.Order{}, retrievedOrder)
}

func TestCreateOrderItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	orderItem := domain.OrderItem{
		ID:        1,
		OrderID:   1,
		ProductID: 1,
		Quantity:  2,
		Price:     10.0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockOrderRepo.EXPECT().CreateOrderItem(gomock.Any(), orderItem).Return(orderItem, nil).Times(1)

	createdOrderItem, err := orderUseCase.CreateOrderItem(context.Background(), orderItem)
	assert.NoError(t, err)
	assert.Equal(t, orderItem, createdOrderItem)
}

func TestCreateOrderItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	orderItem := domain.OrderItem{
		ID:        1,
		OrderID:   1,
		ProductID: 1,
		Quantity:  2,
		Price:     10.0,
		CreatedAt: now,
		UpdatedAt: now,
	}
	expectedError := errors.New("database error")

	mockOrderRepo.EXPECT().CreateOrderItem(gomock.Any(), orderItem).Return(domain.OrderItem{}, expectedError).Times(1)

	createdOrderItem, err := orderUseCase.CreateOrderItem(context.Background(), orderItem)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.OrderItem{}, createdOrderItem)
}

func TestListOrderItemsByOrderID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	now := time.Now()
	expectedOrderItems := []domain.OrderItem{
		{ID: 1, OrderID: 1, ProductID: 1, Quantity: 2, Price: 10.0, CreatedAt: now, UpdatedAt: now},
		{ID: 2, OrderID: 1, ProductID: 2, Quantity: 1, Price: 20.0, CreatedAt: now, UpdatedAt: now},
	}

	mockOrderRepo.EXPECT().ListOrderItemsByOrderID(gomock.Any(), int64(1)).Return(expectedOrderItems, nil).Times(1)

	retrievedOrderItems, err := orderUseCase.ListOrderItemsByOrderID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrderItems, retrievedOrderItems)
}

func TestListOrderItemsByOrderIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := domain.NewMockOrderRepository(ctrl)
	mockProductUC := domain.NewMockProductUseCase(ctrl)
	orderUseCase := NewOrderUsecase(mockOrderRepo, mockProductUC)

	expectedError := errors.New("database error")

	mockOrderRepo.EXPECT().ListOrderItemsByOrderID(gomock.Any(), int64(1)).Return(nil, expectedError).Times(1)

	retrievedOrderItems, err := orderUseCase.ListOrderItemsByOrderID(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, retrievedOrderItems)
}