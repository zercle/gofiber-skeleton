# Go Fiber Backend Mono-Repo Template

## 1. Introduction

This document provides a standardized template for initializing new backend mono-repo projects using **Go Fiber** framework. The boilerplate follows **Domain-Driven Clean Architecture** with each domain having its own complete Clean Architecture implementation within its domain package, promoting **SOLID principles** and domain isolation, mirroring the successful Hono.js mono-repo template architecture.

**Purpose:** Streamline mono-repo setup with domain-specific Clean Architecture using high-performance Go Fiber applications, ensuring scalability, maintainability, and clear domain boundaries.

---

## 2. Core Principles & Architecture

The boilerplate adopts a **mono-repo with domain-specific Clean Architecture** approach where each domain is a complete, self-contained module with its own layered architecture, identical to the Hono.js mono-repo implementation.

### Mono-Repo Domain Architecture

Each domain package contains its complete Clean Architecture implementation:

- **Domain Package** (`internal/domains/{domain}/`): Complete domain module
  - **Entities** (`entities/`): Core domain entities and value objects
  - **Use Cases** (`usecases/`): Domain-specific business logic
  - **Repositories** (`repositories/`): Domain data access interfaces and implementations  
  - **Handlers** (`handlers/`): Domain HTTP handlers
  - **Routes** (`routes/`): Domain-specific route definitions
  - **Models** (`models/`): DTOs and request/response models
  - **Tests** (`tests/`): Domain-specific tests

### Shared Infrastructure

- **Shared Infrastructure** (`internal/shared/`): Cross-domain shared components
  - **Database** (`database/`): GORM client and connection management
  - **Middleware** (`middleware/`): Authentication, validation, CORS
  - **Config** (`config/`): Application configuration
  - **Utils** (`utils/`): Shared utilities and helpers
  - **Types** (`types/`): Shared Go types and interfaces

### Key Tools & Libraries

