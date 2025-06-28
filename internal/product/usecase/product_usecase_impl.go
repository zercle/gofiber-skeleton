package usecase

import (
	"context"
	"gofiber-skeleton/internal/product/domain"
	"gofiber-skeleton/internal/product/infrastructure"
)

type productUsecase struct {
	productRepo infrastructure.ProductRepository
}

func NewProductUsecase(productRepo infrastructure.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (pu *productUsecase) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	return pu.productRepo.GetProduct(ctx, id)
}

func (pu *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	return pu.productRepo.CreateProduct(ctx, product)
}

func (pu *productUsecase) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return pu.productRepo.UpdateProduct(ctx, product)
}

func (pu *productUsecase) DeleteProduct(ctx context.Context, id uint) error {
	return pu.productRepo.DeleteProduct(ctx, id)
}
