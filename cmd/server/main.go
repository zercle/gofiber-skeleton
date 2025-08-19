package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3" // Renamed to avoid conflict with "github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"               // PostgreSQL driver
	"github.com/zercle/gofiber-skeleton/internal/handler"
	"github.com/zercle/gofiber-skeleton/internal/repository"
	sqldb "github.com/zercle/gofiber-skeleton/internal/repository/db"
	"github.com/zercle/gofiber-skeleton/internal/usecase"
)

func main() {
	// Database connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Initialize repositories and use cases
	productRepository := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productHandler := handler.NewProductHandler(productUseCase)

	userRepository := sqldb.NewUserRepository(db)
	userUseCase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	orderRepository := sqldb.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUsecase(orderRepository, productUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	userHandler.RegisterRoutes(api)

	// JWT Middleware for protected routes
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	// Product routes
	productRoutes := api.Group("/products")
	productRoutes.Post("/", productHandler.CreateProduct)
	productRoutes.Get("/", productHandler.GetAllProducts)
	productRoutes.Get("/:id", productHandler.GetProductByID)
	productRoutes.Put("/:id", productHandler.UpdateProduct)
	productRoutes.Delete("/:id", productHandler.DeleteProduct)

	// Order routes
	orderHandler.RegisterRoutes(api)

	log.Fatal(app.Listen(":3000"))
}