# Adding New Domains

This guide shows how to add a new domain to the Go Fiber Skeleton template following the Clean Architecture pattern.

## 🎯 Example: Adding a Product Domain

Let's add a `product` domain with CRUD operations.

## 📁 Step 1: Create Domain Structure

```bash
mkdir -p internal/domains/product/{entity,repository,usecase,handler}
```

## 📝 Step 2: Define Entity

Create `internal/domains/product/entity/product.go`:

```go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type Product struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    Price       float64   `json:"price" db:"price"`
    Stock       int       `json:"stock" db:"stock"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
    Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
    Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
    Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
    Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

func NewProduct(name, description string, price float64, stock int) *Product {
    now := time.Now()
    return &Product{
        ID:          uuid.New(),
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
}

func (p *Product) Update(name, description string, price float64, stock int) {
    p.Name = name
    p.Description = description
    p.Price = price
    p.Stock = stock
    p.UpdatedAt = time.Now()
}
```

## 🗄️ Step 3: Add Database Queries

Add to `db/query.sql`:

```sql
-- name: CreateProduct :one
INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, description, price, stock, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, stock = $5, updated_at = $6
WHERE id = $1
RETURNING id, name, description, price, stock, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
```

Add migration `db/migrations/002_create_products_table.up.sql`:

```sql
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_created_at ON products(created_at);
```

Add down migration `db/migrations/002_create_products_table.down.sql`:

```sql
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_name;
DROP TABLE IF EXISTS products;
```

Update `db/migrations` to include the products table.

Run: `make sqlc`

## 🔌 Step 4: Create Repository

Create `internal/domains/product/repository/product_repository.go`:

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "github.com/zercle/gofiber-skeleton/internal/domains/product/entity"
    "github.com/zercle/gofiber-skeleton/internal/infrastructure/database/sqlc"
)

type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type productRepository struct {
    queries *sqlc.Queries
}

func NewProductRepository(queries *sqlc.Queries) ProductRepository {
    return &productRepository{
        queries: queries,
    }
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
    params := sqlc.CreateProductParams{
        ID:          product.ID,
        Name:        product.Name,
        Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
        Price:       product.Price,
        Stock:       int32(product.Stock),
        CreatedAt:   product.CreatedAt,
        UpdatedAt:   product.UpdatedAt,
    }
    _, err := r.queries.CreateProduct(ctx, params)
    return err
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    dbProduct, err := r.queries.GetProductByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return r.dbProductToEntity(&dbProduct), nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
    dbProducts, err := r.queries.ListProducts(ctx, int32(limit), int32(offset))
    if err != nil {
        return nil, err
    }

    products := make([]*entity.Product, len(dbProducts))
    for i, dbProduct := range dbProducts {
        products[i] = r.dbProductToEntity(&dbProduct)
    }
    return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
    params := sqlc.UpdateProductParams{
        ID:          product.ID,
        Name:        product.Name,
        Description: sql.NullString{String: product.Description, Valid: product.Description != ""},
        Price:       product.Price,
        Stock:       int32(product.Stock),
        UpdatedAt:   product.UpdatedAt,
    }
    _, err := r.queries.UpdateProduct(ctx, params)
    return err
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.queries.DeleteProduct(ctx, id)
}

func (r *productRepository) dbProductToEntity(dbProduct *sqlc.Product) *entity.Product {
    description := ""
    if dbProduct.Description.Valid {
        description = dbProduct.Description.String
    }

    return &entity.Product{
        ID:          dbProduct.ID,
        Name:        dbProduct.Name,
        Description: description,
        Price:       dbProduct.Price,
        Stock:       int(dbProduct.Stock),
        CreatedAt:   dbProduct.CreatedAt,
        UpdatedAt:   dbProduct.UpdatedAt,
    }
}
```

## 🎯 Step 5: Create Usecase

Create `internal/domains/product/usecase/product_usecase.go`:

```go
package usecase

import (
    "context"
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/zercle/gofiber-skeleton/internal/domains/product/entity"
    "github.com/zercle/gofiber-skeleton/internal/domains/product/repository"
)

var (
    ErrProductNotFound = errors.New("product not found")
)

type ProductUsecase interface {
    Create(ctx context.Context, req *entity.CreateProductRequest) (*entity.Product, error)
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
    Update(ctx context.Context, id uuid.UUID, req *entity.UpdateProductRequest) (*entity.Product, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type productUsecase struct {
    productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
    return &productUsecase{
        productRepo: productRepo,
    }
}

func (u *productUsecase) Create(ctx context.Context, req *entity.CreateProductRequest) (*entity.Product, error) {
    product := entity.NewProduct(req.Name, req.Description, req.Price, req.Stock)

    if err := u.productRepo.Create(ctx, product); err != nil {
        return nil, fmt.Errorf("failed to create product: %w", err)
    }

    return product, nil
}

func (u *productUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    product, err := u.productRepo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrProductNotFound
    }
    return product, nil
}

func (u *productUsecase) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
    products, err := u.productRepo.List(ctx, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list products: %w", err)
    }
    return products, nil
}

func (u *productUsecase) Update(ctx context.Context, id uuid.UUID, req *entity.UpdateProductRequest) (*entity.Product, error) {
    product, err := u.productRepo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrProductNotFound
    }

    if req.Name != nil {
        product.Name = *req.Name
    }
    if req.Description != nil {
        product.Description = *req.Description
    }
    if req.Price != nil {
        product.Price = *req.Price
    }
    if req.Stock != nil {
        product.Stock = *req.Stock
    }

    product.UpdatedAt = time.Now()

    if err := u.productRepo.Update(ctx, product); err != nil {
        return nil, fmt.Errorf("failed to update product: %w", err)
    }

    return product, nil
}

func (u *productUsecase) Delete(ctx context.Context, id uuid.UUID) error {
    if err := u.productRepo.Delete(ctx, id); err != nil {
        return fmt.Errorf("failed to delete product: %w", err)
    }
    return nil
}
```

## 🌐 Step 6: Create Handler

Create `internal/domains/product/handler/product_handler.go`:

```go
package handler

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/zercle/gofiber-skeleton/internal/domains/product/entity"
    "github.com/zercle/gofiber-skeleton/internal/domains/product/usecase"
    "github.com/zercle/gofiber-skeleton/internal/shared/middleware"
    "github.com/zercle/gofiber-skeleton/internal/shared/response"
    "github.com/zercle/gofiber-skeleton/internal/shared/validator"
)

type ProductHandler struct {
    productUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
    return &ProductHandler{
        productUsecase: productUsecase,
    }
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param request body entity.CreateProductRequest true "Product data"
// @Security BearerAuth
// @Success 201 {object} response.Response{data=entity.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
    var req entity.CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err.Error())
    }

    if err := validator.Validate(&req); err != nil {
        return response.ValidationError(c, err.Error())
    }

    product, err := h.productUsecase.Create(c.Context(), &req)
    if err != nil {
        return response.InternalServerError(c, "Failed to create product", err.Error())
    }

    return response.Created(c, "Product created successfully", product)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
    productID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid product ID", err.Error())
    }

    product, err := h.productUsecase.GetByID(c.Context(), productID)
    if err != nil {
        if err == usecase.ErrProductNotFound {
            return response.NotFound(c, "Product not found", err.Error())
        }
        return response.InternalServerError(c, "Failed to get product", err.Error())
    }

    return response.OK(c, "Product retrieved successfully", product)
}

// ListProducts godoc
// @Summary List products
// @Description Get a paginated list of products
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]entity.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
    page, err := strconv.Atoi(c.Query("page", "1"))
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(c.Query("limit", "10"))
    if err != nil || limit < 1 || limit > 100 {
        limit = 10
    }

    offset := (page - 1) * limit

    products, err := h.productUsecase.List(c.Context(), limit, offset)
    if err != nil {
        return response.InternalServerError(c, "Failed to list products", err.Error())
    }

    return response.OK(c, "Products retrieved successfully", products)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body entity.UpdateProductRequest true "Product update data"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.Product}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
    productID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid product ID", err.Error())
    }

    var req entity.UpdateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err.Error())
    }

    if err := validator.Validate(&req); err != nil {
        return response.ValidationError(c, err.Error())
    }

    product, err := h.productUsecase.Update(c.Context(), productID, &req)
    if err != nil {
        if err == usecase.ErrProductNotFound {
            return response.NotFound(c, "Product not found", err.Error())
        }
        return response.InternalServerError(c, "Failed to update product", err.Error())
    }

    return response.OK(c, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
    productID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid product ID", err.Error())
    }

    if err := h.productUsecase.Delete(c.Context(), productID); err != nil {
        if err == usecase.ErrProductNotFound {
            return response.NotFound(c, "Product not found", err.Error())
        }
        return response.InternalServerError(c, "Failed to delete product", err.Error())
    }

    return response.NoContent(c, "Product deleted successfully")
}
```

## 🔧 Step 7: Update Main Application

Update `cmd/server/main.go`:

1. Add DI container providers:
```go
do.Provide(di, func(i *do.Injector) repository.ProductRepository {
    queries := do.MustInvoke[*sqlc.Queries](i)
    return repository.NewProductRepository(queries)
})

do.Provide(di, func(i *do.Injector) usecase.ProductUsecase {
    productRepo := do.MustInvoke[repository.ProductRepository](i)
    return usecase.NewProductUsecase(productRepo)
})

do.Provide(di, func(i *do.Injector) *handler.ProductHandler {
    productUsecase := do.MustInject[usecase.ProductUsecase](i)
    return handler.NewProductHandler(productUsecase)
})
```

2. Add routes:
```go
productHandler := do.MustInvoke[*handler.ProductHandler](di)

products := api.Group("/products", authMiddleware.RequireAuth())
products.Post("/", productHandler.CreateProduct)
products.Get("/", productHandler.ListProducts)
products.Get("/:id", productHandler.GetProduct)
products.Put("/:id", productHandler.UpdateProduct)
products.Delete("/:id", productHandler.DeleteProduct)
```

## 🧪 Step 8: Add Tests

Create test files for each layer following the user domain pattern.

## 🚀 Step 9: Run and Test

```bash
# Run migrations
make migrate-up

# Start server
make dev

# Test your new endpoints
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","price":29.99,"stock":100}'
```

## ✅ Checklist

- [ ] Domain structure created
- [ ] Entity with validation
- [ ] Database queries added
- [ ] Migration created and applied
- [ ] Repository interface implemented
- [ ] Usecase with business logic
- [ ] Handler with Swagger docs
- [ ] DI container updated
- [ ] Routes registered
- [ ] Tests written
- [ ] Documentation updated

That's it! Your new domain follows the same Clean Architecture patterns as the user domain.