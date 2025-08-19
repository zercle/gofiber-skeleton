package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	sqldb "github.com/zercle/gofiber-skeleton/internal/repository/db"
)

type productRepository struct {
	db *sql.DB
	*sqldb.Queries
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{
		db:      db,
		Queries: sqldb.New(db),
	}
}

func (pr *productRepository) CreateProduct(product *domain.Product) error {
	ctx := context.Background()
	createdProduct, err := pr.Queries.CreateProduct(ctx, sqldb.CreateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: sql.NullString{String: *product.Description, Valid: product.Description != nil},
		Price:       fmt.Sprintf("%.2f", product.Price), // Convert float64 to string
		Stock:       product.Stock,
		ImageUrl:    sql.NullString{String: *product.ImageURL, Valid: product.ImageURL != nil},
	})
	if err != nil {
		return err
	}

	product.CreatedAt = createdProduct.CreatedAt.Time
	product.UpdatedAt = createdProduct.UpdatedAt.Time
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
		Price:       parseFloat(product.Price), // Convert string to float64
		Stock:       product.Stock,
		ImageURL:    &product.ImageUrl.String,
		CreatedAt:   product.CreatedAt.Time,
		UpdatedAt:   product.UpdatedAt.Time,
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
			Price:       parseFloat(product.Price), // Convert string to float64
			Stock:       product.Stock,
			ImageURL:    &product.ImageUrl.String,
			CreatedAt:   product.CreatedAt.Time,
			UpdatedAt:   product.UpdatedAt.Time,
		})
	}
	return domainProducts, nil
}

func (pr *productRepository) UpdateProduct(product *domain.Product) error {
	ctx := context.Background()
	updatedProduct, err := pr.Queries.UpdateProduct(ctx, sqldb.UpdateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: sql.NullString{String: *product.Description, Valid: product.Description != nil},
		Price:       fmt.Sprintf("%.2f", product.Price), // Convert float64 to string
		Stock:       product.Stock,
		ImageUrl:    sql.NullString{String: *product.ImageURL, Valid: product.ImageURL != nil},
	})
	if err != nil {
		return err
	}

	product.CreatedAt = updatedProduct.CreatedAt.Time
	product.UpdatedAt = updatedProduct.UpdatedAt.Time
	return nil
}

func (pr *productRepository) DeleteProduct(id uuid.UUID) error {
	ctx := context.Background()
	return pr.Queries.DeleteProduct(ctx, id)
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
		return 0.0
	}
	return f
}