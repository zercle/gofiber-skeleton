package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"

	demousecase "github.com/zercle/gofiber-skeleton/internal/demo/usecase"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/fiber/response"
)

// DemoHandler handles HTTP requests for demo operations.
type DemoHandler struct {
	cfg         *config.Config
	demoUsecase demousecase.DemoUseCase
}

// NewDemoHandler creates a new DemoHandler.
func NewDemoHandler(cfg *config.Config, demoUsecase demousecase.DemoUseCase) *DemoHandler {
	return &DemoHandler{
		cfg:         cfg,
		demoUsecase: demoUsecase,
	}
}

// PerformTransactionDemo
// @Summary Perform a transaction demonstration
// @Description Demonstrates a database transaction by creating, updating, and deleting a product.
// @Tags Demo
// @Accept json
// @Produce json
// @Success 200 {object} response.JSendSuccess "Transaction demo successful"
// @Failure 500 {object} response.JSendFail "Internal Server Error"
// @Router /demo/transaction [post]
func (h *DemoHandler) PerformTransactionDemo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.Timeout)
	defer cancel()

	err := h.demoUsecase.PerformTransactionDemo(ctx)
	if err != nil {
		return response.Fail(c, "internal_server_error", fiber.Map{"error": err.Error()}, http.StatusInternalServerError)
	}

	return response.Success(c, fiber.Map{"message": "Transaction demo successful"})
}

// GetJoinedDataDemo
// @Summary Get joined data demonstration
// @Description Retrieves joined data from orders, order items, and products for demonstration.
// @Tags Demo
// @Accept json
// @Produce json
// @Success 200 {object} response.JSendSuccess{data=[]sqlc.GetOrdersWithItemsAndProductsRow} "Joined data retrieved successfully"
// @Failure 500 {object} response.JSendFail "Internal Server Error"
// @Router /demo/joined [get]
func (h *DemoHandler) GetJoinedDataDemo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.cfg.App.Timeout)
	defer cancel()

	data, err := h.demoUsecase.GetJoinedDataDemo(ctx)
	if err != nil {
		return response.Fail(c, "internal_server_error", fiber.Map{"error": err.Error()}, http.StatusInternalServerError)
	}

	return response.Success(c, fiber.Map{"data": data})
}
