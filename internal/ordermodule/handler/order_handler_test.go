package orderhandler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/ordermodule"
	ordermock "github.com/zercle/gofiber-skeleton/internal/ordermodule/mock"
)

// mockAuthMiddleware is a Fiber middleware that sets a user ID in c.Locals for testing.
func mockAuthMiddleware(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUseCase := ordermock.NewMockOrderUseCase(ctrl)
	app := fiber.New()
	handler := NewOrderHandler(mockOrderUseCase)

	userID := uuid.New().String()
	orderItems := []ordermodule.OrderItem{
		{ProductID: uuid.New().String(), Quantity: 1, Price: 10.0},
	}
	createOrderReq := CreateOrderRequest{Items: orderItems}
	bodyBytes, _ := json.Marshal(createOrderReq)

	// Register route with mock middleware for authenticated tests
	app.Post("/api/v1/orders/create", mockAuthMiddleware(userID), handler.CreateOrder)

	t.Run("successful order creation", func(t *testing.T) {
		expectedOrder := &ordermodule.Order{
			ID:        uuid.New().String(),
			UserID:    userID,
			Status:    ordermodule.OrderStatusPending,
			Total:     10.0,
			Items:     orderItems,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockOrderUseCase.EXPECT().CreateOrder(userID, orderItems).Return(expectedOrder, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Order created successfully", responseBody["message"])
		assert.NotNil(t, responseBody["order"])
	})

	t.Run("unauthenticated user", func(t *testing.T) {
		// Create a new app instance without the mockAuthMiddleware for this test case
		unauthApp := fiber.New()
		unauthApp.Post("/api/v1/orders/create", handler.CreateOrder)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := unauthApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "User not authenticated", responseBody["message"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader([]byte(`{"items":123}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Invalid request body", responseBody["message"])
	})

	t.Run("empty items in request", func(t *testing.T) {
		emptyItemsReq := CreateOrderRequest{Items: []ordermodule.OrderItem{}}
		emptyBodyBytes, _ := json.Marshal(emptyItemsReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader(emptyBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Order must contain at least one item", responseBody["message"])
	})

	t.Run("usecase returns error", func(t *testing.T) {
		mockOrderUseCase.EXPECT().CreateOrder(userID, orderItems).Return(nil, errors.New("insufficient stock"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "insufficient stock", responseBody["message"])
	})
}

func TestOrderHandler_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUseCase := ordermock.NewMockOrderUseCase(ctrl)
	app := fiber.New()
	handler := NewOrderHandler(mockOrderUseCase)
	app.Get("/api/v1/orders/:id", handler.GetOrder)

	orderID := uuid.New().String()
	expectedOrder := &ordermodule.Order{
		ID:        orderID,
		UserID:    uuid.New().String(),
		Status:    ordermodule.OrderStatusPending,
		Total:     100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful order retrieval", func(t *testing.T) {
		mockOrderUseCase.EXPECT().GetOrder(orderID).Return(expectedOrder, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID, nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		data := responseBody["data"].(map[string]any)
		assert.NotNil(t, data["order"])
		assert.Equal(t, expectedOrder.ID, data["order"].(map[string]any)["id"])
	})

	t.Run("order not found", func(t *testing.T) {
		mockOrderUseCase.EXPECT().GetOrder(orderID).Return(nil, errors.New("order not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID, nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Order not found", responseBody["message"])
	})
}

func TestOrderHandler_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUseCase := ordermock.NewMockOrderUseCase(ctrl)
	app := fiber.New()
	handler := NewOrderHandler(mockOrderUseCase)

	userID := uuid.New().String()
	expectedOrders := []*ordermodule.Order{
		{ID: uuid.New().String(), UserID: userID, Status: ordermodule.OrderStatusPending, Total: 50.0},
		{ID: uuid.New().String(), UserID: userID, Status: ordermodule.OrderStatusConfirmed, Total: 75.0},
	}

	// Register route with mock middleware for authenticated tests
	app.Get("/api/v1/users/:id/orders", mockAuthMiddleware(userID), handler.GetUserOrders)

	t.Run("successful retrieval of user orders", func(t *testing.T) {
		mockOrderUseCase.EXPECT().GetUserOrders(userID).Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID+"/orders", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.NotNil(t, responseBody["orders"])
		assert.Len(t, responseBody["orders"].([]any), 2)
	})

	t.Run("unauthenticated user", func(t *testing.T) {
		// Create a new app instance without the mockAuthMiddleware for this test case
		unauthApp := fiber.New()
		unauthApp.Get("/api/v1/users/:id/orders", handler.GetUserOrders)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID+"/orders", nil)

		resp, err := unauthApp.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "User not authenticated", responseBody["message"])
	})

	t.Run("usecase returns error", func(t *testing.T) {
		mockOrderUseCase.EXPECT().GetUserOrders(userID).Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID+"/orders", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "Failed to fetch orders", responseBody["message"])
	})
}

func TestOrderHandler_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUseCase := ordermock.NewMockOrderUseCase(ctrl)
	app := fiber.New()
	handler := NewOrderHandler(mockOrderUseCase)
	app.Get("/api/v1/orders", handler.GetAllOrders)

	expectedOrders := []*ordermodule.Order{
		{ID: uuid.New().String(), UserID: uuid.New().String(), Status: ordermodule.OrderStatusPending, Total: 50.0},
		{ID: uuid.New().String(), UserID: uuid.New().String(), Status: ordermodule.OrderStatusConfirmed, Total: 75.0},
	}

	t.Run("successful retrieval of all orders", func(t *testing.T) {
		mockOrderUseCase.EXPECT().GetAllOrders().Return(expectedOrders, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		data := responseBody["data"].(map[string]any)
		assert.NotNil(t, data["orders"])
		assert.Len(t, data["orders"].([]any), 2)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		mockOrderUseCase.EXPECT().GetAllOrders().Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "Failed to fetch orders", responseBody["message"])
	})
}

func TestOrderHandler_UpdateOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUseCase := ordermock.NewMockOrderUseCase(ctrl)
	app := fiber.New()
	handler := NewOrderHandler(mockOrderUseCase)
	app.Put("/api/v1/orders/:id/status", handler.UpdateOrderStatus)

	orderID := uuid.New().String()
	updateStatusReq := UpdateOrderStatusRequest{Status: string(ordermodule.OrderStatusShipped)}
	bodyBytes, _ := json.Marshal(updateStatusReq)

	t.Run("successful order status update", func(t *testing.T) {
		mockOrderUseCase.EXPECT().UpdateOrderStatus(orderID, ordermodule.OrderStatusShipped).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/"+orderID+"/status", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Order status updated successfully", responseBody["message"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/"+orderID+"/status", bytes.NewReader([]byte(`{"status":123}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Invalid request body", responseBody["message"])
	})

	t.Run("usecase returns error", func(t *testing.T) {
		mockOrderUseCase.EXPECT().UpdateOrderStatus(orderID, ordermodule.OrderStatusShipped).Return(errors.New("invalid status"))

		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/"+orderID+"/status", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "invalid status", responseBody["message"])
	})
}
