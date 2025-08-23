package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type productUseCase struct {
	productRepo domain.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
	}
}

func (uc *productUseCase) CreateProduct(name, description string, price float64, stock int, imageURL string) (*domain.Product, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("product name is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	// Create product
	product := &domain.Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    imageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *productUseCase) GetProduct(id string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("product ID is required")
	}

	return uc.productRepo.GetByID(id)
}

func (uc *productUseCase) GetAllProducts() ([]*domain.Product, error) {
	return uc.productRepo.GetAll()
}

func (uc *productUseCase) UpdateProduct(id, name, description string, price float64, stock int, imageURL string) (*domain.Product, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("product ID is required")
	}
	if name == "" {
		return nil, errors.New("product name is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	// Get existing product
	existingProduct, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingProduct.Name = name
	existingProduct.Description = description
	existingProduct.Price = price
	existingProduct.Stock = stock
	existingProduct.ImageURL = imageURL
	existingProduct.UpdatedAt = time.Now()

	if err := uc.productRepo.Update(existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (uc *productUseCase) DeleteProduct(id string) error {
	if id == "" {
		return errors.New("product ID is required")
	}

	return uc.productRepo.Delete(id)
}

func (uc *productUseCase) UpdateStock(id string, quantity int) error {
	if id == "" {
		return errors.New("product ID is required")
	}

	return uc.productRepo.UpdateStock(id, quantity)
}