//go:generate mockgen -source=product_repository.go -destination=mocks/mock_product_repository.go -package=mocks

package repository

import (
	"context"

	"gofiber-skeleton/internal/core/domain"
	"gofiber-skeleton/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	
	// Convert price to pgtype.Numeric
	price := pgtype.Numeric{}
	if err := price.Scan(product.Price); err != nil {
		return err
	}

	// Convert description to pgtype.Text
	description := pgtype.Text{}
	if product.Description != "" {
		description.Scan(product.Description)
	}

	// Convert image URL to pgtype.Text
	imageURL := pgtype.Text{}
	if product.ImageURL != "" {
		imageURL.Scan(product.ImageURL)
	}

	dbProduct, err := r.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:        product.Name,
		Description: description,
		Price:       price,
		Stock:       int32(product.Stock),
		ImageUrl:    imageURL,
	})
	if err != nil {
		return err
	}

	// Update the product with the generated ID and timestamps
	product.ID = dbProduct.ID.Bytes
	product.CreatedAt = dbProduct.CreatedAt.Time
	product.UpdatedAt = dbProduct.UpdatedAt.Time
	return nil
}

func (r *ProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)
	
	dbProduct, err := r.queries.GetProduct(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	// Convert price from pgtype.Numeric to float64
	var price float64
	if dbProduct.Price.Valid {
		if err := dbProduct.Price.Scan(&price); err != nil {
			return nil, err
		}
	}

	return &domain.Product{
		ID:          dbProduct.ID.Bytes,
		Name:        dbProduct.Name,
		Description: dbProduct.Description.String,
		Price:       price,
		Stock:       int(dbProduct.Stock),
		ImageURL:    dbProduct.ImageUrl.String,
		CreatedAt:   dbProduct.CreatedAt.Time,
		UpdatedAt:   dbProduct.UpdatedAt.Time,
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
		// Convert price from pgtype.Numeric to float64
		var price float64
		if dbProduct.Price.Valid {
			if err := dbProduct.Price.Scan(&price); err != nil {
				return nil, err
			}
		}

		products[i] = &domain.Product{
			ID:          dbProduct.ID.Bytes,
			Name:        dbProduct.Name,
			Description: dbProduct.Description.String,
			Price:       price,
			Stock:       int(dbProduct.Stock),
			ImageURL:    dbProduct.ImageUrl.String,
			CreatedAt:   dbProduct.CreatedAt.Time,
			UpdatedAt:   dbProduct.UpdatedAt.Time,
		}
	}

	return products, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	ctx := context.Background()
	
	// Convert price to pgtype.Numeric
	price := pgtype.Numeric{}
	if err := price.Scan(product.Price); err != nil {
		return err
	}

	// Convert description to pgtype.Text
	description := pgtype.Text{}
	if product.Description != "" {
		description.Scan(product.Description)
	}

	// Convert image URL to pgtype.Text
	imageURL := pgtype.Text{}
	if product.ImageURL != "" {
		imageURL.Scan(product.ImageURL)
	}

	// Convert ID to pgtype.UUID
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(product.ID)

	dbProduct, err := r.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          pgUUID,
		Name:        product.Name,
		Description: description,
		Price:       price,
		Stock:       int32(product.Stock),
		ImageUrl:    imageURL,
	})
	if err != nil {
		return err
	}

	// Update the product with the new timestamps
	product.UpdatedAt = dbProduct.UpdatedAt.Time
	return nil
}

func (r *ProductRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)
	
	return r.queries.DeleteProduct(ctx, pgUUID)
}

func (r *ProductRepository) UpdateStock(id uuid.UUID, quantity int) error {
	ctx := context.Background()
	
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)
	
	_, err := r.queries.UpdateProductStock(ctx, db.UpdateProductStockParams{
		ID:     pgUUID,
		Stock:  int32(quantity),
	})
	return err
}