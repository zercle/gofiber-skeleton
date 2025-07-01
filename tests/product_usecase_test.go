package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/product/domain"
	"gofiber-skeleton/internal/product/mocks"
	productUsecase "gofiber-skeleton/internal/product/usecase"
	"gofiber-skeleton/internal/infra/types"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productUc := productUsecase.NewProductUsecase(mockRepo)

	// Test case 1: Product found
	productID := types.NewUUIDv7()
	expectedProduct := &domain.Product{ID: productID, Name: "testproduct", Price: 10.0}
	mockRepo.EXPECT().GetProduct(gomock.Any(), productID).Return(expectedProduct, nil).Times(1)

	product, err := productUc.GetProduct(context.Background(), productID)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)

	// Test case 2: Product not found
	notFoundProductID := types.NewUUIDv7()
	mockRepo.EXPECT().GetProduct(gomock.Any(), notFoundProductID).Return(nil, errors.New("product not found")).Times(1)

	product, err = productUc.GetProduct(context.Background(), notFoundProductID)
	assert.Error(t, err)
	assert.Nil(t, product)
}



