package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type ProductHandler struct {
	productUseCase domain.ProductUseCase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUseCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create product
	product, err := h.productUseCase.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.ImageURL)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": product,
	})
}

// GetProduct handles getting a single product
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	product, err := h.productUseCase.GetProduct(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(fiber.Map{
		"product": product,
	})
}

// GetAllProducts handles getting all products
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productUseCase.GetAllProducts()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
	})
}

// UpdateProduct handles product updates
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update product
	product, err := h.productUseCase.UpdateProduct(id, req.Name, req.Description, req.Price, req.Stock, req.ImageURL)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

// DeleteProduct handles product deletion
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	if err := h.productUseCase.DeleteProduct(id); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}