package productrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/productmodule"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

// productRepository implements the productmodule.ProductRepository interface.
type productRepository struct {
	q     *sqlc.Queries // The generated Queries struct (holds methods)
	rawDB *sql.DB       // The underlying DB connection (passed as DBTX)
}

// NewProductRepository creates a new ProductRepository instance.
func NewProductRepository(db *sql.DB) productmodule.ProductRepository {
	return &productRepository{
		q:     sqlc.New(), // Call the parameterless New()
		rawDB: db,         // Store the actual DB connection
	}
}

// Create adds a new product to the database.
func (r *productRepository) Create(product *productmodule.Product) error {
	ctx := context.Background()

	dbProduct, err := r.q.CreateProduct(ctx, r.rawDB, sqlc.CreateProductParams{
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
	product.CreatedAt = dbProduct.CreatedAt.Time
	product.UpdatedAt = dbProduct.UpdatedAt.Time
	return nil
}

// GetByID retrieves a product by its ID from the database.
func (r *productRepository) GetByID(id string) (*productmodule.Product, error) {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id) // Use uuid.Parse directly
	if err != nil {
		return nil, err
	}

	dbProduct, err := r.q.GetProductByID(ctx, r.rawDB, parsedUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, productmodule.ErrProductNotFound
		}
		return nil, err
	}

	price, err := strconv.ParseFloat(dbProduct.Price, 64) // Handle error for ParseFloat
	if err != nil {
		return nil, fmt.Errorf("failed to parse product price: %w", err)
	}

	return &productmodule.Product{
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

// GetAll retrieves all products from the database.
func (r *productRepository) GetAll() ([]*productmodule.Product, error) {
	ctx := context.Background()

	dbProducts, err := r.q.GetAllProducts(ctx, r.rawDB)
	if err != nil {
		return nil, err
	}

	products := make([]*productmodule.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		price, err := strconv.ParseFloat(dbProduct.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse product price for product ID %s: %w", dbProduct.ID.String(), err)
		}
		products[i] = &productmodule.Product{
			ID:          dbProduct.ID.String(),
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

// Update updates an existing product in the database.
func (r *productRepository) Update(product *productmodule.Product) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(product.ID)
	if err != nil {
		return err
	}

	priceStr := fmt.Sprintf("%.2f", product.Price) // Convert float64 to string for DB
	dbProduct, err := r.q.UpdateProduct(ctx, r.rawDB, sqlc.UpdateProductParams{
		ID:          parsedUUID,
		Name:        product.Name,
		Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:       priceStr,
		Stock:       int32(product.Stock),
		ImageUrl:    sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""},
	})
	if err != nil {
		return err
	}

	product.UpdatedAt = dbProduct.UpdatedAt.Time
	return nil
}

// Delete removes a product by its ID from the database.
func (r *productRepository) Delete(id string) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.q.DeleteProduct(ctx, r.rawDB, parsedUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return productmodule.ErrProductNotFound
		}
		return err
	}
	return nil
}

// UpdateStock updates the stock quantity for a product in the database.
func (r *productRepository) UpdateStock(id string, delta int) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = r.q.UpdateProductStock(ctx, r.rawDB, sqlc.UpdateProductStockParams{
		ID:    parsedUUID,
		Stock: int32(delta),
	})
	return err
}
