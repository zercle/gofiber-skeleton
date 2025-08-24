package orderhandler

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes order routes with dependency injection
func SetupRoutes(router fiber.Router, jwtHandler fiber.Handler, handler *OrderHandler) {
	// Order management routes
	orderAPI := router.Group("/orders")
	orderAPI.Use(jwtHandler)

	// Get All Orders
	// @Summary Get all orders
	// @Description Retrieves a list of all orders
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Success 200 {array} domain.Order
	// @Router /api/v1/orders [get]
	orderAPI.Get("/", handler.GetAllOrders)

	// Get Order by ID
	// @Summary Get order by ID
	// @Description Retrieves a single order by its ID
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Param id path string true "Order ID"
	// @Success 200 {object} domain.Order
	// @Failure 404 {object} jsend.ErrorResponse "Order not found"
	// @Router /api/v1/orders/{id} [get]
	orderAPI.Get("/:id", handler.GetOrder)

	// Update Order Status
	// @Summary Update order status
	// @Description Updates the status of an order
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Param id path string true "Order ID"
	// @Param status body string true "New order status" Enums(pending, confirmed, shipped, delivered, cancelled)
	// @Success 200 {object} domain.Order
	// @Failure 400 {object} jsend.ErrorResponse "Invalid status"
	// @Failure 404 {object} jsend.ErrorResponse "Order not found"
	// @Router /api/v1/orders/{id}/status [put]
	orderAPI.Put("/:id/status", handler.UpdateOrderStatus)

	// Customer order flow
	// Create Order
	// @Summary Create a new order
	// @Description Creates a new order for a customer
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Param order body object true "Order details"
	// @Success 201 {object} domain.Order
	// @Failure 400 {object} jsend.ErrorResponse "Invalid input"
	// @Router /api/v1/orders/create [post]
	orderAPI.Post("/create", handler.CreateOrder)
}
