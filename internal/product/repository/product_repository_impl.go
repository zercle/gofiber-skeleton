package repository

import (
	"context"
	"gofiber-skeleton/internal/product/domain"
	"gofiber-skeleton/pkg/types"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (pr *productRepository) GetProduct(ctx context.Context, id types.UUIDv7) (*domain.Product, error) {
	var product domain.Product
	if err := pr.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	if err := pr.db.WithContext(ctx).Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if err := pr.db.WithContext(ctx).Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (pr *productRepository) DeleteProduct(ctx context.Context, id types.UUIDv7) error {
	if err := pr.db.WithContext(ctx).Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
