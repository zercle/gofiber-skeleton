package delivery

import (
	"context"
	"gofiber-skeleton/api/product"
	
	"gofiber-skeleton/internal/product/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(app *fiber.App, productUsecase usecase.ProductUsecase) {
	h := &ProductHandler{
		productUsecase: productUsecase,
	}

	// REST Endpoints
	app.Get("/products/:id", h.GetProductByID)
	app.Post("/products", h.CreateProduct)
	app.Put("/products/:id", h.UpdateProduct)
	app.Delete("/products/:id", h.DeleteProduct)
}

// REST Handlers
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	// TODO: Implement logic to get product by ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "GetProductByID"})
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// TODO: Implement logic to create product
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "CreateProduct"})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	// TODO: Implement logic to update product
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "UpdateProduct"})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	// TODO: Implement logic to delete product
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "DeleteProduct"})
}

// gRPC Server
type GrpcProductServer struct {
	product.UnimplementedProductServiceServer
	productUsecase usecase.ProductUsecase
}

func NewGrpcProductServer(grpcServer *grpc.Server, productUsecase usecase.ProductUsecase) {
	s := &GrpcProductServer{
		productUsecase: productUsecase,
	}
	product.RegisterProductServiceServer(grpcServer, s)
}

// gRPC Handlers
func (s *GrpcProductServer) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.Product, error) {
	// TODO: Implement logic to get product
	return &product.Product{Id: req.Id, Name: "test", Price: 10.0}, nil
}

func (s *GrpcProductServer) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.Product, error) {
	// TODO: Implement logic to create product
	return &product.Product{Name: req.Name, Price: req.Price}, nil
}

func (s *GrpcProductServer) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.Product, error) {
	// TODO: Implement logic to update product
	return &product.Product{Id: req.Id, Name: req.Name, Price: req.Price}, nil
}

func (s *GrpcProductServer) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductRequest, error) {
	// TODO: Implement logic to delete product
	return &product.DeleteProductRequest{Id: req.Id}, nil
}
