package delivery

import (
	"context"
	"gofiber-skeleton/api/order"
	
	"gofiber-skeleton/internal/order/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(app *fiber.App, orderUsecase usecase.OrderUsecase) {
	h := &OrderHandler{
		orderUsecase: orderUsecase,
	}

	// REST Endpoints
	app.Get("/orders/:id", h.GetOrderByID)
	app.Post("/orders", h.CreateOrder)
	app.Put("/orders/:id", h.UpdateOrder)
	app.Delete("/orders/:id", h.DeleteOrder)
}

// REST Handlers
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	// TODO: Implement logic to get order by ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "GetOrderByID"})
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// TODO: Implement logic to create order
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "CreateOrder"})
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	// TODO: Implement logic to update order
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "UpdateOrder"})
}

func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	// TODO: Implement logic to delete order
	return c.SendStatus(fiber.StatusNoContent)
}

// gRPC Server
type GrpcOrderServer struct {
	order.UnimplementedOrderServiceServer
	orderUsecase usecase.OrderUsecase
}

func NewGrpcOrderServer(grpcServer *grpc.Server, orderUsecase usecase.OrderUsecase) {
	s := &GrpcOrderServer{
		orderUsecase: orderUsecase,
	}
	order.RegisterOrderServiceServer(grpcServer, s)
}

// gRPC Handlers
func (s *GrpcOrderServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.Order, error) {
	// TODO: Implement logic to get order
	return &order.Order{Id: req.Id, UserId: 1, ProductId: 1, Quantity: 1, TotalPrice: 10.0}, nil
}

func (s *GrpcOrderServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
	// TODO: Implement logic to create order
	return &order.Order{UserId: req.UserId, ProductId: req.ProductId, Quantity: req.Quantity, TotalPrice: req.TotalPrice}, nil
}

func (s *GrpcOrderServer) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.Order, error) {
	// TODO: Implement logic to update order
	return &order.Order{Id: req.Id, UserId: req.UserId, ProductId: req.ProductId, Quantity: req.Quantity, TotalPrice: req.TotalPrice}, nil
}

func (s *GrpcOrderServer) DeleteOrder(ctx context.Context, req *order.DeleteOrderRequest) (*order.DeleteOrderRequest, error) {
	// TODO: Implement logic to delete order
	return &order.DeleteOrderRequest{Id: req.Id}, nil
}
