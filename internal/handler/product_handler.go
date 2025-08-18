package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type ProductHandler struct {
	productUseCase domain.ProductUseCase
}

func NewProductHandler(productUseCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

func (ph *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	product := new(domain.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := ph.productUseCase.CreateProduct(product); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create product",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(product)
}

func (ph *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	product, err := ph.productUseCase.GetProductByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(product)
}

func (ph *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := ph.productUseCase.GetAllProducts()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve products",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(products)
}

func (ph *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	product := new(domain.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	product.ID = id

	if err := ph.productUseCase.UpdateProduct(product); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update product",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(product)
}

func (ph *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	if err := ph.productUseCase.DeleteProduct(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete product",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}