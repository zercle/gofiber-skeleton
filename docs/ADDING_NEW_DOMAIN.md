# Adding a New Domain

This guide demonstrates how to add a new business domain to the project by following the established architecture patterns.

## Overview

Each domain follows Clean Architecture principles with these layers:

1. **Entity**: Domain models and DTOs
2. **Repository**: Data access interfaces and implementations
3. **Usecase**: Business logic and workflows
4. **Delivery**: HTTP handlers and routes
5. **Tests**: Comprehensive test coverage
6. **Mocks**: Generated test mocks

## Step-by-Step Guide

Let's create a "Product" domain as an example.

### Step 1: Create Directory Structure

```bash
mkdir -p internal/domains/product/{entity,repository,usecase,delivery,tests,mocks}
```

Your structure should look like:
```
internal/domains/product/
├── entity/
├── repository/
├── usecase/
├── delivery/
├── tests/
└── mocks/
```

### Step 2: Define Entities

Create `internal/domains/product/entity/product.go`:

```go
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product entity
type Product struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CategoryID  uuid.UUID  `json:"category_id"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateProductRequest represents a product creation request
type CreateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description *string    `json:"description,omitempty"`
	Price       float64    `json:"price" validate:"required,gte=0"`
	Stock       int        `json:"stock" validate:"required,gte=0"`
	CategoryID  uuid.UUID  `json:"category_id" validate:"required"`
}

// UpdateProductRequest represents a product update request
type UpdateProductRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string    `json:"description,omitempty"`
	Price       *float64   `json:"price,omitempty" validate:"omitempty,gte=0"`
	Stock       *int       `json:"stock,omitempty" validate:"omitempty,gte=0"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
}

// ProductListResponse represents a paginated product list
type ProductListResponse struct {
	Products   []*Product `json:"products"`
	TotalCount int64      `json:"total_count"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
}
```

### Step 3: Create Database Migration

```bash
make migrate-create NAME=create_products_table
```

Edit the generated up migration (`db/migrations/XXXXXX_create_products_table.up.sql`):

```sql
-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_name ON products(name);

-- Create updated_at trigger
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

Edit the down migration (`db/migrations/XXXXXX_create_products_table.down.sql`):

```sql
-- Drop trigger
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

-- Drop indexes
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_category_id;

-- Drop table
DROP TABLE IF EXISTS products;
```

Apply the migration:
```bash
make migrate-up
```

### Step 4: Define SQL Queries

Create `db/queries/products.sql`:

```sql
-- name: CreateProduct :one
INSERT INTO products (
    name, description, price, stock, category_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE is_active = TRUE;

-- name: UpdateProduct :one
UPDATE products
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    price = COALESCE(sqlc.narg('price'), price),
    stock = COALESCE(sqlc.narg('stock'), stock),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: UpdateProductStock :exec
UPDATE products
SET stock = $2, updated_at = NOW()
WHERE id = $1;
```

Generate SQL code:
```bash
make sqlc
```

### Step 5: Define Repository Interface

Create `internal/domains/product/repository/repository.go`:

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-username/your-project/internal/domains/product/entity"
)

//go:generate mockgen -source=repository.go -destination=../mocks/repository_mock.go -package=mocks

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Create creates a new product
	Create(ctx context.Context, product *entity.Product) error

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)

	// List retrieves a paginated list of products
	List(ctx context.Context, limit, offset int) ([]*entity.Product, error)

	// Count returns the total number of active products
	Count(ctx context.Context) (int64, error)

	// Update updates a product
	Update(ctx context.Context, product *entity.Product) error

	// Delete soft deletes a product
	Delete(ctx context.Context, productID uuid.UUID) error

	// UpdateStock updates product stock
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error
}
```

### Step 6: Implement Repository

Create `internal/domains/product/repository/postgres_repository.go`:

```go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/your-project/internal/domains/product/entity"
)

// PostgresProductRepository implements ProductRepository for PostgreSQL
type PostgresProductRepository struct {
	db *pgxpool.Pool
}

// NewPostgresProductRepository creates a new PostgreSQL product repository
func NewPostgresProductRepository(db *pgxpool.Pool) ProductRepository {
	return &PostgresProductRepository{db: db}
}

// Create creates a new product
func (r *PostgresProductRepository) Create(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
	).Scan(
		&product.ID,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *PostgresProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category_id, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &entity.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// List retrieves a paginated list of products
func (r *PostgresProductRepository) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category_id, is_active, created_at, updated_at
		FROM products
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product := &entity.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// Count returns the total number of active products
func (r *PostgresProductRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM products WHERE is_active = TRUE`

	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// Update updates a product
func (r *PostgresProductRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, price = $4, stock = $5,
		    category_id = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// Delete soft deletes a product
func (r *PostgresProductRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	query := `UPDATE products SET is_active = FALSE, updated_at = NOW() WHERE id = $1`

	_, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// UpdateStock updates product stock
func (r *PostgresProductRepository) UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error {
	query := `UPDATE products SET stock = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.Exec(ctx, query, productID, stock)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
```

### Step 7: Define Usecase Interface

Create `internal/domains/product/usecase/usecase.go`:

```go
package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-username/your-project/internal/domains/product/entity"
)

//go:generate mockgen -source=usecase.go -destination=../mocks/usecase_mock.go -package=mocks

// ProductUsecase defines the interface for product business logic
type ProductUsecase interface {
	// Create creates a new product
	Create(ctx context.Context, req *entity.CreateProductRequest) (*entity.Product, error)

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)

	// List retrieves a paginated list of products
	List(ctx context.Context, page, perPage int) (*entity.ProductListResponse, error)

	// Update updates a product
	Update(ctx context.Context, productID uuid.UUID, req *entity.UpdateProductRequest) (*entity.Product, error)

	// Delete soft deletes a product
	Delete(ctx context.Context, productID uuid.UUID) error

	// UpdateStock updates product stock
	UpdateStock(ctx context.Context, productID uuid.UUID, stock int) error
}
```

### Step 8: Implement Usecase

Create `internal/domains/product/usecase/product_usecase.go` - See user usecase for reference implementation pattern.

### Step 9: Create HTTP Handler

Create `internal/domains/product/delivery/http_handler.go` - See user handler for reference implementation pattern.

### Step 10: Register Routes

Create `internal/domains/product/delivery/routes.go`:

```go
package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-username/your-project/internal/config"
	"github.com/your-username/your-project/internal/middleware"
)

