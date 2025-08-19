package handler_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq" // Added this import
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid" // Added this import

	// "github.com/zercle/gofiber-skeleton/internal/domain" // Removed this import
	"github.com/zercle/gofiber-skeleton/internal/handler"
	"github.com/zercle/gofiber-skeleton/internal/repository"
	"github.com/zercle/gofiber-skeleton/internal/repository/db"
	"github.com/zercle/gofiber-skeleton/internal/usecase"
)

var (
	testDB *sql.DB
	app    *fiber.App
)

func TestMain(m *testing.M) {
	// Initialize Fiber app and handlers
	// Note: Database setup and teardown will be handled within individual tests
	// to ensure isolation and a clean state for each test.

	// No need for validate or queries here, as they are not used in this TestMain scope.

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func runMigrations(connStr string) {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database for migrations: %v", err)
	}
	defer dbConn.Close()

	// Ensure UUID extension is enabled
	_, err = dbConn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	// SQL for creating tables
	createUsersTableSQL := `
CREATE TABLE users (
	   id SERIAL PRIMARY KEY,
	   username VARCHAR(255) NOT NULL UNIQUE,
	   password_hash VARCHAR(255) NOT NULL,
	   role VARCHAR(50) NOT NULL DEFAULT 'user',
	   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	   updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

	createProductsTableSQL := `
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);`

	createOrdersTableSQL := `
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    total_price NUMERIC NOT NULL DEFAULT 0.00,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

	createOrderItemsTableSQL := `
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price NUMERIC NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

	// Execute SQL statements in correct order for foreign key dependencies
	_, err = dbConn.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	_, err = dbConn.Exec(createProductsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create products table: %v", err)
	}
	_, err = dbConn.Exec(createOrdersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create orders table: %v", err)
	}
	_, err = dbConn.Exec(createOrderItemsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create order_items table: %v", err)
	}

	log.Println("Migrations applied successfully by direct SQL execution.")
}

func cleanupDatabase(connStr string) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database for cleanup: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`SET client_min_messages TO WARNING;`)
	if err != nil {
		log.Fatalf("Failed to set client_min_messages: %v", err)
	}

	// Get all table names
	rows, err := db.Query(`SELECT tablename FROM pg_tables WHERE schemaname = 'public';`)
	if err != nil {
		log.Fatalf("Failed to query table names: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Failed to scan table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	// Truncate tables to remove data and reset identity sequences
	for _, table := range tables {
		if table == "schema_migrations" {
			continue // Don't truncate the migrations table itself
		}
		_, err := db.Exec(fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY CASCADE;`, table))
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}

	// Drop all tables
	for _, table := range tables {
		if table == "schema_migrations" {
			continue // Don't drop the migrations table itself
		}
		_, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE;`, table))
		if err != nil {
			log.Fatalf("Failed to drop table %s: %v", table, err)
		}
	}

	// Recreate schema_migrations table if it was dropped by CASCADE
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version bigint not null primary key, dirty boolean not null);`)
	if err != nil {
		log.Fatalf("Failed to recreate schema_migrations table: %v", err)
	}

	log.Println("Database cleaned up successfully by dropping all tables.")
}

func TestOrderCreationAndStockUpdate(t *testing.T) {
	// Setup test database for this specific test
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://root:root@localhost:5432/gofiber_boilerplate?sslmode=disable"
	}

	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer testDB.Close()

	// Clean up and run migrations for this test
	cleanupDatabase(connStr)
	runMigrations(connStr)

	// Initialize Fiber app and handlers for this test
	productRepo := repository.NewProductRepository(testDB)
	productUsecase := usecase.NewProductUseCase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	orderRepo := db.NewOrderRepository(testDB)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	app = fiber.New()
	api := app.Group("/api/v1")
	productHandler.RegisterRoutes(api)
	orderHandler.RegisterRoutes(api)

	// Seed a user
	userID := 1 // Assuming auto-incrementing ID for users
	_, err = testDB.Exec(`INSERT INTO users (username, password_hash) VALUES ($1, $2)`, "testuser", "hashedpassword")
	assert.NoError(t, err)

	// Seed a product
	productID := uuid.New()
	productName := "Test Product"
	productDescription := "A product for testing"
	productPrice := 10.00
	productStock := 100
	productImageURL := "http://example.com/image.jpg"

	_, err = testDB.Exec(`INSERT INTO products (id, name, description, price, stock, image_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		productID, productName, productDescription, productPrice, productStock, productImageURL, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create an order
	orderPayload := fmt.Sprintf(`{
		"user_id": "%d",
		"items": [
			{
				"product_id": "%s",
				"quantity": 10
			}
		]
	}`, userID, productID.String())

	req := httptest.NewRequest("POST", "/api/v1/orders/create", bytes.NewBufferString(orderPayload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Verify order details (optional, but good practice)
	// You might want to decode the response body and assert specific fields

	// Verify product stock update
	var updatedStock int
	err = testDB.QueryRow(`SELECT stock FROM products WHERE id = $1`, productID).Scan(&updatedStock)
	assert.NoError(t, err)
	assert.Equal(t, productStock-10, updatedStock)

	// Teardown for this specific test
	cleanupDatabase(connStr)
}