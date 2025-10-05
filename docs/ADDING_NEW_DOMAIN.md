# Adding a New Domain to Go Fiber Skeleton

This guide walks you through adding a new business domain to the Go Fiber Skeleton template. Following this guide ensures your new domain adheres to the Clean Architecture pattern and integrates seamlessly with the existing codebase.

## Overview

Each domain in this template follows a 4-layer architecture:
1. **Entity Layer** - Domain models (pure Go structs)
2. **Repository Layer** - Data access (interfaces + implementations)
3. **Usecase Layer** - Business logic (interfaces + implementations)
4. **Handler Layer** - HTTP handlers and route registration

## Prerequisites

- Basic understanding of Clean Architecture
- Familiarity with Go and Fiber framework
- Project dependencies installed (`make install-tools`)

## Step-by-Step Guide

Let's create a new domain called `product` as an example.

### Step 1: Create Directory Structure

```bash
mkdir -p internal/product/entity
mkdir -p internal/product/repository/mocks
mkdir -p internal/product/usecase/mocks
mkdir -p internal/product/handler
mkdir -p internal/product/tests
```

### Step 2: Define the Entity

Create `internal/product/entity/product.go`:

```go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

### Step 3: Create Database Migration

Create migration file:
```bash
make migrate-create name=create_products
```

Edit the generated file `db/migrations/000X_create_products.up.sql`:

```sql
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_user_id ON products(user_id);
CREATE INDEX idx_products_created_at ON products(created_at DESC);
```

Create the down migration `db/migrations/000X_create_products.down.sql`:

```sql
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_user_id;
DROP TABLE IF EXISTS products;
```

### Step 4: Define SQL Queries

Create `db/queries/products.sql`:

```sql
-- name: CreateProduct :one
INSERT INTO products (name, description, price, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, price, user_id, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, price, user_id, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProductsByUser :many
SELECT id, name, description, price, user_id, created_at, updated_at
FROM products
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, price, user_id, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
```

### Step 5: Generate Type-Safe Code

```bash
make sqlc
```

This generates type-safe Go code in `internal/db/products.sql.go`.

### Step 6: Implement Repository

Create `internal/product/repository/postgres.go`:

```go
package repository

//go:generate mockgen -source=postgres.go -destination=mocks/repository.go -package=mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/product/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, name, description string, price float64, userID uuid.UUID) (*entity.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Product, error)
	Update(ctx context.Context, id uuid.UUID, name, description string, price float64) (*entity.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type productRepository struct {
	queries *db.Queries
}

func NewProductRepository(queries *db.Queries) ProductRepository {
	return &productRepository{queries: queries}
}

func (r *productRepository) Create(ctx context.Context, name, description string, price float64, userID uuid.UUID) (*entity.Product, error) {
	product, err := r.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// Implement other methods similarly...
```

### Step 7: Implement Usecase

Create `internal/product/usecase/product.go`:

```go
package usecase

//go:generate mockgen -source=product.go -destination=mocks/usecase.go -package=mocks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/product/entity"
	"github.com/zercle/gofiber-skeleton/internal/product/repository"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, name, description string, price float64, userID uuid.UUID) (*entity.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	ListProductsByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, userID uuid.UUID, name, description string, price float64) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) CreateProduct(ctx context.Context, name, description string, price float64, userID uuid.UUID) (*entity.Product, error) {
	return u.repo.Create(ctx, name, description, price, userID)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, id uuid.UUID, userID uuid.UUID, name, description string, price float64) (*entity.Product, error) {
	// Verify ownership
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing.UserID != userID {
		return nil, fmt.Errorf("unauthorized: user does not own this product")
	}

	return u.repo.Update(ctx, id, name, description, price)
}

// Implement other methods...
```

### Step 8: Generate Mocks

```bash
make generate-mocks
```

This creates mock implementations for testing in the `mocks/` directories.

### Step 9: Implement HTTP Handlers

Create `internal/product/handler/product_handler.go`:

```go
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/product/usecase"
	"github.com/zercle/gofiber-skeleton/internal/response"
	"github.com/zercle/gofiber-skeleton/pkg/validator"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(usecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: usecase}
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Product data"
// @Security BearerAuth
// @Success 201 {object} response.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, fiber.Map{"error": "Invalid request body"})
	}

	if err := validator.ValidateRequest(c, &req); err != nil {
		return err
	}

	userID := c.Locals("user_id").(uuid.UUID)

	product, err := h.usecase.CreateProduct(c.Context(), req.Name, req.Description, req.Price, userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create product", err)
	}

	return response.Success(c, fiber.StatusCreated, product)
}