// RegisterRoutes registers product routes
func RegisterRoutes(router fiber.Router, handler *ProductHandler, cfg *config.Config) {
	// Protected routes
	products := router.Group("/products", middleware.AuthMiddleware(cfg))
	products.Post("", handler.CreateProduct)
	products.Get("/:id", handler.GetProductByID)
	products.Get("", handler.ListProducts)
	products.Put("/:id", handler.UpdateProduct)
	products.Delete("/:id", handler.DeleteProduct)
	products.Patch("/:id/stock", handler.UpdateStock)
}
```

### Step 11: Register Domain in Main Server

Edit `cmd/server/main.go`:

```go
// Add imports
import (
	productDelivery "github.com/your-username/your-project/internal/domains/product/delivery"
	productRepo "github.com/your-username/your-project/internal/domains/product/repository"
	productUsecase "github.com/your-username/your-project/internal/domains/product/usecase"
)

// In main() function, after user domain initialization:

// Initialize product repository
productRepository := productRepo.NewPostgresProductRepository(db.Pool)

// Initialize product usecase
productUsecaseInstance := productUsecase.NewProductUsecase(
	productRepository,
	validatorInstance,
	cfg,
)

// Initialize product handler
productHandler := productDelivery.NewProductHandler(productUsecaseInstance)

// Register product routes
productDelivery.RegisterRoutes(v1, productHandler, cfg)
```

### Step 12: Write Tests

Create tests following the pattern in `internal/domains/user/usecase/user_usecase_test.go`.

### Step 13: Generate Mocks and Documentation

```bash
# Generate mocks
make mocks

# Generate API documentation
make swag

# Run tests
make test
```

## Testing Your New Domain

### 1. Start the Server

```bash
make dev
```

### 2. Test Endpoints

**Create Product:**
```bash
curl -X POST http://localhost:3000/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "price": 29.99,
    "stock": 100,
    "category_id": "UUID_HERE"
  }'
```

**Get Product:**
```bash
curl http://localhost:3000/api/v1/products/PRODUCT_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**List Products:**
```bash
curl "http://localhost:3000/api/v1/products?page=1&per_page=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Best Practices

### 1. Follow Existing Patterns

- Use the user domain as a reference
- Maintain consistent naming conventions
- Follow the same layer structure

### 2. Validation

- Add validation tags to all request structs
- Use custom validators when needed
- Return clear validation error messages

### 3. Error Handling

- Wrap errors with context
- Use domain-specific error types
- Return appropriate HTTP status codes

### 4. Testing

- Write unit tests for usecase layer
- Mock external dependencies
- Aim for 90%+ coverage
- Test both success and failure cases

### 5. Documentation

- Add Swagger comments to all handlers
- Document complex business logic
- Keep API docs up to date

### 6. Database

- Always create up and down migrations
- Add appropriate indexes
- Use transactions for multi-step operations
- Implement soft deletes where appropriate

### 7. Security

- Require authentication for sensitive endpoints
- Validate all input data
- Implement proper authorization checks
- Sanitize user input

## Common Patterns

### Soft Delete

```go
// In repository
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE table_name SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
```

### Pagination

```go
// In usecase
func (u *usecase) List(ctx context.Context, page, perPage int) (*entity.ListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	items, err := u.repo.List(ctx, perPage, offset)
	// ... handle error and count
}
```

### Filtering

```go
// Add filter parameters to repository methods
func (r *PostgresRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*entity.Product, error) {
	query := `
		SELECT * FROM products
		WHERE category_id = $1 AND is_active = TRUE
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	// ... implementation
}
```

## Checklist

Before considering your domain complete:

- [ ] Entity models defined with validation tags
- [ ] Database migration created and applied
- [ ] SQL queries written and generated
- [ ] Repository interface and implementation
- [ ] Usecase interface and implementation
- [ ] HTTP handler with Swagger documentation
- [ ] Routes registered
- [ ] Domain integrated in main server
- [ ] Unit tests written with mocks
- [ ] All tests passing
- [ ] API documentation generated
- [ ] Manual endpoint testing completed
- [ ] Error cases handled properly
- [ ] Input validation implemented
- [ ] Authorization checks in place

## Additional Resources

- **User Domain**: Complete reference implementation in `internal/domains/user/`
- **Testing**: Examples in `internal/domains/user/usecase/user_usecase_test.go`
- **Validation**: Custom validators in `pkg/validator/validator.go`
- **Response Format**: Standard responses in `pkg/response/response.go`
- **Middleware**: Authentication and other middleware in `internal/middleware/`

## Support

If you encounter issues:
1. Check the user domain implementation
2. Review error messages carefully
3. Check database connectivity
4. Verify migration status
5. Run tests to identify issues

For help:
- Create an issue on GitHub
- Review documentation in `docs/`
- Check existing domain implementations
