package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/repository/db"
)

type productRepository struct {
	db *sql.DB
	*db.Queries
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{
		db:      db,
		Queries: db.New(db),
	}
}

func (pr *productRepository) CreateProduct(product *domain.Product) error {
	ctx := context.Background()
	createdProduct, err := pr.Queries.CreateProduct(ctx, db.CreateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: sql.NullString{String: *product.Description, Valid: product.Description != nil},
		Price:       product.Price,
		Stock:       product.Stock,
		ImageUrl:    sql.NullString{String: *product.ImageURL, Valid: product.ImageURL != nil},
	})
	if err != nil {
		return err
	}

	product.CreatedAt = createdProduct.CreatedAt
	product.UpdatedAt = createdProduct.UpdatedAt
	return nil
}

func (pr *productRepository) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	ctx := context.Background()
	product, err := pr.Queries.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: &product.Description.String,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    &product.ImageUrl.String,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (pr *productRepository) GetAllProducts() ([]*domain.Product, error) {
	ctx := context.Background()
	products, err := pr.Queries.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	var domainProducts []*domain.Product
	for _, product := range products {
		domainProducts = append(domainProducts, &domain.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: &product.Description.String,
			Price:       product.Price,
			Stock:       product.Stock,
			ImageURL:    &product.ImageUrl.String,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		})
	}
	return domainProducts, nil
}

func (pr *productRepository) UpdateProduct(product *domain.Product) error {
	ctx := context.Background()
	updatedProduct, err := pr.Queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: sql.NullString{String: *product.Description, Valid: product.Description != nil},
		Price:       product.Price,
		Stock:       product.Stock,
		ImageUrl:    sql.NullString{String: *product.ImageURL, Valid: product.ImageURL != nil},
	})
	if err != nil {
		return err
	}

	product.CreatedAt = updatedProduct.CreatedAt
	product.UpdatedAt = updatedProduct.UpdatedAt
	return nil
}

func (pr *productRepository) DeleteProduct(id uuid.UUID) error {
	ctx := context.Background()
	return pr.Queries.DeleteProduct(ctx, id)
}