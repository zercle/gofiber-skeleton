package domain

import (
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

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(product *Product) error
	GetByID(id string) (*Product, error)
	GetAll() ([]*Product, error)
	Update(product *Product) error
	Delete(id string) error
	UpdateStock(id string, quantity int) error
}

// ProductUseCase defines the interface for product business logic
type ProductUseCase interface {
	CreateProduct(name, description string, price float64, stock int, imageURL string) (*Product, error)
	GetProduct(id string) (*Product, error)
	GetAllProducts() ([]*Product, error)
	UpdateProduct(id, name, description string, price float64, stock int, imageURL string) (*Product, error)
	DeleteProduct(id string) error
	UpdateStock(id string, quantity int) error
}