package producthandler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/productmodule"
	productmock "github.com/zercle/gofiber-skeleton/internal/productmodule/mock"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := productmock.NewMockProductUseCase(ctrl)
	app := fiber.New()
	handler := NewProductHandler(mockProductUseCase, validator.New())
	app.Post("/api/v1/products", handler.CreateProduct)

	productInput := CreateProductRequest{
		Name:        "Test Product",
		Description: "A description",
		Price:       10.99,
		Stock:       100,
		ImageURL:    "http://example.com/image.jpg",
	}

	t.Run("successful product creation", func(t *testing.T) {
		expectedProduct := &productmodule.Product{
			ID:          uuid.New().String(),
			Name:        productInput.Name,
			Description: productInput.Description,
			Price:       productInput.Price,
			Stock:       productInput.Stock,
			ImageURL:    productInput.ImageURL,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockProductUseCase.EXPECT().CreateProduct(productInput.Name, productInput.Description, productInput.Price, productInput.Stock, productInput.ImageURL).Return(expectedProduct, nil)

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
		assert.Equal(t, expectedProduct.Name, data["name"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader([]byte(`{"name":123}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Invalid request body", responseBody["data"].(map[string]any)["message"])
	})

	t.Run("validation errors", func(t *testing.T) {
		invalidInput := CreateProductRequest{
			Name:  "", // Empty name
			Price: -10.00,
		}
		body, _ := json.Marshal(invalidInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Contains(t, responseBody["data"].(map[string]any)["message"].(string), "required")
		assert.Contains(t, responseBody["data"].(map[string]any)["message"].(string), "min")
	})

	t.Run("error from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().CreateProduct(productInput.Name, productInput.Description, productInput.Price, productInput.Stock, productInput.ImageURL).Return(nil, errors.New("internal server error"))

		body, _ := json.Marshal(productInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "internal server error", responseBody["message"])
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := productmock.NewMockProductUseCase(ctrl)
	app := fiber.New()
	handler := NewProductHandler(mockProductUseCase, validator.New())
	app.Get("/api/v1/products/:id", handler.GetProduct)

	productID := uuid.New()
	expectedProduct := &productmodule.Product{
		ID:          productID.String(),
		Name:        "Test Product",
		Description: "A description",
		Price:       10.99,
		Stock:       100,
		ImageURL:    "http://example.com/image.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("successful product retrieval", func(t *testing.T) {
		mockProductUseCase.EXPECT().GetProduct(productID.String()).Return(expectedProduct, nil)

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
	})

	t.Run("invalid product ID", func(t *testing.T) {
		// Use uuidv4.Parse for validation of existing UUIDs if needed, or remove if all are UUIDv7
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/invalid-uuid", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Invalid product ID format", responseBody["data"].(map[string]any)["message"])
	})

	t.Run("product not found", func(t *testing.T) {
		mockProductUseCase.EXPECT().GetProduct(productID.String()).Return(nil, productmodule.ErrProductNotFound)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Product not found", responseBody["message"])
	})

	t.Run("error from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().GetProduct(productID.String()).Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "database error", responseBody["message"])
	})
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := productmock.NewMockProductUseCase(ctrl)
	app := fiber.New()
	handler := NewProductHandler(mockProductUseCase, validator.New())
	app.Get("/api/v1/products", handler.GetAllProducts)

	t.Run("successful retrieval of all products", func(t *testing.T) {
		expectedProducts := []*productmodule.Product{
			{ID: uuid.New().String(), Name: "Product 1", Price: 10.00, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: uuid.New().String(), Name: "Product 2", Price: 20.00, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockProductUseCase.EXPECT().GetAllProducts().Return(expectedProducts, nil)

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
		assert.Equal(t, expectedProducts[0].Name, data[0].(map[string]any)["name"])
	})

	t.Run("error from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().GetAllProducts().Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "Failed to fetch products", responseBody["message"])
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := productmock.NewMockProductUseCase(ctrl)
	app := fiber.New()
	handler := NewProductHandler(mockProductUseCase, validator.New())
	app.Put("/api/v1/products/:id", handler.UpdateProduct)

	productID := uuid.New()
	productInput := UpdateProductRequest{
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       15.99,
		Stock:       50,
		ImageURL:    "http://example.com/updated_image.jpg",
	}

	t.Run("successful product update", func(t *testing.T) {
		expectedProduct := &productmodule.Product{
			ID:          productID.String(),
			Name:        productInput.Name,
			Description: productInput.Description,
			Price:       productInput.Price,
			Stock:       productInput.Stock,
			ImageURL:    productInput.ImageURL,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockProductUseCase.EXPECT().UpdateProduct(productID.String(), productInput.Name, productInput.Description, productInput.Price, productInput.Stock, productInput.ImageURL).Return(expectedProduct, nil)

		body, _ := json.Marshal(productInput)
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
		assert.Equal(t, expectedProduct.Name, data["name"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+productID.String(), bytes.NewReader([]byte(`{"name":123}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Invalid request body", responseBody["data"].(map[string]any)["message"])
	})

	t.Run("validation errors", func(t *testing.T) {
		invalidInput := UpdateProductRequest{
			Name:  "", // Empty name
			Price: -10.00,
		}
		body, _ := json.Marshal(invalidInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Contains(t, responseBody["data"].(map[string]any)["message"].(string), "min")
	})

	t.Run("product not found from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().UpdateProduct(productID.String(), productInput.Name, productInput.Description, productInput.Price, productInput.Stock, productInput.ImageURL).Return(nil, productmodule.ErrProductNotFound)

		body, _ := json.Marshal(productInput)
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
	})

	t.Run("error from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().UpdateProduct(productID.String(), productInput.Name, productInput.Description, productInput.Price, productInput.Stock, productInput.ImageURL).Return(nil, errors.New("internal server error"))

		body, _ := json.Marshal(productInput)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "internal server error", responseBody["message"])
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := productmock.NewMockProductUseCase(ctrl)
	app := fiber.New()
	handler := NewProductHandler(mockProductUseCase, validator.New())
	app.Delete("/api/v1/products/:id", handler.DeleteProduct)

	productID := uuid.New()

	t.Run("successful product deletion", func(t *testing.T) {
		mockProductUseCase.EXPECT().DeleteProduct(productID.String()).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		assert.Equal(t, "Product deleted successfully", responseBody["data"].(map[string]any)["message"])
	})

	t.Run("product not found from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().DeleteProduct(productID.String()).Return(productmodule.ErrProductNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Product not found", responseBody["message"])
	})

	t.Run("error from usecase", func(t *testing.T) {
		mockProductUseCase.EXPECT().DeleteProduct(productID.String()).Return(errors.New("internal server error"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+productID.String(), nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "internal server error", responseBody["message"])
	})
}
