package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	orderhandler "github.com/zercle/gofiber-skeleton/internal/order/handler"
	orderrepository "github.com/zercle/gofiber-skeleton/internal/order/repository"
	orderusecase "github.com/zercle/gofiber-skeleton/internal/order/usecase"
	productrepository "github.com/zercle/gofiber-skeleton/internal/product/repository"
)

func setupOrderIntegrationTest(t *testing.T) (*fiber.App, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	orderRepo := orderrepository.NewOrderRepository(db)
	productRepo := productrepository.NewProductRepository(db)
	orderUseCase := orderusecase.NewOrderUseCase(orderRepo, productRepo)
	orderHandler := orderhandler.NewOrderHandler(orderUseCase)

	app := fiber.New()
	// Public routes (no auth)
	app.Get("/api/v1/orders", orderHandler.GetAllOrders)
	app.Get("/api/v1/orders/:id", orderHandler.GetOrder)

	// Authenticated routes group
	authRoutes := app.Group("/api/v1/orders")
	authRoutes.Use(func(c *fiber.Ctx) error {
		// This is a mock JWT middleware for integration tests
		// In a real app, this would validate a JWT token and extract user_id
		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "User not authenticated",
			})
		}
		// For simplicity in testing, assume the token directly contains the userID
		userID := authHeader[7:]
		c.Locals("user_id", userID)
		return c.Next()
	})
	authRoutes.Put("/:id/status", orderHandler.UpdateOrderStatus)
	authRoutes.Post("/create", orderHandler.CreateOrder)

	return app, mock, db
}

