package service_test

import (
	"testing"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/product/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockProductRepository is a simple mock for testing
type MockProductRepository struct {
	products map[uuid.UUID]*domain.Product
	create   func(*domain.Product) error
	getByID  func(uuid.UUID) (*domain.Product, error)
	getAll   func() ([]*domain.Product, error)
	update   func(*domain.Product) error
	delete   func(uuid.UUID) error
	updateStock func(uuid.UUID, int) error
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	if m.create != nil {
		return m.create(product)
	}
	return nil
}

func (m *MockProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	return m.products[id], nil
}

func (m *MockProductRepository) GetAll() ([]*domain.Product, error) {
	if m.getAll != nil {
		return m.getAll()
	}
	products := make([]*domain.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	if m.update != nil {
		return m.update(product)
	}
	return nil
}

func (m *MockProductRepository) Delete(id uuid.UUID) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}

func (m *MockProductRepository) UpdateStock(id uuid.UUID, quantity int) error {
	if m.updateStock != nil {
		return m.updateStock(id, quantity)
	}
	return nil
}

func TestProductService_CreateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	productService := service.NewProductService(mockRepo)

	t.Run("successful product creation", func(t *testing.T) {
		req := &domain.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       29.99,
			Stock:       100,
			ImageURL:    "https://example.com/image.jpg",
		}

		product, err := productService.CreateProduct(req)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, req.Name, product.Name)
		assert.Equal(t, req.Description, product.Description)
		assert.Equal(t, req.Price, product.Price)
		assert.Equal(t, req.Stock, product.Stock)
		assert.Equal(t, req.ImageURL, product.ImageURL)
		assert.NotEqual(t, uuid.Nil, product.ID)
	})

	t.Run("empty product name", func(t *testing.T) {
		req := &domain.CreateProductRequest{
			Name:        "",
			Description: "Test Description",
			Price:       29.99,
			Stock:       100,
		}

		product, err := productService.CreateProduct(req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product name is required", err.Error())
	})

	t.Run("negative price", func(t *testing.T) {
		req := &domain.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       -10.0,
			Stock:       100,
		}

		product, err := productService.CreateProduct(req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product price cannot be negative", err.Error())
	})

	t.Run("negative stock", func(t *testing.T) {
		req := &domain.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       29.99,
			Stock:       -10,
		}

		product, err := productService.CreateProduct(req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product stock cannot be negative", err.Error())
	})
}

func TestProductService_GetProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	productService := service.NewProductService(mockRepo)

	t.Run("nil product ID", func(t *testing.T) {
		product, err := productService.GetProduct(uuid.Nil)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "invalid product ID", err.Error())
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	productService := service.NewProductService(mockRepo)

	t.Run("nil product ID", func(t *testing.T) {
		req := &domain.UpdateProductRequest{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       39.99,
			Stock:       150,
		}

		product, err := productService.UpdateProduct(uuid.Nil, req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "invalid product ID", err.Error())
	})

	t.Run("empty product name", func(t *testing.T) {
		req := &domain.UpdateProductRequest{
			Name:        "",
			Description: "Updated Description",
			Price:       39.99,
			Stock:       150,
		}

		product, err := productService.UpdateProduct(uuid.New(), req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product name is required", err.Error())
	})

	t.Run("negative price", func(t *testing.T) {
		req := &domain.UpdateProductRequest{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       -10.0,
			Stock:       150,
		}

		product, err := productService.UpdateProduct(uuid.New(), req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product price cannot be negative", err.Error())
	})

	t.Run("negative stock", func(t *testing.T) {
		req := &domain.UpdateProductRequest{
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       39.99,
			Stock:       -10,
		}

		product, err := productService.UpdateProduct(uuid.New(), req)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product stock cannot be negative", err.Error())
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	productService := service.NewProductService(mockRepo)

	t.Run("nil product ID", func(t *testing.T) {
		err := productService.DeleteProduct(uuid.Nil)
		assert.Error(t, err)
		assert.Equal(t, "invalid product ID", err.Error())
	})
}

func TestProductService_UpdateProductStock(t *testing.T) {
	mockRepo := &MockProductRepository{}
	productService := service.NewProductService(mockRepo)

	t.Run("nil product ID", func(t *testing.T) {
		err := productService.UpdateProductStock(uuid.Nil, 10)
		assert.Error(t, err)
		assert.Equal(t, "invalid product ID", err.Error())
	})

	t.Run("zero quantity", func(t *testing.T) {
		err := productService.UpdateProductStock(uuid.New(), 0)
		assert.Error(t, err)
		assert.Equal(t, "quantity must be positive", err.Error())
	})

	t.Run("negative quantity", func(t *testing.T) {
		err := productService.UpdateProductStock(uuid.New(), -10)
		assert.Error(t, err)
		assert.Equal(t, "quantity must be positive", err.Error())
	})
}