# Architecture Overview

## Design Pattern
**Domain-Driven Clean Architecture** within a **mono-repo structure**

### Core Principles
- **Domain Isolation**: Each business domain (user, post, thread) is self-contained with zero cross-domain dependencies
- **Dependency Inversion**: Inner layers (entities, usecases) never depend on outer layers (handlers, repositories)
- **Interface Segregation**: Repository and usecase interfaces enable mockgen-based testing
- **SOLID Compliance**: Clear separation of concerns across entity, repository, usecase, and handler layers

## Directory Structure

```
gofiber-skeleton/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point, DB initialization, graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go            # Viper-based config (env vars > .env file)
│   ├── db/                      # sqlc-generated code
│   │   ├── db.go                # DBTX interface, Queries struct
│   │   ├── models.go            # Generated DB models
│   │   ├── users.sql.go         # Generated user queries
│   │   └── posts.sql.go         # Generated post queries
│   ├── logger/
│   │   └── logger.go            # Zerolog structured logger
│   ├── middleware/
│   │   ├── logging.go           # Request logging with duration/status
│   │   ├── rate_limit.go        # API and auth rate limiting
│   │   └── request_id.go        # Request ID generation
│   ├── response/
│   │   └── jsend.go             # JSend response format (success/fail/error)
│   ├── server/
│   │   └── router.go            # Fiber app setup, route registration, health checks
│   ├── shared/
│   │   └── types/               # Shared types (currently empty)
│   ├── user/                    # USER DOMAIN
│   │   ├── entity/
│   │   │   └── user.go          # User business entity
│   │   ├── repository/
│   │   │   └── postgres.go      # UserRepository interface + Postgres implementation
│   │   ├── usecase/
│   │   │   └── auth.go          # AuthUsecase interface + implementation (Register, Login)
│   │   ├── handler/
│   │   │   └── auth_handler.go  # HTTP handlers + route registration
│   │   └── middleware/
│   │       └── auth_middleware.go # JWT authentication middleware
│   └── post/                    # POST DOMAIN
│       ├── entity/
│       │   └── post.go          # Post business entity
│       ├── repository/
│       │   ├── postgres.go      # PostRepository interface + implementation
│       │   └── mocks/
│       │       └── repository.go # mockgen-generated mocks
│       ├── usecase/
│       │   ├── post.go          # PostUsecase interface + implementation
│       │   └── mocks/
│       │       └── usecase.go   # mockgen-generated mocks
│       ├── handler/
│       │   └── post_handler.go  # HTTP handlers + route registration
│       └── tests/
│           └── post_test.go     # Placeholder tests
├── db/
│   ├── migrations/              # SQL migration files
│   │   ├── 001_create_users.up.sql
│   │   ├── 002_create_roles.up.sql
│   │   ├── 003_create_threads_posts_comments.up.sql
│   │   └── 004_create_sessions.up.sql
│   └── queries/                 # sqlc query definitions
│       ├── users.sql            # User CRUD queries
│       └── posts.sql            # Post CRUD queries
├── docs/                        # Swagger-generated API docs
├── .env.example                 # Environment variable template
├── compose.yml                  # PostgreSQL + Valkey/Redis services
├── Dockerfile                   # Multi-stage build (alpine-based)
├── Makefile                     # Development commands (build, test, migrate, sqlc, etc.)
├── sqlc.yaml                    # sqlc configuration
└── go.mod                       # Go module definition
```

## Layer Responsibilities

### 1. Entity Layer (`entity/`)
- **Role**: Business domain objects
- **Dependencies**: None (pure Go structs)
- **Example**: `user.go` defines User struct with ID, Username, Email, timestamps

### 2. Repository Layer (`repository/`)
- **Role**: Data access abstraction
- **Dependencies**: Entity layer, sqlc-generated code
- **Pattern**: Interface + Postgres implementation
- **Key Feature**: `//go:generate mockgen` annotation for test mocks
- **Example**: `UserRepository` interface with `Create()`, `GetByID()`, `GetByEmail()`

### 3. Usecase Layer (`usecase/`)
- **Role**: Business logic orchestration
- **Dependencies**: Entity, Repository interfaces
- **Pattern**: Interface + concrete implementation
- **Key Feature**: `//go:generate mockgen` annotation for test mocks
- **Example**: `AuthUsecase` handles registration (password hashing via bcrypt) and login (JWT generation)

