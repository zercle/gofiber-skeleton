package handler

import (
	"log" // Added this import
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid" // Added this import
	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/usecase"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (oh *OrderHandler) RegisterRoutes(api fiber.Router) {
	orderRoutes := api.Group("/orders")
	orderRoutes.Post("/create", oh.CreateOrder)
	orderRoutes.Get("/", oh.ListOrders)
	orderRoutes.Get("/:id", oh.GetOrderByID)
	orderRoutes.Put("/:id/status", oh.UpdateOrderStatus)
}

type CreateOrderPayload struct {
	UserID string `json:"user_id"`
	Items  []struct {
		ProductID uuid.UUID `json:"product_id"` // Changed to uuid.UUID
		Quantity  int64     `json:"quantity"`
	} `json:"items"`
}

func (oh *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var payload CreateOrderPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Error parsing request body: %v", err) // Add logging
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	orderItems := make([]domain.OrderItem, len(payload.Items))
	for i, item := range payload.Items {
		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			// Price will be set in usecase after fetching product details
		}
	}

	createdOrder, err := oh.orderUsecase.ProcessOrder(c.Context(), payload.UserID, orderItems)
	if err != nil {
		log.Printf("Error processing order: %v", err) // Add logging
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create order",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(createdOrder)
}

func (oh *OrderHandler) ListOrders(c *fiber.Ctx) error {
	orders, err := oh.orderUsecase.ListOrders(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve orders",
			"error":   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(orders)
}

func (oh *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
			"error":   err.Error(),
		})
	}

	order, err := oh.orderUsecase.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
			"error":   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(order)
}

func (oh *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
			"error":   err.Error(),
		})
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	order, err := oh.orderUsecase.UpdateOrderStatus(c.Context(), id, domain.OrderStatus(payload.Status))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update order status",
			"error":   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(order)
}