func TestOrderIntegration_GetAllOrders(t *testing.T) {
	app, mock, db := setupOrderIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	order1ID := uuid.New()
	order2ID := uuid.New()
	mockTime := time.Now()

	t.Run("successful end-to-end get all orders", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "status", "total", "created_at", "updated_at"}).
			AddRow(order1ID, uuid.New(), domain.OrderStatusPending, "100.00", mockTime, mockTime).
			AddRow(order2ID, uuid.New(), domain.OrderStatusConfirmed, "200.00", mockTime, mockTime)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, user_id, status, total, created_at, updated_at FROM orders ORDER BY created_at DESC`,
		)).WillReturnRows(rows)

		// Mock order items for order1ID
		order1ItemsRows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price"}).
			AddRow(uuid.New(), order1ID, uuid.New(), 1, "50.00")
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`,
		)).WithArgs(order1ID).WillReturnRows(order1ItemsRows)

		// Mock order items for order2ID
		order2ItemsRows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price"}).
			AddRow(uuid.New(), order2ID, uuid.New(), 2, "75.00")
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`,
		)).WithArgs(order2ID).WillReturnRows(order2ItemsRows)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)["orders"].([]any)
		assert.Len(t, data, 2)
		assert.Equal(t, order1ID.String(), data[0].(map[string]any)["id"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderIntegration_GetOrder(t *testing.T) {
	app, mock, db := setupOrderIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	orderID := uuid.New()
	expectedOrder := domain.Order{
		ID:        orderID.String(),
		UserID:    uuid.New().String(),
		Status:    domain.OrderStatusPending,
		Total:     99.99,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful end-to-end order retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "status", "total", "created_at", "updated_at"}).
			AddRow(expectedOrder.ID, expectedOrder.UserID, expectedOrder.Status, fmt.Sprintf("%.2f", expectedOrder.Total), expectedOrder.CreatedAt, expectedOrder.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, user_id, status, total, created_at, updated_at FROM orders WHERE id = $1`,
		)).
			WithArgs(orderID).
			WillReturnRows(rows)

		// Mock order items for the expected order
		orderItemsRows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price"}).
			AddRow(uuid.New(), expectedOrder.ID, uuid.New(), 1, "50.00")
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`,
		)).WithArgs(orderID).WillReturnRows(orderItemsRows)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)["order"].(map[string]any)
		assert.Equal(t, expectedOrder.ID, data["id"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("order not found end-to-end", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, user_id, status, total, created_at, updated_at FROM orders WHERE id = $1`,
		)).
			WithArgs(orderID).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Order not found", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderIntegration_UpdateOrderStatus(t *testing.T) {
	app, mock, db := setupOrderIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	userID := uuid.New().String() // Add userID for authentication
	orderID := uuid.New()
	updateStatusInput := orderhandler.UpdateOrderStatusRequest{
		Status: string(domain.OrderStatusShipped),
	}
	mockTime := time.Now()

	t.Run("successful end-to-end order status update", func(t *testing.T) {

		// Expect UpdateStatus call
		updatedOrderRows := sqlmock.NewRows([]string{"id", "user_id", "status", "total", "created_at", "updated_at"}).
			AddRow(orderID, uuid.New(), domain.OrderStatusShipped, "100.00", mockTime, mockTime) // Assuming mockTime is still available

		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING id, user_id, status, total, created_at, updated_at`,
		)).
			WithArgs(orderID, domain.OrderStatusShipped).
			WillReturnRows(updatedOrderRows)

		body, _ := json.Marshal(updateStatusInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/"+orderID.String()+"/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID)) // Add Authorization header

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Order status updated successfully", responseBody["message"]) // No "status" for success on this endpoint
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("order not found end-to-end status update", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING id, user_id, status, total, created_at, updated_at`,
		)).
			WithArgs(orderID, domain.OrderStatusShipped).
			WillReturnError(sql.ErrNoRows) // Simulate order not found for update

		body, _ := json.Marshal(updateStatusInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/"+orderID.String()+"/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID)) // Add Authorization header

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode) // UseCase returns bad request on not found due to repository error

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "failed to update order status", responseBody["message"]) // Error message from repository
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid status end-to-end status update", func(t *testing.T) {
		// No GetByID mock expected here as the use case directly validates status
		invalidStatusInput := orderhandler.UpdateOrderStatusRequest{
			Status: "INVALID_STATUS",
		}

		body, _ := json.Marshal(invalidStatusInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/"+orderID.String()+"/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID)) // Add Authorization header

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "invalid order status", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderIntegration_CreateOrder(t *testing.T) {
	userID := uuid.New().String()
	app, mock, db := setupOrderIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	productID := uuid.New()
	mockTime := time.Now()

	createOrderInput := orderhandler.CreateOrderRequest{
		Items: []domain.OrderItem{
			{ProductID: productID.String(), Quantity: 2},
		},
	}

	t.Run("successful end-to-end order creation", func(t *testing.T) {

		// Mock GetByID for product validation and stock check
		productRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, "Test Product", "Desc", "50.00", 10, "url", mockTime, mockTime)
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnRows(productRows).
			RowsWillBeClosed()

		// Mock UpdateStock for product (happens before transaction in usecase)
		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE products SET stock = stock + $2, updated_at = NOW() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productID, -int32(createOrderInput.Items[0].Quantity)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
				AddRow(productID, "Test Product", "Desc", "50.00", 8, "url", mockTime, mockTime))

		mock.ExpectBegin() // Expect a transaction to begin

		// Mock Create order (happens within transaction)
		orderRows := sqlmock.NewRows([]string{"id", "user_id", "status", "total", "created_at", "updated_at"}).
			AddRow(uuid.New(), uuid.MustParse(userID), domain.OrderStatusPending, "100.00", mockTime, mockTime)
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO orders (user_id, status, total) VALUES ($1, $2, $3) RETURNING id, user_id, status, total, created_at, updated_at`,
		)).
			WithArgs(uuid.MustParse(userID), domain.OrderStatusPending, "100.00").
			WillReturnRows(orderRows)

		// Mock Create order items (happens within transaction)
		orderItemRows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price"}).
			AddRow(uuid.New(), uuid.New(), productID, int32(2), "50.00")
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id, order_id, product_id, quantity, price`,
		)).
			WithArgs(sqlmock.AnyArg(), productID, int32(2), "50.00").
			WillReturnRows(orderItemRows)

		mock.ExpectCommit() // Expect the transaction to commit

		body, _ := json.Marshal(createOrderInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		// Auth middleware is now applied globally in setupOrderIntegrationTest for authenticated routes
		// Ensure user_id is set for the mock auth middleware
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID)) // Simulate JWT token

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Order created successfully", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insufficient stock", func(t *testing.T) {
		// Mock GetByID for product validation and stock check (insufficient stock)
		productRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, "Test Product", "Desc", "50.00", 1, "url", mockTime, mockTime) // Stock of 1
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnRows(productRows)


		body, _ := json.Marshal(createOrderInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID)) // Simulate JWT token

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "insufficient stock for product: Test Product", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
