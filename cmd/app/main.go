package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	orderDelivery "gofiber-skeleton/internal/order/delivery"
	orderInfrastructure "gofiber-skeleton/internal/order/infrastructure"
	orderUsecase "gofiber-skeleton/internal/order/usecase"
	productDelivery "gofiber-skeleton/internal/product/delivery"
	productInfrastructure "gofiber-skeleton/internal/product/infrastructure"
	productUsecase "gofiber-skeleton/internal/product/usecase"
	userDelivery "gofiber-skeleton/internal/user/delivery"
	userInfrastructure "gofiber-skeleton/internal/user/infrastructure"
	userUsecase "gofiber-skeleton/internal/user/usecase"
	"gofiber-skeleton/pkg/app"
	"gofiber-skeleton/pkg/config"
	"gofiber-skeleton/pkg/database"

	"google.golang.org/grpc"
)

func main() {
	migrationsPath := os.Getenv("PWD") + "/database/migrations"
	fiberApp := app.SetupApp(migrationsPath)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db := database.ConnectDB(cfg.DATABASE_URL)

	// Initialize User Module
	userRepo := userInfrastructure.NewUserRepository(db)
	userUc := userUsecase.NewUserUsecase(userRepo)
	userDelivery.NewUserHandler(fiberApp, userUc)

	// Initialize Product Module
	productRepo := productInfrastructure.NewProductRepository(db)
	productUc := productUsecase.NewProductUsecase(productRepo)
	productDelivery.NewProductHandler(fiberApp, productUc)

	// Initialize Order Module
	orderRepo := orderInfrastructure.NewOrderRepository(db)
	orderUc := orderUsecase.NewOrderUsecase(orderRepo)
	orderDelivery.NewOrderHandler(fiberApp, orderUc)

	// Setup gRPC Server
	grpcPort := "50051" // You can make this configurable
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userDelivery.NewGrpcUserServer(grpcServer, userUc)
	productDelivery.NewGrpcProductServer(grpcServer, productUc)
	orderDelivery.NewGrpcOrderServer(grpcServer, orderUc)

	go func() {
		log.Printf("gRPC server listening on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start Fiber app in a goroutine
	go func() {
		if err := fiberApp.Listen(fmt.Sprintf(":%s", cfg.APP_PORT)); err != nil {
			log.Fatalf("failed to serve Fiber: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Shutdown gRPC server
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")

	// Shutdown Fiber app
	if err := fiberApp.Shutdown(); err != nil {
		log.Fatalf("Fiber app shutdown error: %v", err)
	}
	log.Println("Fiber app stopped.")

	log.Println("Servers gracefully stopped.")
}