| Feature               | Tool / Library                    |
| --------------------- | --------------------------------- |
| HTTP Framework        | Go Fiber v2                       |
| Dependency Injection  | fx (Uber's DI framework)          |
| Database ORM          | GORM                              |
| Configuration         | Viper                             |
| Migrations            | [golang-migrate/migrate](https://github.com/golang-migrate/migrate)            |
| SQL Code Generation   | [sqlc](https://sqlc.dev/)                   |
| JSON Response Format  | [omniti-labs/jsend](https://github.com/omniti-labs/jsend)                 |
| Mock Generation       | [go.uber.org/mock/mockgen](https://github.com/uber-go/mockgen)          |
| SQL Mocking           | [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)    |
| JWT Authentication    | golang-jwt                        |
| UUID Generation       | **UUIDv7** (index-friendly)       |
| Validation            | go-playground/validator           |
| Testing               | Go testing + Testify             |
| Hot Reload            | Air                              |
| Linting & Formatting  | golangci-lint, gofmt             |
| OpenAPI Documentation | swaggo/swag                     |

---

## 3. Getting Started: Project Initialization

### Prerequisites

- Go 1.21+ installed
- Docker & Docker Compose
- Air for hot reload (optional)

### Initial Setup

1. **Initialize Project**

```bash
go mod init hono-skeleton-go
```

2. **Configure Environment**

Create `.env` and `.env.example` with database URL and other configs.

3. **Setup Database**

- Define models in `domain/entities/`.
- Run migrations:

```bash
go run cmd/migrate/main.go
```

4. **Create Application Structure**

Follow layered architecture folders under project root.

5. **Run Application**

```bash
go run main.go
# Or with hot reload
air
```

---

## 4. Mono-Repo Project Structure Overview

```
.
├── cmd/
│   ├── server/
│   │   └── main.go          # Application entry point
│   └── migrate/
│       └── main.go          # Database migration runner
├── internal/
│   ├── app/
│   │   └── app.go           # Fiber app configuration
│   ├── domains/             # Domain-specific modules
│   │   ├── auth/            # Authentication domain
│   │   │   ├── entities/
│   │   │   │   └── user.go
│   │   │   ├── usecases/
│   │   │   │   ├── auth_usecase.go
│   │   │   │   └── interfaces/
│   │   │   ├── repositories/
│   │   │   │   ├── user_repository.go
│   │   │   │   └── interfaces/
│   │   │   ├── handlers/
│   │   │   │   └── auth_handler.go
│   │   │   ├── routes/
│   │   │   │   └── auth_routes.go
│   │   │   ├── models/
│   │   │   │   └── auth_models.go
│   │   │   └── tests/
│   │   │       ├── auth_usecase_test.go
│   │   │       └── auth_handler_test.go
│   │   ├── posts/           # Posts domain
│   │   │   ├── entities/
│   │   │   │   └── post.go
│   │   │   ├── usecases/
│   │   │   ├── repositories/
│   │   │   ├── handlers/
│   │   │   ├── routes/
│   │   │   ├── models/
│   │   │   └── tests/
│   │   └── greeting/        # Greeting domain (example)
│   │       ├── entities/
│   │       ├── usecases/
│   │       ├── repositories/
│   │       ├── handlers/
│   │       ├── routes/
│   │       ├── models/
│   │       └── tests/
│   └── infrastructure/              # Shared infrastructure
│       ├── database/
│       │   ├── connection.go
│       │   └── base_repository.go
│       ├── middleware/
│       │   ├── auth_middleware.go
│       │   ├── cors_middleware.go
│       │   └── validation_middleware.go
│       ├── config/
│       │   └── app_config.go
│       ├── utils/
│       │   ├── uuidv7.go    # UUIDv7 generator
│       │   ├── response.go
│       │   └── validation.go
│       ├── types/
│       │   └── common.go
│       └── container/
│           └── di_container.go
├── pkg/
│   └── utils/              # Public utilities
├── migrations/             # Database migration files
├── docs/                   # Swagger documentation
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .air.toml               # Air configuration for hot reload
└── README.md
```

---

## 5. Adding a New Domain (Mono-Repo Approach)

1. **Create Domain Directory Structure**

```bash
mkdir -p internal/domains/{domain-name}/{entities,usecases,repositories,handlers,routes,models,tests}
mkdir -p internal/domains/{domain-name}/usecases/interfaces
mkdir -p internal/domains/{domain-name}/repositories/interfaces
```

2. **Define Domain Entity with UUIDv7**

Create entity in `internal/domains/{domain}/entities/` with UUIDv7 ID:

```go
package entities

import (
    "time"
    "github.com/your-org/your-repo/internal/shared/utils"
)

type DomainEntity struct {
    ID        string    `json:"id" gorm:"type:uuid;primaryKey" example:"01234567-89ab-cdef-0123-456789abcdef"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (e *DomainEntity) BeforeCreate(tx *gorm.DB) error {
    if e.ID == "" {
        e.ID = utils.GenerateUUIDv7()
    }
    return nil
}
```

3. **Create Domain Models**

Define request/response DTOs in `internal/domains/{domain}/models/`.

4. **Implement Repository Interface & Implementation**

- Interface in `internal/domains/{domain}/repositories/interfaces/`
- Implementation in `internal/domains/{domain}/repositories/`

5. **Implement Use Case Interface & Implementation**

- Interface in `internal/domains/{domain}/usecases/interfaces/`
- Implementation in `internal/domains/{domain}/usecases/`

6. **Implement Handlers**

Create HTTP handlers in `internal/domains/{domain}/handlers/`.

7. **Define Routes**

Create routes in `internal/domains/{domain}/routes/` and register in app.

8. **Add Database Model**

Define GORM model with UUIDv7 default:

```go
type DomainModel struct {
    ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    // ... other fields
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

9. **Register Dependencies**

Update `internal/shared/container/di_container.go` with new domain services using fx.

10. **Write Domain Tests**

Create comprehensive tests in `internal/domains/{domain}/tests/`.

11. **Generate API Documentation**

Add Swagger annotations and regenerate docs with swaggo.

---

## 6. Development Commands

- **Run development server**

```bash
go run cmd/server/main.go
# Or with hot reload
air
```

- **Run database migrations**

```bash
go run cmd/migrate/main.go
```

- **Run tests**

```bash
go test ./...
go test -v ./internal/tests/...
```

- **Lint and format**

```bash
golangci-lint run
gofmt -s -w .
```

- **Generate API documentation**

```bash
swag init -g cmd/server/main.go -o docs
```

- **Build for production**

```bash
go build -o bin/server cmd/server/main.go
```

## 7. Example Domain Implementation

### User Domain (Authentication)

**Entity (`internal/domain/entities/user.go`):**
```go
type User struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Repository Interface (`internal/domain/interfaces/user_repository.go`):**
```go
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    GetByEmail(ctx context.Context, email string) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**Use Case (`internal/usecases/auth_usecase.go`):**
```go
type AuthUseCase struct {
    userRepo interfaces.UserRepository
    jwtService *JWTService
}

func (uc *AuthUseCase) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
    // Implementation
}

func (uc *AuthUseCase) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
    // Implementation
}
```

**Handler (`internal/handlers/auth_handler.go`):**
```go
type AuthHandler struct {
    authUC usecases.AuthUseCase
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // Implementation
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // Implementation
}
```

### Post Domain (Example CRUD)

**Entity (`internal/domain/entities/post.go`):**
```go
type Post struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
    Title     string    `json:"title" gorm:"not null"`
    Content   string    `json:"content"`
    AuthorID  uuid.UUID `json:"author_id" gorm:"not null"`
    Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

---

## 8. Key Features Implementation

### JWT Authentication Middleware
```go
func JWTMiddleware(secret string) fiber.Handler {
    return jwtware.New(jwtware.Config{
        SigningKey: []byte(secret),
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        },
    })
}
```

### Database Configuration with GORM
```go
func NewDatabase(config *Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto-migrate models
    err = db.AutoMigrate(
        &entities.User{},
        &entities.Post{},
    )
    
    return db, err
}
```

### Dependency Injection with fx
```go
var Module = fx.Options(
    fx.Provide(
        config.NewConfig,
        database.NewDatabase,
        repositories.NewUserRepository,
        repositories.NewPostRepository,
        usecases.NewAuthUseCase,
        usecases.NewPostUseCase,
        handlers.NewAuthHandler,
        handlers.NewPostHandler,
        app.NewApp,
    ),
)
```

### Validation with go-playground/validator
```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Name     string `json:"name" validate:"required"`
}
```

---

## 9. Testing Strategy

### Repository Tests
```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := repositories.NewUserRepository(db)
    
    user := &entities.User{
        Email: "test@example.com",
        Name: "Test User",
    }
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, user.ID)
}
```

### Handler Tests
```go
func TestAuthHandler_Register(t *testing.T) {
    app := setupTestApp(t)
    
    req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{
        "email": "test@example.com",
        "password": "password123",
        "name": "Test User"
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

---

## 10. Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/server .
COPY --from=builder /app/.env.example .env

CMD ["./server"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://user:pass@postgres:5432/dbname?sslmode=disable
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: hono_skeleton_go
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

---

## 11. Performance Optimizations

### Database Connection Pooling
```go
func NewDatabase(config *Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}
```

### Response Caching Middleware
```go
func CacheMiddleware(duration time.Duration) fiber.Handler {
    return cache.New(cache.Config{
        Expiration: duration,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.OriginalURL()
        },
    })
}
```

---

## 12. API Documentation with Swagger

### Swagger Annotations Example
```go
// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // Implementation
}
```

---

## 13. Environment Configuration

### .env.example
```env
# Application
PORT=3000
ENV=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/hono_skeleton_go?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h

# Redis
REDIS_URL=redis://localhost:6379

# CORS
CORS_ORIGINS=*
```

---

## 14. Documentation References

- **Go**: https://golang.org/doc/
- **Fiber**: https://gofiber.io/
- **GORM**: https://gorm.io/docs/
- **fx**: https://uber-go.github.io/fx/
- **Viper**: https://github.com/spf13/viper

---