package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/order/domain"
	"gofiber-skeleton/internal/order/mocks"
	orderUsecase "gofiber-skeleton/internal/order/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderUc := orderUsecase.NewOrderUsecase(mockRepo)

	// Test case 1: Order found
	expectedOrder := &domain.Order{ID: 1, UserID: 1, ProductID: 1, Quantity: 1, TotalPrice: 10.0}
	mockRepo.EXPECT().GetOrder(gomock.Any(), uint(1)).Return(expectedOrder, nil).Times(1)

	order, err := orderUc.GetOrder(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)

	// Test case 2: Order not found
	mockRepo.EXPECT().GetOrder(gomock.Any(), uint(2)).Return(nil, errors.New("order not found")).Times(1)

	order, err = orderUc.GetOrder(context.Background(), 2)
	assert.Error(t, err)
	assert.Nil(t, order)
}