### 4. Handler Layer (`handler/`)
- **Role**: HTTP request/response handling
- **Dependencies**: Usecase interfaces, response package
- **Responsibilities**: Request parsing, validation, usecase invocation, JSend response formatting
- **Example**: `RegisterAuthRoutes()` registers POST /auth/register and /auth/login

## Data Flow

**Request Flow** (HTTP Request → Response):
```
HTTP Request
  ↓
Handler (parse request, validate)
  ↓
Usecase (business logic)
  ↓
Repository (database operations via sqlc)
  ↓
PostgreSQL Database
  ↓
Repository (map DB models to entities)
  ↓
Usecase (process, enrich data)
  ↓
Handler (format JSend response)
  ↓
HTTP Response
```

## Key Technology Integration

### Database: sqlc + PostgreSQL
- **sqlc.yaml**: Schema from `db/migrations/`, queries from `db/queries/`, output to `internal/db/`
- **Workflow**: Write SQL → Run `make sqlc` → Type-safe Go code generated
- **Benefits**: Compile-time SQL validation, zero ORM overhead

### Authentication: JWT + bcrypt
- **Registration**: bcrypt hashes password (cost=10), stores in users table
- **Login**: Verifies hash, generates JWT with 72-hour expiry, returns token
- **Middleware**: `auth_middleware.go` validates Bearer token, extracts user_id to context

### Configuration: godotenv + os.LookupEnv
- **Priority**: Runtime env vars > .env file > hardcoded defaults
- **Structure**: `Config` struct with Port, DatabaseDSN, JWTSecret

### Logging: zerolog
- **Features**: Structured JSON logging, request ID tracking, duration metrics
- **Integration**: Middleware logs all HTTP requests with method, path, status, duration

### API Documentation: swaggo/swag
- **Workflow**: Annotate handlers with `// @Summary`, `// @Router` → Run `make generate-docs` → Swagger UI at `/swagger`
- **Location**: `cmd/server/main.go` has global API metadata

### Testing: mockgen + go-sqlmock
- **Unit Testing**: Interfaces annotated with `//go:generate mockgen` for mocking
- **Repository Testing**: Use `go-sqlmock` to simulate database responses
- **Current State**: Mock infrastructure present, full test suite pending implementation

## Domain Pattern (Using Post as Example)

```
internal/post/
├── entity/post.go              # Post{ID, ThreadID, UserID, Content, CreatedAt, UpdatedAt}
├── repository/
│   ├── postgres.go             # PostRepository interface: Create, GetByID, ListByUser, Update, Delete
│   └── mocks/repository.go     # Auto-generated mock
├── usecase/
│   ├── post.go                 # PostUsecase interface: CRUD + ownership validation
│   └── mocks/usecase.go        # Auto-generated mock
├── handler/
│   └── post_handler.go         # HTTP handlers, route registration
└── tests/
    └── post_test.go            # Placeholder test
```

## Deployment Architecture

**Development**: Docker Compose
- **Services**: app (Go Fiber), db (PostgreSQL 18), redis (Valkey 8)
- **Networking**: Internal Docker network, exposed ports 8080 (app), 6379 (redis)
- **Health Checks**: All services have health checks with retry logic

**Production** (Containerized):
- **Build**: Multi-stage Dockerfile (builder: golang:alpine → runtime: alpine)
- **Security**: Non-root user (appuser), static binary compilation
- **Configuration**: Environment variables injected at runtime

## Critical Paths

### Source Files (Most Frequently Modified)
- `cmd/server/main.go`: Application bootstrap
- `internal/server/router.go`: Route definitions, middleware registration
- `internal/config/config.go`: Configuration struct
- `db/migrations/*.sql`: Database schema
- `db/queries/*.sql`: SQL query definitions
- `internal/{domain}/`: Domain-specific implementations

### Generated Files (Never Edit Manually)
- `internal/db/*.go`: sqlc-generated code
- `docs/*`: Swagger-generated documentation
- `internal/*/mocks/*.go`: mockgen-generated test mocks

## Design Decisions

1. **No ORM**: Prefer raw SQL via sqlc for performance and explicitness
2. **No Uber fx (yet)**: Current implementation uses manual dependency injection in `router.go` (brief mentions fx but not implemented)
3. **JSend Response Format**: Standardized API responses (status: success/fail/error)
4. **UUID v7**: Time-sortable UUIDs for better database performance
5. **Graceful Shutdown**: 30-second timeout for in-flight requests
6. **Middleware Order**: Recover → RequestID → Logging → RateLimit