package delivery

import (
	"context"
	"gofiber-skeleton/api/product"
	"gofiber-skeleton/pkg/jsend"

	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/infra/middleware"
	"gofiber-skeleton/internal/product/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(app *fiber.App, uc usecase.ProductUsecase, jwtService auth.JWTService) {
	handler := &ProductHandler{
		productUsecase: uc,
	}

	// Group all product routes and protect them
	productsAPI := app.Group("/api/v1/products", middleware.Protected(jwtService))

	productsAPI.Get("/", handler.GetAllProducts)
	productsAPI.Get("/:id", handler.GetProductByID)
	productsAPI.Post("/", handler.CreateProduct)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	// products, err := h.productUsecase.FindAll(c.Context())
	// if err != nil {
	// 	return jsend.Error(c, "Could not retrieve products", http.StatusInternalServerError)
	// }
	products := []fiber.Map{{"id": "prod_123", "name": "Sample Product"}}
	return jsend.Success(c, products)
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
	return c.SendStatus(fiber.StatusNoContent)
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

func (s *GrpcProductServer) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*emptypb.Empty, error) {
	// TODO: Implement logic to delete product
	return &emptypb.Empty{}, nil
}