// Implement other handlers...

func RegisterProductRoutes(app fiber.Router, handler *ProductHandler, authMiddleware fiber.Handler) {
	products := app.Group("/products")

	products.Post("/", authMiddleware, handler.CreateProduct)
	products.Get("/:id", handler.GetProductByID)
	products.Get("/user/:user_id", handler.ListProductsByUser)
	products.Put("/:id", authMiddleware, handler.UpdateProduct)
	products.Delete("/:id", authMiddleware, handler.DeleteProduct)
}
```

### Step 10: Write Tests

Create `internal/product/tests/product_usecase_test.go`:

```go
package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/product/entity"
	mockRepo "github.com/zercle/gofiber-skeleton/internal/product/repository/mocks"
	"github.com/zercle/gofiber-skeleton/internal/product/usecase"
)

func TestProductUsecase_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mockRepo.NewMockProductRepository(ctrl)
	productUsecase := usecase.NewProductUsecase(mockProductRepo)

	ctx := context.Background()
	name := "Test Product"
	description := "Test Description"
	price := 99.99
	userID := uuid.New()

	expectedProduct := &entity.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
	}

	mockProductRepo.EXPECT().
		Create(ctx, name, description, price, userID).
		Return(expectedProduct, nil).
		Times(1)

	result, err := productUsecase.CreateProduct(ctx, name, description, price, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, name, result.Name)
}
```

### Step 11: Register Routes in Router

Edit `internal/server/router.go` to register your new domain:

```go
// Add to imports
import (
	productHandler "github.com/zercle/gofiber-skeleton/internal/product/handler"
	productRepo "github.com/zercle/gofiber-skeleton/internal/product/repository"
	productUsecase "github.com/zercle/gofiber-skeleton/internal/product/usecase"
)

// In SetupRouter function, after existing domain setup:
func SetupRouter(app *fiber.App, db *sql.DB, cfg *config.Config) {
	// ... existing code ...

	// Product domain
	productRepository := productRepo.NewProductRepository(queries)
	productUsecaseInstance := productUsecase.NewProductUsecase(productRepository)
	productHandlerInstance := productHandler.NewProductHandler(productUsecaseInstance)
	productHandler.RegisterProductRoutes(api, productHandlerInstance, authMiddleware.Authenticate)
}
```

### Step 12: Run Migrations and Test

```bash
# Run database migrations
make migrate-up

# Run tests
make test

# Generate updated documentation
make generate-docs

# Run the application
make dev
```

## Testing Your Domain

1. **Unit Tests**: Test usecases with mocked repositories
2. **Repository Tests**: Test repositories with go-sqlmock
3. **Integration Tests**: Test handlers with test database

Example test execution:
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run with race detector
make test-race
```

## Best Practices

1. **Always use interfaces** for repositories and usecases
2. **Add validation** to all request DTOs
3. **Implement ownership checks** in usecases where applicable
4. **Use sqlc** for all database queries (no raw SQL in repositories)
5. **Write tests** before considering the domain complete
6. **Add Swagger annotations** to all handlers
7. **Use JSend format** for all API responses
8. **Follow naming conventions** (e.g., `CreateProduct`, not `AddProduct`)

## Checklist

- [ ] Entity defined in `entity/` directory
- [ ] Database migration created (up and down)
- [ ] SQL queries defined in `db/queries/`
- [ ] sqlc code generated (`make sqlc`)
- [ ] Repository interface and implementation created
- [ ] Usecase interface and implementation created
- [ ] HTTP handlers implemented
- [ ] Routes registered in router
- [ ] Mocks generated (`make generate-mocks`)
- [ ] Unit tests written
- [ ] Swagger documentation added
- [ ] Integration tested with `make dev`

## Common Pitfalls

- **Forgetting to add foreign key constraints** in migrations
- **Not checking ownership** before update/delete operations
- **Missing validation tags** on request structs
- **Forgetting to regenerate mocks** after interface changes
- **Not handling database errors** properly in repositories
- **Circular dependencies** between domains (avoid!)

## Example Domains

Refer to these existing domains for reference:
- `internal/user` - Authentication and user management
- `internal/post` - Posts with ownership validation

## Need Help?

- Check existing domains for patterns
- Review the architecture documentation in `ARCHITECTURE.md`
- Consult the Memory Bank in `.agents/rules/memory-bank/`

---

**Time to Complete**: Approximately 15-20 minutes for a simple CRUD domain.

**Next Steps**: After successfully adding a domain, consider adding additional features like pagination, search, or complex business logic.
