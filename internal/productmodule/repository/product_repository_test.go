package productrepository

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zercle/gofiber-skeleton/internal/productmodule"
)

func TestProductRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	repo := NewProductRepository(db)

	product := &productmodule.Product{
		Name:        "New Product",
		Description: "A description",
		Price:       99.99,
		Stock:       10,
		ImageURL:    "http://example.com/new.jpg",
	}

	t.Run("successful product creation", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(uuid.New(), product.Name, sql.NullString{String: product.Description, Valid: product.Description != ""}, fmt.Sprintf("%.2f", product.Price), int32(product.Stock), sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})

		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO products (name, description, price, stock, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(product.Name, sql.NullString{String: product.Description, Valid: product.Description != ""}, fmt.Sprintf("%.2f", product.Price), int32(product.Stock), sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""}).
			WillReturnRows(rows)

		err := repo.Create(product)
		require.NoError(t, err)
		assert.NotEmpty(t, product.ID)
		assert.False(t, product.CreatedAt.IsZero())
		assert.False(t, product.UpdatedAt.IsZero())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error on create", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO products (name, description, price, stock, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(product.Name, sql.NullString{String: product.Description, Valid: product.Description != ""}, fmt.Sprintf("%.2f", product.Price), int32(product.Stock), sql.NullString{String: product.ImageURL, Valid: product.ImageURL != ""}).
			WillReturnError(errors.New("db insert error"))

		err := repo.Create(product)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db insert error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	repo := NewProductRepository(db)

	productID := uuid.New()
	expectedProduct := productmodule.Product{
		ID:          productID.String(),
		Name:        "Test Product",
		Description: "A description",
		Price:       10.99,
		Stock:       100,
		ImageURL:    "http://example.com/image.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful product retrieval by ID", func(t *testing.T) {
		rowsProduct := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, expectedProduct.Name, sql.NullString{String: expectedProduct.Description, Valid: expectedProduct.Description != ""}, fmt.Sprintf("%.2f", expectedProduct.Price), int32(expectedProduct.Stock), sql.NullString{String: expectedProduct.ImageURL, Valid: expectedProduct.ImageURL != ""}, sql.NullTime{Time: expectedProduct.CreatedAt, Valid: true}, sql.NullTime{Time: expectedProduct.UpdatedAt, Valid: true})
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).WillReturnRows(rowsProduct)

		product, err := repo.GetByID(productID.String())
		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, expectedProduct.ID, product.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("product not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).WillReturnError(sql.ErrNoRows) // Simulate no rows found

		product, err := repo.GetByID(productID.String())
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, productmodule.ErrProductNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		product, err := repo.GetByID("invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("database error on get by ID", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).WillReturnError(errors.New("db select error"))

		product, err := repo.GetByID(productID.String())
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "db select error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	repo := NewProductRepository(db)

	product1ID := uuid.New()
	product2ID := uuid.New()

	t.Run("successful retrieval of all products", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products ORDER BY created_at DESC`,
		)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
				AddRow(product1ID, "Product 1", sql.NullString{String: "Desc 1", Valid: true}, fmt.Sprintf("%.2f", 10.00), int32(10), sql.NullString{String: "url1", Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}).
				AddRow(product2ID, "Product 2", sql.NullString{String: "Desc 2", Valid: true}, fmt.Sprintf("%.2f", 20.00), int32(20), sql.NullString{String: "url2", Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}))

		products, err := repo.GetAll()
		require.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 2)
		assert.Equal(t, product1ID.String(), products[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error on get all", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products ORDER BY created_at DESC`,
		)).WillReturnError(errors.New("db select all error"))

		products, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "db select all error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	repo := NewProductRepository(db)

	productID := uuid.New()
	productToUpdate := &productmodule.Product{
		ID:          productID.String(),
		Name:        "Updated Name",
		Description: "Updated Desc",
		Price:       12.34,
		Stock:       50,
		ImageURL:    "http://example.com/updated.jpg",
	}

	t.Run("successful product update", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE products SET name = $2, description = $3, price = $4, stock = $5, image_url = $6, updated_at = NOW() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productID, productToUpdate.Name, sql.NullString{String: productToUpdate.Description, Valid: productToUpdate.Description != ""}, fmt.Sprintf("%.2f", productToUpdate.Price), int32(productToUpdate.Stock), sql.NullString{String: productToUpdate.ImageURL, Valid: productToUpdate.ImageURL != ""}).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, productToUpdate.Name, sql.NullString{String: productToUpdate.Description, Valid: true}, fmt.Sprintf("%.2f", productToUpdate.Price), int32(productToUpdate.Stock), sql.NullString{String: productToUpdate.ImageURL, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}))

		err := repo.Update(productToUpdate)
		require.NoError(t, err)
		assert.False(t, productToUpdate.UpdatedAt.IsZero())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		invalidProduct := &productmodule.Product{ID: "invalid-uuid"}
		err := repo.Update(invalidProduct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("database error on update", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE products SET name = $2, description = $3, price = $4, stock = $5, image_url = $6, updated_at = NOW() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productID, productToUpdate.Name, sql.NullString{String: productToUpdate.Description, Valid: productToUpdate.Description != ""}, fmt.Sprintf("%.2f", productToUpdate.Price), int32(productToUpdate.Stock), sql.NullString{String: productToUpdate.ImageURL, Valid: productToUpdate.ImageURL != ""}).WillReturnError(errors.New("db error"))

		err := repo.Update(productToUpdate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	repo := NewProductRepository(db)

	productID := uuid.New()

	t.Run("successful product deletion", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnResult(sqlmock.NewResult(0, 1)) // lastInsertId, rowsAffected

		err := repo.Delete(productID.String())
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		err := repo.Delete("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("database error on delete", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(productID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepository_UpdateStock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	repo := NewProductRepository(db)

	productID := uuid.New()
	quantity := 5

	t.Run("successful stock update", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE products SET stock = stock + $2, updated_at = NOW() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productID, int32(quantity)).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, "Test Product", sql.NullString{String: "Desc", Valid: true}, fmt.Sprintf("%.2f", 10.00), int32(15), sql.NullString{String: "url", Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}))

		err := repo.UpdateStock(productID.String(), quantity)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		err := repo.UpdateStock("invalid-uuid", quantity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("database error on update stock", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE products SET stock = stock + $2, updated_at = NOW() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productID, int32(quantity)).WillReturnError(errors.New("db update stock error"))

		err := repo.UpdateStock(productID.String(), quantity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db update stock error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
