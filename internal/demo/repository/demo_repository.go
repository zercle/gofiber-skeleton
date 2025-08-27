package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

type DemoRepository interface {
	PerformTransaction(ctx context.Context) error
	GetJoinedData(ctx context.Context) ([]sqlc.GetOrdersWithItemsAndProductsRow, error)
}

type demoRepository struct {
	q     *sqlc.Queries // The generated Queries struct
	rawDB *sql.DB       // The underlying raw DB connection
}

func NewDemoRepository(db *sql.DB) DemoRepository {
	return &demoRepository{
		q:     sqlc.New(),
		rawDB: db,
	}
}

func (r *demoRepository) PerformTransaction(ctx context.Context) error {
	tx, err := r.rawDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	// Pass tx directly to each sqlc.Queries method call
	productID := uuid.New()
	_, err = r.q.CreateProduct(ctx, tx, sqlc.CreateProductParams{
		Name:        "Product for Transaction Demo",
		Description: sql.NullString{String: "This is a product for transaction demo.", Valid: true},
		Price:       "100.00", // Price is string in sqlc.Product
		Stock:       10,
		ImageUrl:    sql.NullString{String: "http://example.com/demo.jpg", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	_, err = r.q.UpdateProductStock(ctx, tx, sqlc.UpdateProductStockParams{
		ID:    productID,
		Stock: 5, // Reduce stock by 5
	})
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	err = r.q.DeleteProduct(ctx, tx, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *demoRepository) GetJoinedData(ctx context.Context) ([]sqlc.GetOrdersWithItemsAndProductsRow, error) {
	// Use the non-transactional querier (r.q) and pass the raw DB connection (r.rawDB)
	data, err := r.q.GetOrdersWithItemsAndProducts(ctx, r.rawDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get joined data: %w", err)
	}
	return data, nil
}
