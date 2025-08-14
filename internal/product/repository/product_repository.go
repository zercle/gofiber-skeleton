package repository

import (
	"context"
	"database/sql"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/platform/db"

	"github.com/google/uuid"
)

type ProductRepository struct {
	queries *db.Queries
}

func NewProductRepository(queries *db.Queries) *ProductRepository {
	return &ProductRepository{
		queries: queries,
	}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	ctx := context.Background()
	
	dbProduct, err := r.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:        product.Name,
		Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:       product.Price,
		Stock:       int32(product.Stock),
		ImageUrl:    sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""},
	})
	if err != nil {
		return err
	}

	// Update the product with the generated ID and timestamps
	product.ID = dbProduct.ID
	product.CreatedAt = dbProduct.CreatedAt
	product.UpdatedAt = dbProduct.UpdatedAt
	return nil
}

func (r *ProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	ctx := context.Background()
	
	dbProduct, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          dbProduct.ID,
		Name:        dbProduct.Name,
		Description: dbProduct.Description.String,
		Price:       float64(dbProduct.Price),
		Stock:       int(dbProduct.Stock),
		ImageURL:    dbProduct.ImageUrl.String,
		CreatedAt:   dbProduct.CreatedAt,
		UpdatedAt:   dbProduct.UpdatedAt,
	}, nil
}

func (r *ProductRepository) GetAll() ([]*domain.Product, error) {
	ctx := context.Background()
	
	dbProducts, err := r.queries.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*domain.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = &domain.Product{
			ID:          dbProduct.ID,
			Name:        dbProduct.Name,
			Description: dbProduct.Description.String,
			Price:       float64(dbProduct.Price),
			Stock:       int(dbProduct.Stock),
			ImageURL:    dbProduct.ImageUrl.String,
			CreatedAt:   dbProduct.CreatedAt,
			UpdatedAt:   dbProduct.UpdatedAt,
		}
	}

	return products, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	ctx := context.Background()
	
	dbProduct, err := r.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:       product.Price,
		Stock:       int32(product.Stock),
		ImageUrl:    sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""},
	})
	if err != nil {
		return err
	}

	// Update the product with the new timestamps
	product.UpdatedAt = dbProduct.UpdatedAt
	return nil
}

func (r *ProductRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()
	return r.queries.DeleteProduct(ctx, id)
}

func (r *ProductRepository) UpdateStock(id uuid.UUID, quantity int) error {
	ctx := context.Background()
	
	_, err := r.queries.UpdateProductStock(ctx, db.UpdateProductStockParams{
		ID:     id,
		Stock:  int32(quantity),
	})
	return err
}