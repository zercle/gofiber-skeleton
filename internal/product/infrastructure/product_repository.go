//go:generate mockgen -source=product_repository.go -destination=../mocks/mock_product_repository.go -package=mocks
package infrastructure

import (
	"context"
	"gofiber-skeleton/internal/product/domain"
)

type ProductRepository interface {
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}
