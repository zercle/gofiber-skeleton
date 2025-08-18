package usecase

import (
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