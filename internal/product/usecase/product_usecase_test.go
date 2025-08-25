package productusecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/domain/mock"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewProductUseCase(mockRepo)

	testCases := []struct {
		name        string
		inputName   string
		description string
		price       float64
		stock       int
		imageURL    string
		mockSetup   func()
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "empty name",
			inputName:   "",
			description: "desc",
			price:       10.0,
			stock:       10,
			imageURL:    "url",
			mockSetup:   func() {},
			wantErr:     true,
			errMsg:      "product name is required",
		},
		{
			name:        "negative price",
			inputName:   "name",
			description: "desc",
			price:       -10.0,
			stock:       10,
			imageURL:    "url",
			mockSetup:   func() {},
			wantErr:     true,
			errMsg:      "price cannot be negative",
		},
		{
			name:        "negative stock",
			inputName:   "name",
			description: "desc",
			price:       10.0,
			stock:       -10,
			imageURL:    "url",
			mockSetup:   func() {},
			wantErr:     true,
			errMsg:      "stock cannot be negative",
		},
		{
			name:        "repository error",
			inputName:   "Test Product",
			description: "desc error",
			price:       10.99,
			stock:       100,
			imageURL:    "url",
			mockSetup: func() {
				mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name:        "successful creation",
			inputName:   "Test Product",
			description: "A test description",
			price:       10.99,
			stock:       100,
			imageURL:    "http://example.com/image.jpg",
			mockSetup: func() {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, "Test Product", p.Name)
					assert.Equal(t, "A test description", p.Description)
					assert.Equal(t, 10.99, p.Price)
					assert.Equal(t, 100, p.Stock)
					assert.Equal(t, "http://example.com/image.jpg", p.ImageURL)
					assert.NotEmpty(t, p.ID)
					assert.False(t, p.CreatedAt.IsZero())
					assert.False(t, p.UpdatedAt.IsZero())
					return nil
				})
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			product, err := usecase.CreateProduct(tc.inputName, tc.description, tc.price, tc.stock, tc.imageURL)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, product)
				assert.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, tc.inputName, product.Name)
			}
		})
	}
}

func TestProductUseCase_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewProductUseCase(mockProductRepo)

	productID := uuid.New().String()
	expectedProduct := &domain.Product{
		ID:          productID,
		Name:        "Existing Product",
		Description: "Desc",
		Price:       20.00,
		Stock:       50,
		ImageURL:    "url",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful product retrieval", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(expectedProduct, nil)

		product, err := usecase.GetProduct(productID)
		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, expectedProduct.ID, product.ID)
	})

	t.Run("invalid product ID", func(t *testing.T) {
		product, err := usecase.GetProduct("invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "invalid product ID")
	})

	t.Run("empty product ID", func(t *testing.T) {
		product, err := usecase.GetProduct("")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "product ID is required")
	})

	t.Run("repository returns not found error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(nil, domain.ErrProductNotFound)

		product, err := usecase.GetProduct(productID)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, domain.ErrProductNotFound.Error())
	})

	t.Run("repository returns generic error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(nil, errors.New("db error"))

		product, err := usecase.GetProduct(productID)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "db error")
	})
}

