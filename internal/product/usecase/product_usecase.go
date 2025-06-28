//go:generate mockgen -source=product_usecase.go -destination=../mocks/mock_product_usecase.go -package=mocks
package usecase

import (
	"context"
	"gofiber-skeleton/internal/product/domain"
)

type ProductUsecase interface {
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}