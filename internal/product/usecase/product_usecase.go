package productusecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

// productUseCase implements the domain.ProductUseCase interface.
type productUseCase struct {
	productRepo domain.ProductRepository
}

// NewProductUseCase creates a new ProductUseCase instance.
func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product with the given details.
// It performs basic validation and then calls the product repository to persist the product.
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

// GetProduct retrieves a product by its ID.
// It validates the ID and then calls the product repository to fetch the product.
func (uc *productUseCase) GetProduct(id string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("product ID is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid product ID")
	}

	return uc.productRepo.GetByID(id)
}

// GetAllProducts retrieves all products from the repository.
func (uc *productUseCase) GetAllProducts() ([]*domain.Product, error) {
	return uc.productRepo.GetAll()
}

// UpdateProduct updates an existing product with the given details.
// It performs validation, retrieves the existing product, updates its fields, and then persists the changes.
func (uc *productUseCase) UpdateProduct(id, name, description string, price float64, stock int, imageURL string) (*domain.Product, error) {
	// Validate input
	if id == "" {
		return nil, errors.New("product ID is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid product ID")
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

// DeleteProduct deletes a product by its ID.
// It validates the ID and then calls the product repository to delete the product.
func (uc *productUseCase) DeleteProduct(id string) error {
	if id == "" {
		return errors.New("product ID is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid product ID")
	}

	return uc.productRepo.Delete(id)
}

// UpdateStock updates the stock quantity for a product.
// It validates the ID and quantity, then calls the product repository to update the stock.
func (uc *productUseCase) UpdateStock(id string, quantity int) error {
	if id == "" {
		return errors.New("product ID is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid product ID")
	}

	return uc.productRepo.UpdateStock(id, quantity)
}
