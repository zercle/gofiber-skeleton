package repository

import (
	"context"
	"database/sql"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/repository/db"
)

type productRepository struct {
	db *db.Queries
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *db.Queries) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	ctx := context.Background()
	
	dbProduct, err := r.db.CreateProduct(ctx, db.CreateProductParams{
		Name:        product.Name,
		Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:       fmt.Sprintf("%.2f", product.Price),
		Stock:       int32(product.Stock),
		ImageUrl:    sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""},
	})
	if err != nil {
		return err
	}

	product.ID = dbProduct.ID.String()
	product.CreatedAt = dbProduct.CreatedAt
	product.UpdatedAt = dbProduct.UpdatedAt
	return nil
}

func (r *productRepository) GetByID(id string) (*domain.Product, error) {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	dbProduct, err := r.db.GetProductByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	price, _ := strconv.ParseFloat(dbProduct.Price, 64)
	return &domain.Product{
		ID:          dbProduct.ID.String(),
		Name:        dbProduct.Name,
		Description: dbProduct.Description.String,
		Price:       price,
		Stock:       int(dbProduct.Stock),
		ImageURL:    dbProduct.ImageUrl.String,
		CreatedAt:   dbProduct.CreatedAt.Time,
		UpdatedAt:   dbProduct.UpdatedAt.Time,
	}, nil
}

func (r *productRepository) GetAll() ([]*domain.Product, error) {
	ctx := context.Background()
	
	dbProducts, err := r.db.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*domain.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = &domain.Product{
			ID:          dbProduct.ID.String(),
			Name:        dbProduct.Name,
			Description: dbProduct.Description,
			Price:       dbProduct.Price,
			Stock:       int(dbProduct.Stock),
			ImageURL:    dbProduct.ImageUrl.String,
			CreatedAt:   dbProduct.CreatedAt,
			UpdatedAt:   dbProduct.UpdatedAt,
		}
	}

	return products, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(product.ID)
	if err != nil {
		return err
	}

	dbProduct, err := r.db.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
		ImageUrl:    sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""},
	})
	if err != nil {
		return err
	}

	product.UpdatedAt = dbProduct.UpdatedAt
	return nil
}

func (r *productRepository) Delete(id string) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return err
	}

	return r.db.DeleteProduct(ctx, uuid)
}

func (r *productRepository) UpdateStock(id string, quantity int) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateProductStock(ctx, db.UpdateProductStockParams{
		ID:     uuid,
		Stock:  int32(quantity),
	})
	return err
}