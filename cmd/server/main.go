package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/handler"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure"
	"github.com/zercle/gofiber-skeleton/internal/repository"
	"github.com/zercle/gofiber-skeleton/internal/usecase"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	// Initialize database
	dbConfig := infrastructure.NewDatabaseConfig()
	db, err := infrastructure.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize SQLC queries
	queries := infrastructure.NewQueries(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(queries)
	productRepo := repository.NewProductRepository(queries)
	orderRepo := repository.NewOrderRepository(queries)

	// Initialize use cases
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	userUseCase := usecase.NewUserUseCase(userRepo, jwtSecret)
	productUseCase := usecase.NewProductUseCase(productRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, productRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// JWT middleware
	app.Use(infrastructure.JWTMiddleware(jwtSecret))

	// API routes
	api := app.Group("/api/v1")

	// User routes (public)
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)

	// Product routes
	products := api.Group("/products")
	products.Get("/", productHandler.GetAllProducts)                    // Public
	products.Get("/:id", productHandler.GetProduct)                     // Public
	products.Post("/", infrastructure.AdminMiddleware(), productHandler.CreateProduct)      // Admin only
	products.Put("/:id", infrastructure.AdminMiddleware(), productHandler.UpdateProduct)    // Admin only
	products.Delete("/:id", infrastructure.AdminMiddleware(), productHandler.DeleteProduct) // Admin only

	// Order routes
	orders := api.Group("/orders")
	orders.Post("/create", orderHandler.CreateOrder)                    // Authenticated users
	orders.Get("/", orderHandler.GetUserOrders)                        // Authenticated users
	orders.Get("/:id", orderHandler.GetOrder)                          // Authenticated users
	orders.Get("/admin/all", infrastructure.AdminMiddleware(), orderHandler.GetAllOrders)           // Admin only
	orders.Put("/:id/status", infrastructure.AdminMiddleware(), orderHandler.UpdateOrderStatus)    // Admin only

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "E-commerce API is running",
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}