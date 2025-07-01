package delivery

import (
	"context"
	"gofiber-skeleton/api/order"

	"gofiber-skeleton/internal/order/usecase"
	"gofiber-skeleton/internal/infra/jsend"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderDelivery struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderDelivery(orderUsecase usecase.OrderUsecase) *OrderDelivery {
	return &OrderDelivery{
		orderUsecase: orderUsecase,
	}
}

// REST Handlers
// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} jsend.JSendResponse{data=domain.Order} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /orders/{id} [get]
func (d *OrderDelivery) GetOrderByID(c *fiber.Ctx) error {
	// TODO: Implement logic to get order by ID
	return jsend.Success(c, fiber.Map{"message": "GetOrderByID"})
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body domain.Order true "Order object"
// @Success 201 {object} jsend.JSendResponse{data=domain.Order} "Created"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /orders [post]
func (d *OrderDelivery) CreateOrder(c *fiber.Ctx) error {
	// TODO: Implement logic to create order
	return jsend.Success(c, fiber.Map{"message": "CreateOrder"})
}

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body domain.Order true "Order object"
// @Success 200 {object} jsend.JSendResponse{data=domain.Order} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /orders/{id} [put]
func (d *OrderDelivery) UpdateOrder(c *fiber.Ctx) error {
	// TODO: Implement logic to update order
	return jsend.Success(c, fiber.Map{"message": "UpdateOrder"})
}

// DeleteOrder godoc
// @Summary Delete an order by ID
// @Description Delete an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} jsend.JSendResponse{data=string} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /orders/{id} [delete]
func (d *OrderDelivery) DeleteOrder(c *fiber.Ctx) error {
	// TODO: Implement logic to delete order
	return jsend.Success(c, fiber.Map{"message": "DeleteOrder"})
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
	return &order.Order{Id: req.Id, UserId: 0, ProductId: 0, Quantity: 1, Status: "pending"}, nil
}

func (s *GrpcOrderServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
	// TODO: Implement logic to create order
	return &order.Order{UserId: req.UserId, ProductId: req.ProductId, Quantity: req.Quantity, TotalPrice: req.TotalPrice, Status: "pending"}, nil
}

func (s *GrpcOrderServer) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.Order, error) {
	// TODO: Implement logic to update order
	return &order.Order{Id: req.Id, UserId: req.UserId, ProductId: req.ProductId, Quantity: req.Quantity, TotalPrice: req.TotalPrice, Status: req.Status}, nil
}

func (s *GrpcOrderServer) DeleteOrder(ctx context.Context, req *order.DeleteOrderRequest) (*emptypb.Empty, error) {
	// TODO: Implement logic to delete order
	return &emptypb.Empty{}, nil
}