func TestProductUseCase_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewProductUseCase(mockProductRepo)

	expectedProducts := []*domain.Product{
		{ID: uuid.New().String(), Name: "Product 1"},
		{ID: uuid.New().String(), Name: "Product 2"},
	}

	t.Run("successful retrieval of all products", func(t *testing.T) {
		mockProductRepo.EXPECT().GetAll().Return(expectedProducts, nil)

		products, err := usecase.GetAllProducts()
		require.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 2)
		assert.Equal(t, expectedProducts[0].Name, products[0].Name)
	})

	t.Run("repository returns an error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetAll().Return(nil, errors.New("db error"))

		products, err := usecase.GetAllProducts()
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.EqualError(t, err, "db error")
	})
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewProductUseCase(mockProductRepo)

	productID := uuid.New().String()
	originalProduct := &domain.Product{
		ID:          productID,
		Name:        "Original Name",
		Description: "Original Desc",
		Price:       10.00,
		Stock:       10,
		ImageURL:    "original.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful product update", func(t *testing.T) {
		updatedName := "Updated Name"
		updatedDesc := "Updated Desc"
		updatedPrice := 20.00
		updatedStock := 20
		updatedImageURL := "updated.jpg"

		mockProductRepo.EXPECT().GetByID(productID).Return(originalProduct, nil)
		mockProductRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
			assert.Equal(t, updatedName, p.Name)
			assert.Equal(t, updatedDesc, p.Description)
			assert.Equal(t, updatedPrice, p.Price)
			assert.Equal(t, updatedStock, p.Stock)
			assert.Equal(t, updatedImageURL, p.ImageURL)
			return nil
		})

		product, err := usecase.UpdateProduct(productID, updatedName, updatedDesc, updatedPrice, updatedStock, updatedImageURL)
		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, updatedName, product.Name)
	})

	t.Run("invalid input - invalid ID", func(t *testing.T) {
		product, err := usecase.UpdateProduct("invalid-uuid", "name", "desc", 10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "invalid product ID")
	})

	t.Run("invalid input - empty ID", func(t *testing.T) {
		product, err := usecase.UpdateProduct("", "name", "desc", 10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "product ID is required")
	})

	t.Run("invalid input - empty name", func(t *testing.T) {
		product, err := usecase.UpdateProduct(productID, "", "desc", 10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "product name is required")
	})

	t.Run("invalid input - negative price", func(t *testing.T) {
		product, err := usecase.UpdateProduct(productID, "name", "desc", -10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "price cannot be negative")
	})

	t.Run("invalid input - negative stock", func(t *testing.T) {
		product, err := usecase.UpdateProduct(productID, "name", "desc", 10.0, -10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "stock cannot be negative")
	})

	t.Run("GetByID returns not found error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(nil, domain.ErrProductNotFound)

		product, err := usecase.UpdateProduct(productID, "name", "desc", 10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, domain.ErrProductNotFound.Error())
	})

	t.Run("GetByID returns generic error", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(nil, errors.New("db error"))

		product, err := usecase.UpdateProduct(productID, "name", "desc", 10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "db error")
	})

	t.Run("repository returns error on update", func(t *testing.T) {
		mockProductRepo.EXPECT().GetByID(productID).Return(originalProduct, nil)
		mockProductRepo.EXPECT().Update(gomock.Any()).Return(errors.New("db update error"))

		product, err := usecase.UpdateProduct(productID, "name", "desc", 10.0, 10, "url")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "db update error")
	})
}

func TestProductUseCase_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewProductUseCase(mockProductRepo)

	productID := uuid.New().String()

	t.Run("successful product deletion", func(t *testing.T) {
		mockProductRepo.EXPECT().Delete(productID).Return(nil)

		err := usecase.DeleteProduct(productID)
		require.NoError(t, err)
	})

	t.Run("invalid product ID", func(t *testing.T) {
		err := usecase.DeleteProduct("invalid-uuid")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid product ID")
	})

	t.Run("empty product ID", func(t *testing.T) {
		err := usecase.DeleteProduct("")
		assert.Error(t, err)
		assert.EqualError(t, err, "product ID is required")
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockProductRepo.EXPECT().Delete(productID).Return(errors.New("db error"))

		err := usecase.DeleteProduct(productID)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}

func TestProductUseCase_UpdateStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock.NewMockProductRepository(ctrl)
	usecase := NewProductUseCase(mockProductRepo)

	productID := uuid.New().String()
	quantity := 5

	t.Run("successful stock update", func(t *testing.T) {
		mockProductRepo.EXPECT().UpdateStock(productID, quantity).Return(nil)

		err := usecase.UpdateStock(productID, quantity)
		require.NoError(t, err)
	})

	t.Run("invalid product ID", func(t *testing.T) {
		err := usecase.UpdateStock("invalid-uuid", quantity)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid product ID")
	})

	t.Run("empty product ID", func(t *testing.T) {
		err := usecase.UpdateStock("", quantity)
		assert.Error(t, err)
		assert.EqualError(t, err, "product ID is required")
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockProductRepo.EXPECT().UpdateStock(productID, quantity).Return(errors.New("db error"))

		err := usecase.UpdateStock(productID, quantity)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}
