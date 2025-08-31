package producthandler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/productmodule"
	"github.com/zercle/gofiber-skeleton/pkg/jsend"
)

// ProductHandler handles HTTP requests related to products.
type ProductHandler struct {
	productUseCase productmodule.ProductUseCase
	validator      *validator.Validate
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(productUseCase productmodule.ProductUseCase, validator *validator.Validate) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		validator:      validator,
	}
}

// CreateProductRequest represents the request body for creating a product.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	ImageURL    string  `json:"image_url" validate:"url"`
}

// UpdateProductRequest represents the request body for updating a product.
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"min=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	ImageURL    string  `json:"image_url" validate:"url"`
}

// CreateProduct handles the creation of a new product via HTTP POST request.
// It parses the request body, validates the input, and calls the product use case to create the product.
// Returns a 201 Created status with the new product on success, or an error status on failure.
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	product, err := h.productUseCase.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.ImageURL)
	if err != nil {
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.SuccessWithStatus(c, fiber.Map{
		"message": "Product created successfully",
		"product": product,
	}, http.StatusCreated)
}

// GetProduct handles retrieving a single product by its ID via HTTP GET request.
// It extracts the product ID from the URL parameters, validates it, and calls the product use case to fetch the product.
// Returns a 200 OK status with the product details on success, a 404 Not Found if the product does not exist,
// or other error statuses on failure.
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "Product ID is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid product ID format")
	}

	product, err := h.productUseCase.GetProduct(id)
	if err != nil {
		if errors.Is(err, productmodule.ErrProductNotFound) {
			return jsend.Fail(c, jsend.Empty, "Product not found")
		}
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"product": product,
	})
}

// GetAllProducts handles retrieving all products via HTTP GET request.
// It calls the product use case to fetch all products.
// Returns a 200 OK status with a list of products on success, or an error status on failure.
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productUseCase.GetAllProducts()
	if err != nil {
		return jsend.Error(c, "Failed to fetch products", 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"products": products,
	})
}

// UpdateProduct handles updating an existing product by its ID via HTTP PUT request.
// It extracts the product ID from the URL parameters, parses and validates the request body,
// and calls the product use case to update the product.
// Returns a 200 OK status with the updated product on success, a 404 Not Found if the product does not exist,
// or other error statuses on failure.
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "Product ID is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid product ID format")
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	product, err := h.productUseCase.UpdateProduct(id, req.Name, req.Description, req.Price, req.Stock, req.ImageURL)
	if err != nil {
		if errors.Is(err, productmodule.ErrProductNotFound) {
			return jsend.Fail(c, jsend.Empty, "Product not found")
		}
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

// DeleteProduct handles deleting a product by its ID via HTTP DELETE request.
// It extracts the product ID from the URL parameters, validates it, and calls the product use case to delete the product.
// Returns a 200 OK status on successful deletion, a 404 Not Found if the product does not exist,
// or other error statuses on failure.
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "Product ID is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid product ID format")
	}

	if err := h.productUseCase.DeleteProduct(id); err != nil {
		if errors.Is(err, productmodule.ErrProductNotFound) {
			return jsend.Fail(c, jsend.Empty, "Product not found")
		}
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"message": "Product deleted successfully",
	})
}

// UpdateStock handles updating the stock quantity of a product via HTTP PUT request.
// It extracts the product ID from the URL parameters, parses and validates the request body,
// and calls the product use case to update the product's stock.
// Returns a 200 OK status with the updated product on success, a 404 Not Found if the product does not exist,
// or other error statuses on failure.
func (h *ProductHandler) UpdateStock(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "Product ID is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid product ID format")
	}

	var req struct {
		Stock int `json:"stock" validate:"min=0"`
	}
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	if err := h.productUseCase.UpdateStock(id, req.Stock); err != nil {
		if errors.Is(err, productmodule.ErrProductNotFound) {
			return jsend.Fail(c, jsend.Empty, "Product not found")
		}
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"message": "Product stock updated successfully",
	})
}
