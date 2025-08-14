package service

import (
	"errors"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/product/repository"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Price < 0 {
		return nil, errors.New("product price cannot be negative")
	}
	if req.Stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	// Create product domain object
	product := &domain.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}

	// Save to repository
	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProduct(id uuid.UUID) (*domain.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid product ID")
	}

	return s.productRepo.GetByID(id)
}

func (s *ProductService) GetAllProducts() ([]*domain.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) UpdateProduct(id uuid.UUID, req *domain.UpdateProductRequest) (*domain.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid product ID")
	}

	// Validate request
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Price < 0 {
		return nil, errors.New("product price cannot be negative")
	}
	if req.Stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	// Get existing product
	existingProduct, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.Price = req.Price
	existingProduct.Stock = req.Stock
	existingProduct.ImageURL = req.ImageURL

	// Save changes
	if err := s.productRepo.Update(existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid product ID")
	}

	return s.productRepo.Delete(id)
}

func (s *ProductService) UpdateProductStock(id uuid.UUID, quantity int) error {
	if id == uuid.Nil {
		return errors.New("invalid product ID")
	}
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	return s.productRepo.UpdateStock(id, quantity)
}