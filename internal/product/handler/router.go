package producthandler

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes product routes with dependency injection
func SetupRoutes(router fiber.Router, jwtHandler fiber.Handler, handler *ProductHandler) {
	// Product management routes
	productAPI := router.Group("/products")
	productAPI.Use(jwtHandler)

	// Create Product
	// @Summary Create a new product
	// @Description Adds a new product to the inventory
	// @Tags products
	// @Accept json
	// @Produce json
	// @Param product body object true "Product details"
	// @Success 201 {object} domain.Product
	// @Failure 400 {object} jsend.ErrorResponse "Invalid input"
	// @Router /api/v1/products [post]
	productAPI.Post("/", handler.CreateProduct)

	// Update Product
	// @Summary Update product information
	// @Description Updates an existing product's details
	// @Tags products
	// @Accept json
	// @Produce json
	// @Param id path string true "Product ID"
	// @Param product body object true "Updated product details"
	// @Success 200 {object} domain.Product
	// @Failure 400 {object} jsend.ErrorResponse "Invalid input"
	// @Failure 404 {object} jsend.ErrorResponse "Product not found"
	// @Router /api/v1/products/{id} [put]
	productAPI.Put("/:id", handler.UpdateProduct)

	// Delete Product
	// @Summary Delete a product
	// @Description Deletes a product by its ID
	// @Tags products
	// @Accept json
	// @Produce json
	// @Param id path string true "Product ID"
	// @Success 204 "No Content"
	// @Failure 404 {object} jsend.ErrorResponse "Product not found"
	// @Router /api/v1/products/{id} [delete]
	productAPI.Delete("/:id", handler.DeleteProduct)

	// Get All Products
	// @Summary Get all products
	// @Description Retrieves a list of all products
	// @Tags products
	// @Accept json
	// @Produce json
	// @Success 200 {array} domain.Product
	// @Router /api/v1/products [get]
	productAPI.Get("/", handler.GetAllProducts)

	// Get Product by ID
	// @Summary Get product by ID
	// @Description Retrieves a single product by its ID
	// @Tags products
	// @Accept json
	// @Produce json
	// @Param id path string true "Product ID"
	// @Success 200 {object} domain.Product
	// @Failure 404 {object} jsend.ErrorResponse "Product not found"
	// @Router /api/v1/products/{id} [get]
	productAPI.Get("/:id", handler.GetProduct)
}
