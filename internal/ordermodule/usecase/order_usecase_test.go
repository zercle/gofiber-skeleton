package orderusecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/ordermodule"
	"github.com/zercle/gofiber-skeleton/internal/productmodule"
	ordermock "github.com/zercle/gofiber-skeleton/internal/ordermodule/mock"
	productmock "github.com/zercle/gofiber-skeleton/internal/productmodule/mock"
)

func TestOrderUseCase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := ordermock.NewMockOrderRepository(ctrl)
	mockProductRepo := productmock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	userID := uuid.New().String()
	productID := uuid.New().String()
	orderItems := []ordermodule.OrderItem{
		{ProductID: productID, Quantity: 2},
	}
	productPrice := 10.0

	t.Run("successful order creation", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&productmodule.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 10,
		}, nil)
		mockProductRepo.EXPECT().UpdateStock(productID, -2).Return(nil)
		mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(ordermodule.Order{
			ID:        uuid.New().String(),
			UserID:    userID,
			Status:    ordermodule.OrderStatusPending,
			Total:     productPrice * 2,
			Items:     orderItems,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

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
		order, err := usecase.CreateOrder(userID, []ordermodule.OrderItem{})
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "order must contain at least one item")
	})

	t.Run("zero quantity item", func(t *testing.T) {
		items := []ordermodule.OrderItem{{ProductID: productID, Quantity: 0}}
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
		mockProductRepo.EXPECT().GetByID(productID).Return(&productmodule.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 1,
		}, nil)

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "insufficient stock for product: Test Product")
	})

	t.Run("product stock update error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&productmodule.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 10,
		}, nil)
		mockProductRepo.EXPECT().UpdateStock(productID, -2).Return(errors.New("stock update failed"))

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "stock update failed")
	})

	t.Run("order repository creation error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(&productmodule.Product{
			ID: productID, Name: "Test Product", Price: productPrice, Stock: 10,
		}, nil)
		mockProductRepo.EXPECT().UpdateStock(productID, -2).Return(nil)
		mockOrderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(ordermodule.Order{}, errors.New("db create error"))

		order, err := usecase.CreateOrder(userID, orderItems)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "db create error")
	})
}

func TestOrderUseCase_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := ordermock.NewMockOrderRepository(ctrl)
	mockProductRepo := productmock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	orderID := uuid.New().String()
	expectedOrder := &ordermodule.Order{
		ID: orderID, UserID: uuid.New().String(), Status: ordermodule.OrderStatusPending, Total: 100.0,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	t.Run("successful order retrieval", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetOrderByID(gomock.Any(), uuid.MustParse(orderID)).Return(*expectedOrder, nil)

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
		mockOrderRepo.EXPECT().GetOrderByID(gomock.Any(), uuid.MustParse(orderID)).Return(ordermodule.Order{}, errors.New("db error"))

		order, err := usecase.GetOrder(orderID)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.EqualError(t, err, "failed to fetch order")
	})
}

func TestOrderUseCase_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := ordermock.NewMockOrderRepository(ctrl)
	mockProductRepo := productmock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	userID := uuid.New().String()
	expectedOrders := []*ordermodule.Order{
		{ID: uuid.New().String(), UserID: userID, Status: ordermodule.OrderStatusPending, Total: 50.0},
		{ID: uuid.New().String(), UserID: userID, Status: ordermodule.OrderStatusConfirmed, Total: 75.0},
	}

	t.Run("successful retrieval of user orders", func(t *testing.T) {
		// The use case now returns []ordermodule.Order, so the mock should return that.
		// Then, we convert it to []*ordermodule.Order for comparison in the test.
		mockOrderRepo.EXPECT().GetOrdersByUserID(gomock.Any(), uuid.MustParse(userID)).Return([]ordermodule.Order{
			*expectedOrders[0], *expectedOrders[1],
		}, nil)

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
		mockOrderRepo.EXPECT().GetOrdersByUserID(gomock.Any(), uuid.MustParse(userID)).Return(nil, errors.New("failed to fetch user orders"))

		orders, err := usecase.GetUserOrders(userID)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "failed to fetch user orders")
	})
}

func TestOrderUseCase_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := ordermock.NewMockOrderRepository(ctrl)
	mockProductRepo := productmock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	expectedOrders := []*ordermodule.Order{
		{ID: uuid.New().String(), UserID: uuid.New().String(), Status: ordermodule.OrderStatusPending, Total: 50.0},
		{ID: uuid.New().String(), UserID: uuid.New().String(), Status: ordermodule.OrderStatusConfirmed, Total: 75.0},
	}

	t.Run("successful retrieval of all orders", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetAllOrders(gomock.Any()).Return([]ordermodule.Order{
			*expectedOrders[0], *expectedOrders[1],
		}, nil)

		orders, err := usecase.GetAllOrders()
		require.NoError(t, err)
		assert.NotNil(t, orders)
		assert.Len(t, orders, 2)
		assert.Equal(t, expectedOrders[0].ID, orders[0].ID)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockOrderRepo.EXPECT().GetAllOrders(gomock.Any()).Return(nil, errors.New("failed to fetch all orders"))

		orders, err := usecase.GetAllOrders()
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.EqualError(t, err, "failed to fetch all orders")
	})
}

func TestOrderUseCase_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := ordermock.NewMockOrderRepository(ctrl)
	mockProductRepo := productmock.NewMockProductRepository(ctrl)
	usecase := NewOrderUseCase(mockOrderRepo, mockProductRepo)

	orderID := uuid.New().String()

	t.Run("successful status update", func(t *testing.T) {
		mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), uuid.MustParse(orderID), string(ordermodule.OrderStatusShipped)).Return(ordermodule.Order{}, nil)

		err := usecase.UpdateOrderStatus(orderID, ordermodule.OrderStatusShipped)
		require.NoError(t, err)
	})

	t.Run("empty order ID", func(t *testing.T) {
		err := usecase.UpdateOrderStatus("", ordermodule.OrderStatusShipped)
		assert.Error(t, err)
		assert.EqualError(t, err, "order ID is required")
	})

	t.Run("invalid status", func(t *testing.T) {
		err := usecase.UpdateOrderStatus(orderID, "invalid_status")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid order status")
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockOrderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), uuid.MustParse(orderID), string(ordermodule.OrderStatusShipped)).Return(ordermodule.Order{}, errors.New("failed to update order status"))

		err := usecase.UpdateOrderStatus(orderID, ordermodule.OrderStatusShipped)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to update order status")
	})
}
