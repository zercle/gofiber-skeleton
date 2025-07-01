package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/product/domain"
	"gofiber-skeleton/internal/product/mocks"
	productUsecase "gofiber-skeleton/internal/product/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productUc := productUsecase.NewProductUsecase(mockRepo)

	// Test case 1: Product found
	expectedProduct := &domain.Product{ID: 1, Name: "testproduct", Price: 10.0}
	mockRepo.EXPECT().GetProduct(gomock.Any(), uint(1)).Return(expectedProduct, nil).Times(1)

	product, err := productUc.GetProduct(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)

	// Test case 2: Product not found
	mockRepo.EXPECT().GetProduct(gomock.Any(), uint(2)).Return(nil, errors.New("product not found")).Times(1)

	product, err = productUc.GetProduct(context.Background(), 2)
	assert.Error(t, err)
	assert.Nil(t, product)
}
