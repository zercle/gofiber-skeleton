package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type productUseCase struct {
	productRepository domain.ProductRepository
}

func NewProductUseCase(productRepository domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{
		productRepository: productRepository,
	}
}

func (pu *productUseCase) CreateProduct(product *domain.Product) error {
	product.ID = uuid.New()
	return pu.productRepository.CreateProduct(product)
}

func (pu *productUseCase) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	return pu.productRepository.GetProductByID(id)
}

func (pu *productUseCase) GetAllProducts() ([]*domain.Product, error) {
	return pu.productRepository.GetAllProducts()
}

func (pu *productUseCase) UpdateProduct(product *domain.Product) error {
	return pu.productRepository.UpdateProduct(product)
}

func (pu *productUseCase) DeleteProduct(id uuid.UUID) error {
	return pu.productRepository.DeleteProduct(id)
}

func (pu *productUseCase) ReduceStock(ctx context.Context, productID string, quantity int) error {
	id, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	product, err := pu.productRepository.GetProductByID(id)
	if err != nil {
		return err
	}
	product.Stock -= int32(quantity)
	return pu.productRepository.UpdateProduct(product)
}