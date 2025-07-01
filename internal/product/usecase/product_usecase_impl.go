package usecase

import (
	"context"
	"gofiber-skeleton/internal/product/domain"
	"gofiber-skeleton/internal/product/repository"
	"gofiber-skeleton/internal/infra/types"
)

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (pu *productUsecase) GetProduct(ctx context.Context, id types.UUIDv7) (*domain.Product, error) {
	return pu.productRepo.GetProduct(ctx, id)
}

func (pu *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	return pu.productRepo.CreateProduct(ctx, product)
}

func (pu *productUsecase) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return pu.productRepo.UpdateProduct(ctx, product)
}

func (pu *productUsecase) DeleteProduct(ctx context.Context, id types.UUIDv7) error {
	return pu.productRepo.DeleteProduct(ctx, id)
}
