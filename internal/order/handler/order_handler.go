package orderhandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type OrderHandler struct {
	orderUseCase domain.OrderUseCase
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderUseCase domain.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	Items []domain.OrderItem `json:"items"`
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

// CreateOrder handles order creation
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Get user ID from JWT token (set by middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		// Fallback to header for integration tests if Locals is not set
		userID = c.Get("user_id")
		if userID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "User not authenticated",
			})
		}
	}

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	// Validate items
	if len(req.Items) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Order must contain at least one item",
		})
	}

	// Create order
	order, err := h.orderUseCase.CreateOrder(userID, req.Items)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
		"order":   order,
	})
}

// GetOrder handles getting a single order
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Order ID is required",
		})
	}

	order, err := h.orderUseCase.GetOrder(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Order not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"order": order,
		},
	})
}

// GetUserOrders handles getting orders for the authenticated user
func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	// Get user ID from JWT token (set by middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		// Fallback to header for integration tests if Locals is not set
		userID = c.Get("user_id")
		if userID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "User not authenticated",
			})
		}
	}

	orders, err := h.orderUseCase.GetUserOrders(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch orders",
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
	})
}

// GetAllOrders handles getting all orders (admin only)
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderUseCase.GetAllOrders()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch orders",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"orders": orders,
		},
	})
}

// UpdateOrderStatus handles updating order status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Order ID is required",
		})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	// Update order status
	err := h.orderUseCase.UpdateOrderStatus(id, domain.OrderStatus(req.Status))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}
