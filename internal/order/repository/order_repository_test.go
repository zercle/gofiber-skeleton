package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

func TestOrderRepository_CreateOrder(t *testing.T) {
	productID := uuid.New()
	userID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name         string
		order        domain.Order
		orderItems   []domain.OrderItem
		mockSetup    func(sqlmock.Sqlmock)
		expectedErr  error
		expectedFunc func(domain.Order, error)
	}{
		{
			name: "successful order creation",
			order: domain.Order{
				ID:     orderID.String(),
				UserID: userID.String(),
				Total:  150.00,
				Status: domain.OrderStatusPending,
			},
			orderItems: []domain.OrderItem{
				{ProductID: productID.String(), Quantity: 1, Price: 150.00},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO orders (user_id, status, total) VALUES ($1, $2, $3) RETURNING id, user_id, status, total, created_at, updated_at`,
				)).
					WithArgs(userID, "pending", "150.00").
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "total", "created_at", "updated_at"}).
						AddRow(orderID, userID, "pending", "150.00", time.Now(), time.Now()))

				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id, order_id, product_id, quantity, price`,
				)).
					WithArgs(orderID, productID, int32(1), "150.00").
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "price"}).
						AddRow(uuid.New(), orderID, productID, 1, "150.00"))
				mock.ExpectCommit()
			},
			expectedErr: nil,
			expectedFunc: func(order domain.Order, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, orderID.String(), order.ID)
				assert.Equal(t, userID.String(), order.UserID)
				assert.Equal(t, 150.00, order.Total)
				assert.Equal(t, domain.OrderStatusPending, order.Status)
			},
		},
		{
			name: "begin transaction error",
			order: domain.Order{
				ID:     orderID.String(),
				UserID: userID.String(),
				Total:  150.00,
				Status: domain.OrderStatusPending,
			},
			orderItems: []domain.OrderItem{
				{ProductID: productID.String(), Quantity: 1, Price: 150.00},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			expectedErr: errors.New("failed to begin transaction: begin error"),
			expectedFunc: func(order domain.Order, err error) {
				assert.Error(t, err)
				assert.Empty(t, order)
				assert.Contains(t, err.Error(), "begin error")
			},
		},
		{
			name: "create order query error",
			order: domain.Order{
				ID:     orderID.String(),
				UserID: userID.String(),
				Total:  150.00,
				Status: domain.OrderStatusPending,
			},
			orderItems: []domain.OrderItem{
				{ProductID: productID.String(), Quantity: 1, Price: 150.00},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO orders (user_id, status, total) VALUES ($1, $2, $3) RETURNING id, user_id, status, total, created_at, updated_at`,
				)).
					WithArgs(userID, "pending", "150.00").
					WillReturnError(errors.New("create order error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("failed to create order: create order error"),
			expectedFunc: func(order domain.Order, err error) {
				assert.Error(t, err)
				assert.Empty(t, order)
				assert.Contains(t, err.Error(), "create order error")
			},
		},
		{
			name: "create order item query error",
			order: domain.Order{
				ID:     orderID.String(),
				UserID: userID.String(),
				Total:  150.00,
				Status: domain.OrderStatusPending,
			},
			orderItems: []domain.OrderItem{
				{ProductID: productID.String(), Quantity: 1, Price: 150.00},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO orders (user_id, status, total) VALUES ($1, $2, $3) RETURNING id, user_id, status, total, created_at, updated_at`,
				)).
					WithArgs(userID, "pending", "150.00").
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "total", "created_at", "updated_at"}).
						AddRow(orderID, userID, "pending", "150.00", time.Now(), time.Now()))

				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id, order_id, product_id, quantity, price`,
				)).
					WithArgs(orderID, productID, int32(1), "150.00").
					WillReturnError(errors.New("create order item error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("failed to create order item: create order item error"),
			expectedFunc: func(order domain.Order, err error) {
				assert.Error(t, err)
				assert.Empty(t, order)
				assert.Contains(t, err.Error(), "create order item error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				_ = db.Close()
			}()

			queries := sqlc.New() // Initialize with the mock DB
			repo := &orderRepository{q: queries, rawDB: db}

			tt.mockSetup(mock)

			createdOrder, err := repo.CreateOrder(context.Background(), tt.order, tt.orderItems)
			tt.expectedFunc(createdOrder, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
