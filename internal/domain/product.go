package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	ImageURL    *string   `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepository interface {
	CreateProduct(product *Product) error
	GetProductByID(id uuid.UUID) (*Product, error)
	GetAllProducts() ([]*Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id uuid.UUID) error
}

type ProductUseCase interface {
	CreateProduct(product *Product) error
	GetProductByID(id uuid.UUID) (*Product, error)
	GetAllProducts() ([]*Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id uuid.UUID) error
	ReduceStock(ctx context.Context, productID string, quantity int) error
}