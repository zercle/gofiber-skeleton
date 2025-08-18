package usecase

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

func TestCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	product := &domain.Product{
		Name:        "Test Product",
		Description: new(string),
		Price:       10.99,
		Stock:       100,
		ImageURL:    new(string),
	}
	*product.Description = "A test product"
	*product.ImageURL = "http://example.com/image.jpg"

	// Expect the CreateProduct method to be called once with any product and return no error
	mockProductRepo.EXPECT().CreateProduct(gomock.Any()).Return(nil).Times(1)

	err := productUseCase.CreateProduct(product)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, product.ID) // Check if ID is set
}

func TestCreateProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	product := &domain.Product{
		Name:        "Test Product",
		Description: new(string),
		Price:       10.99,
		Stock:       100,
		ImageURL:    new(string),
	}
	*product.Description = "A test product"
	*product.ImageURL = "http://example.com/image.jpg"

	expectedError := errors.New("database error")
	mockProductRepo.EXPECT().CreateProduct(gomock.Any()).Return(expectedError).Times(1)

	err := productUseCase.CreateProduct(product)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestGetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	productID := uuid.New()
	expectedProduct := &domain.Product{
		ID:          productID,
		Name:        "Retrieved Product",
		Description: new(string),
		Price:       20.00,
		Stock:       50,
		ImageURL:    new(string),
	}
	*expectedProduct.Description = "Retrieved Description"
	*expectedProduct.ImageURL = "http://example.com/retrieved.jpg"

	mockProductRepo.EXPECT().GetProductByID(productID).Return(expectedProduct, nil).Times(1)

	retrievedProduct, err := productUseCase.GetProductByID(productID)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, retrievedProduct)
}

func TestGetProductByIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	productID := uuid.New()
	expectedError := errors.New("product not found")

	mockProductRepo.EXPECT().GetProductByID(productID).Return(nil, expectedError).Times(1)

	retrievedProduct, err := productUseCase.GetProductByID(productID)
	assert.Error(t, err)
	assert.Nil(t, retrievedProduct)
	assert.Equal(t, expectedError, err)
}

func TestGetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	products := []*domain.Product{
		{ID: uuid.New(), Name: "Product 1", Price: 10.0, Stock: 10},
		{ID: uuid.New(), Name: "Product 2", Price: 20.0, Stock: 20},
	}

	mockProductRepo.EXPECT().GetAllProducts().Return(products, nil).Times(1)

	retrievedProducts, err := productUseCase.GetAllProducts()
	assert.NoError(t, err)
	assert.Equal(t, products, retrievedProducts)
}

func TestGetAllProductsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	expectedError := errors.New("database error")
	mockProductRepo.EXPECT().GetAllProducts().Return(nil, expectedError).Times(1)

	retrievedProducts, err := productUseCase.GetAllProducts()
	assert.Error(t, err)
	assert.Nil(t, retrievedProducts)
	assert.Equal(t, expectedError, err)
}

func TestUpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	productID := uuid.New()
	updatedProduct := &domain.Product{
		ID:          productID,
		Name:        "Updated Product",
		Description: new(string),
		Price:       30.00,
		Stock:       75,
		ImageURL:    new(string),
	}
	*updatedProduct.Description = "Updated Description"
	*updatedProduct.ImageURL = "http://example.com/updated.jpg"

	mockProductRepo.EXPECT().UpdateProduct(updatedProduct).Return(nil).Times(1)

	err := productUseCase.UpdateProduct(updatedProduct)
	assert.NoError(t, err)
}

func TestUpdateProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	productID := uuid.New()
	product := &domain.Product{ID: productID, Name: "Product to Update"}
	expectedError := errors.New("update failed")

	mockProductRepo.EXPECT().UpdateProduct(product).Return(expectedError).Times(1)

	err := productUseCase.UpdateProduct(product)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	productID := uuid.New()

	mockProductRepo.EXPECT().DeleteProduct(productID).Return(nil).Times(1)

	err := productUseCase.DeleteProduct(productID)
	assert.NoError(t, err)
}

func TestDeleteProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := domain.NewMockProductRepository(ctrl)
	productUseCase := NewProductUseCase(mockProductRepo)

	productID := uuid.New()
	expectedError := errors.New("delete failed")

	mockProductRepo.EXPECT().DeleteProduct(productID).Return(expectedError).Times(1)

	err := productUseCase.DeleteProduct(productID)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}