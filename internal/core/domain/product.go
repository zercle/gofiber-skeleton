package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the e-commerce system
type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest represents the request to update an existing product
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	ImageURL    string  `json:"image_url"`
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(product *Product) error
	GetByID(id uuid.UUID) (*Product, error)
	GetAll() ([]*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
	UpdateStock(id uuid.UUID, quantity int) error
}

// ProductService defines the interface for product business logic
type ProductService interface {
	CreateProduct(req *CreateProductRequest) (*Product, error)
	GetProduct(id uuid.UUID) (*Product, error)
	GetAllProducts() ([]*Product, error)
	UpdateProduct(id uuid.UUID, req *UpdateProductRequest) (*Product, error)
	DeleteProduct(id uuid.UUID) error
	UpdateProductStock(id uuid.UUID, quantity int) error
}