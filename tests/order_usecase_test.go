package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/order/domain"
	"gofiber-skeleton/internal/order/mocks"
	orderUsecase "gofiber-skeleton/internal/order/usecase"
	"gofiber-skeleton/pkg/types"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderUc := orderUsecase.NewOrderUsecase(mockRepo)

	// Test case 1: Order found
	orderID := types.NewUUIDv7()
	expectedOrder := &domain.Order{ID: orderID, UserID: types.NewUUIDv7(), ProductID: types.NewUUIDv7(), Quantity: 1, TotalPrice: 10.0, Status: "pending"}
	mockRepo.EXPECT().GetOrder(gomock.Any(), orderID).Return(expectedOrder, nil).Times(1)

	order, err := orderUc.GetOrder(context.Background(), orderID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)

	// Test case 2: Order not found
	notFoundOrderID := types.NewUUIDv7()
	mockRepo.EXPECT().GetOrder(gomock.Any(), notFoundOrderID).Return(nil, errors.New("order not found")).Times(1)

	order, err = orderUc.GetOrder(context.Background(), notFoundOrderID)
	assert.Error(t, err)
	assert.Nil(t, order)
}



