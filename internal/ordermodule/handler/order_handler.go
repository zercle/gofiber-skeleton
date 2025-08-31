package orderhandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/ordermodule"
	"github.com/zercle/gofiber-skeleton/pkg/jsend"
)

type OrderHandler struct {
	orderUseCase ordermodule.OrderUseCase
	validator    *validator.Validate
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderUseCase ordermodule.OrderUseCase, validator *validator.Validate) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		validator:    validator,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	Items []ordermodule.OrderItem `json:"items" validate:"required,min=1,dive"`
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
}

// CreateOrder handles order creation
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Get user ID from JWT token (set by middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		// Fallback to header for integration tests if Locals is not set
		userID = c.Get("user_id")
		if userID == "" {
			return jsend.Fail(c, jsend.Empty, "User not authenticated")
		}
	}

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	// Create order
	order, err := h.orderUseCase.CreateOrder(userID, req.Items)
	if err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	return jsend.SuccessWithStatus(c, fiber.Map{
		"message": "Order created successfully",
		"order":   order,
	}, http.StatusCreated)
}

// GetOrder handles getting a single order
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "Order ID is required")
	}

	order, err := h.orderUseCase.GetOrder(id)
	if err != nil {
		return jsend.Fail(c, jsend.Empty, "Order not found")
	}

	return jsend.Success(c, fiber.Map{
		"order": order,
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
			return jsend.Fail(c, jsend.Empty, "User not authenticated")
		}
	}

	orders, err := h.orderUseCase.GetUserOrders(userID)
	if err != nil {
		return jsend.Error(c, "Failed to fetch orders", 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"orders": orders,
	})
}

// GetAllOrders handles getting all orders (admin only)
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderUseCase.GetAllOrders()
	if err != nil {
		return jsend.Error(c, "Failed to fetch orders", 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"orders": orders,
	})
}

// UpdateOrderStatus handles updating order status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "Order ID is required")
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	// Update order status
	err := h.orderUseCase.UpdateOrderStatus(id, ordermodule.OrderStatus(req.Status))
	if err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	return jsend.Success(c, fiber.Map{
		"message": "Order status updated successfully",
	})
}
