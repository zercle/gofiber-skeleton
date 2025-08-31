//go:generate go run go.uber.org/mock/mockgen -source=product.go -destination=./mock/product_mock.go -package=productmock
package productmodule

import (
	"errors"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductRepository defines the interface for product data operations.
type ProductRepository interface {
	// Create adds a new product to the repository.
	Create(product *Product) error
	// GetByID retrieves a product by its ID.
	GetByID(id string) (*Product, error)
	// GetAll retrieves all products.
	GetAll() ([]*Product, error)
	// Update updates an existing product.
	Update(product *Product) error
	// Delete removes a product by its ID.
	Delete(id string) error
	// UpdateStock updates the stock quantity for a product.
	UpdateStock(id string, quantity int) error
}

// ProductUseCase defines the interface for product business logic.
type ProductUseCase interface {
	// CreateProduct creates a new product with the given details.
	CreateProduct(name, description string, price float64, stock int, imageURL string) (*Product, error)
	// GetProduct retrieves a product by its ID.
	GetProduct(id string) (*Product, error)
	// GetAllProducts retrieves all products.
	GetAllProducts() ([]*Product, error)
	// UpdateProduct updates an existing product with the given details.
	UpdateProduct(id, name, description string, price float64, stock int, imageURL string) (*Product, error)
	// DeleteProduct deletes a product by its ID.
	DeleteProduct(id string) error
	// UpdateStock updates the stock quantity for a product.
	UpdateStock(id string, quantity int) error
}

// ErrProductNotFound indicates that a product was not found.
var ErrProductNotFound = errors.New("product not found")
