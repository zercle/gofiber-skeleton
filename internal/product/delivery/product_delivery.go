package delivery

import (
	"context"
	"gofiber-skeleton/api/product"

	"gofiber-skeleton/internal/product/usecase"
	"gofiber-skeleton/internal/infra/jsend"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductDelivery struct {
	productUsecase usecase.ProductUsecase
}

func NewProductDelivery(productUsecase usecase.ProductUsecase) *ProductDelivery {
	return &ProductDelivery{
		productUsecase: productUsecase,
	}
}

// REST Handlers
// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} jsend.JSendResponse{data=domain.Product} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /products/{id} [get]
func (d *ProductDelivery) GetProductByID(c *fiber.Ctx) error {
	// TODO: Implement logic to get product by ID
	return jsend.Success(c, fiber.Map{"message": "GetProductByID"})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product object"
// @Success 201 {object} jsend.JSendResponse{data=domain.Product} "Created"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /products [post]
func (d *ProductDelivery) CreateProduct(c *fiber.Ctx) error {
	// TODO: Implement logic to create product
	return jsend.Success(c, fiber.Map{"message": "CreateProduct"})
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.Product true "Product object"
// @Success 200 {object} jsend.JSendResponse{data=domain.Product} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /products/{id} [put]
func (d *ProductDelivery) UpdateProduct(c *fiber.Ctx) error {
	// TODO: Implement logic to update product
	return jsend.Success(c, fiber.Map{"message": "UpdateProduct"})
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} jsend.JSendResponse{data=string} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /products/{id} [delete]
func (d *ProductDelivery) DeleteProduct(c *fiber.Ctx) error {
	// TODO: Implement logic to delete product
	return jsend.Success(c, fiber.Map{"message": "DeleteProduct"})
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