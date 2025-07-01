//go:generate mockgen -source=product_usecase.go -destination=../mocks/mock_product_usecase.go -package=mocks
package usecase

import (
	"context"
	"gofiber-skeleton/internal/product/domain"
	"gofiber-skeleton/pkg/types"
)

type ProductUsecase interface {
	GetProduct(ctx context.Context, id types.UUIDv7) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id types.UUIDv7) error
}