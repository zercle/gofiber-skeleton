package integration

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/productmodule"
	producthandler "github.com/zercle/gofiber-skeleton/internal/productmodule/handler"
	productrepository "github.com/zercle/gofiber-skeleton/internal/productmodule/repository"
	productusecase "github.com/zercle/gofiber-skeleton/internal/productmodule/usecase"
)

func setupProductIntegrationTest(t *testing.T) (*fiber.App, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	productRepo := productrepository.NewProductRepository(db)
	productUseCase := productusecase.NewProductUseCase(productRepo)
	productHandler := producthandler.NewProductHandler(productUseCase, validator.New())

	app := fiber.New()
	app.Post("/api/v1/products", productHandler.CreateProduct)
	app.Get("/api/v1/products/:id", productHandler.GetProduct)
	app.Get("/api/v1/products", productHandler.GetAllProducts)
	app.Put("/api/v1/products/:id", productHandler.UpdateProduct)
	app.Delete("/api/v1/products/:id", productHandler.DeleteProduct)

	return app, mock, db
}

func TestProductIntegration_CreateProduct(t *testing.T) {
	app, mock, db := setupProductIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	productInput := producthandler.CreateProductRequest{
		Name:        "Integration Product",
		Description: "A description for integration test",
		Price:       19.99,
		Stock:       200,
		ImageURL:    "http://example.com/integration.jpg",
	}

	expectedProductUUID := uuid.New()
	mockTime := time.Now()

	t.Run("successful end-to-end product creation", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(expectedProductUUID, productInput.Name, productInput.Description, "19.99", int32(productInput.Stock), productInput.ImageURL, mockTime, mockTime)

		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO products (name, description, price, stock, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productInput.Name, productInput.Description, "19.99", int32(productInput.Stock), productInput.ImageURL).
			WillReturnRows(rows)

		body, _ := json.Marshal(productInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)["product"].(map[string]any)
		assert.Equal(t, productInput.Name, data["name"])
		assert.Equal(t, expectedProductUUID.String(), data["id"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductIntegration_GetProduct(t *testing.T) {
	app, mock, db := setupProductIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	productID := uuid.New()
	expectedProduct := productmodule.Product{
		ID:          productID.String(),
		Name:        "Fetched Product",
		Description: "Fetched description",
		Price:       25.50,
		Stock:       75,
		ImageURL:    "http://example.com/fetched.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful end-to-end product retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, expectedProduct.Name, expectedProduct.Description, fmt.Sprintf("%.2f", expectedProduct.Price), int32(expectedProduct.Stock), expectedProduct.ImageURL, expectedProduct.CreatedAt, expectedProduct.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnRows(rows)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)["product"].(map[string]any)
		assert.Equal(t, expectedProduct.Name, data["name"])
		assert.Equal(t, expectedProduct.ID, data["id"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("product not found end-to-end", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows) // Simulate product not found for GetByID

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Product not found", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductIntegration_GetAllProducts(t *testing.T) {
	app, mock, db := setupProductIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	product1ID := uuid.New()
	product2ID := uuid.New()
	mockTime := time.Now()

	t.Run("successful end-to-end get all products", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(product1ID, "Product One", "Desc One", "10.00", 10, "url1", mockTime, mockTime).
			AddRow(product2ID, "Product Two", "Desc Two", "20.00", 20, "url2", mockTime, mockTime)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products ORDER BY created_at DESC`,
		)).
			WillReturnRows(rows)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)["products"].([]any)
		assert.Len(t, data, 2)
		assert.Equal(t, product1ID.String(), data[0].(map[string]any)["id"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductIntegration_UpdateProduct(t *testing.T) {
	app, mock, db := setupProductIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	productID := uuid.New()
	updateInput := producthandler.UpdateProductRequest{
		Name:        "Updated Name",
		Description: "Updated Desc",
		Price:       30.00,
		Stock:       30,
		ImageURL:    "http://example.com/updated.jpg",
	}
	mockTime := time.Now()

	t.Run("successful end-to-end product update", func(t *testing.T) {
		// Expect GetByID call first within the usecase
		originalProductRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, "Original Name", "Original Desc", "10.00", 10, "original.jpg", mockTime, mockTime)
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnRows(originalProductRows)

		// Expect Update call
		updatedProductRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "image_url", "created_at", "updated_at"}).
			AddRow(productID, updateInput.Name, updateInput.Description, fmt.Sprintf("%.2f", updateInput.Price), int32(updateInput.Stock), updateInput.ImageURL, mockTime, mockTime.Add(time.Second)) // Simulate updated_at
		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE products SET name = $2, description = $3, price = $4, stock = $5, image_url = $6, updated_at = NOW() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, created_at, updated_at`,
		)).
			WithArgs(productID, updateInput.Name, updateInput.Description, fmt.Sprintf("%.2f", updateInput.Price), int32(updateInput.Stock), updateInput.ImageURL).
			WillReturnRows(updatedProductRows)

		body, _ := json.Marshal(updateInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)["product"].(map[string]any)
		assert.Equal(t, updateInput.Name, data["name"])
		assert.Equal(t, productID.String(), data["id"]) // Added assertion for ID
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("product not found end-to-end update", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, price, stock, image_url, created_at, updated_at FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows) // Simulate product not found for GetByID

		body, _ := json.Marshal(updateInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Product not found", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductIntegration_DeleteProduct(t *testing.T) {
	productID := uuid.New()

	t.Run("successful end-to-end product deletion", func(t *testing.T) {
		app, mock, db := setupProductIntegrationTest(t)
		defer func() {
			_ = db.Close()
		}()

		// Expect Delete call, which returns rows affected
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		assert.Equal(t, "Product deleted successfully", responseBody["data"].(map[string]any)["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("product not found end-to-end delete", func(t *testing.T) {
		app, mock, db := setupProductIntegrationTest(t)
		defer func() {
			_ = db.Close()
		}()

		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM products WHERE id = $1`,
		)).
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows) // Simulate product not found

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Product not found", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
