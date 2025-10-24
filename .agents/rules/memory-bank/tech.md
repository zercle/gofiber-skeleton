# Technology Stack: Go Fiber Microservice Template

## Language & Runtime
- **Go Version:** 1.25.3+ (stable)
- **Go Modules:** Yes, using go.mod/go.sum
- **Module Name:** github.com/zercle/template-go-fiber

## Core Framework

### Fiber (Web Framework)
- **Package:** github.com/gofiber/fiber/v2
- **Version:** v2.x (latest stable)
- **Purpose:** HTTP server and routing
- **Why:** Fastest Go web framework, Express-like API, minimal overhead
- **Key Features Used:**
  - Routing and middleware
  - Built-in middleware (CORS, compression)
  - Error handling
  - Request/response binding
- **Documentation:** https://docs.gofiber.io/

### Database Access

#### sqlc (Type-Safe SQL Generation)
- **Package:** github.com/sqlc-dev/sqlc/cmd/sqlc
- **Version:** v1.x (latest stable)
- **Purpose:** Generate type-safe Go code from SQL queries
- **Installation:** `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- **Configuration:** sqlc.yaml (in sql/ directory)
- **Key Features:**
  - Compile-time SQL verification
  - Automatic Go struct generation from queries
  - No ORM learning curve
  - Prevents SQL injection by default
- **Documentation:** https://sqlc.dev/
- **Supported Databases:**
  - PostgreSQL 12+
  - MySQL 5.7+
  - MariaDB 10.3+
  - SQLite

#### golang-migrate (Database Migrations)
- **Package:** github.com/golang-migrate/migrate/v4
- **Version:** v4.x (latest stable)
- **Purpose:** Version control for database schema
- **Installation:** `go install -tags 'mysql,postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- **Key Features:**
  - Forward and backward migrations
  - Multiple database support
  - CLI and programmatic API
  - Transactional migrations
- **Documentation:** https://github.com/golang-migrate/migrate
- **File Format:** SQL or YAML migrations
- **Location:** sql/migrations/

### Database Drivers

#### For MySQL/MariaDB
- **Package:** github.com/go-sql-driver/mysql
- **Version:** v1.x (latest stable)
- **Usage:** Import in database connection code
- **Default Choice:** MariaDB 11+ for small/medium systems

#### For PostgreSQL
- **Package:** github.com/lib/pq
- **Version:** v1.x (latest stable)
- **Usage:** Import in database connection code
- **Alternative Choice:** PostgreSQL 18+ for large systems with advanced features

#### Database Connection
- **sql Package:** Standard library `database/sql`
- **Connection Pooling:** Built-in via database/sql
- **Transactions:** database/sql.Tx interface

## Authentication & Security

### JWT (JSON Web Tokens)
- **Package:** github.com/golang-jwt/jwt/v5
- **Version:** v5.x (latest stable)
- **Purpose:** Token-based authentication
- **Key Functions:**
  - Token creation and signing
  - Token validation and verification
  - Claims extraction
- **Algorithm:** HS256 (HMAC-SHA256) for signing key, RS256 (RSA) for public key infrastructure
- **Documentation:** https://github.com/golang-jwt/jwt

### OIDC (OpenID Connect) [Optional]
- **Package:** github.com/zitadel/oidc/v2
- **Version:** v2.x
- **Purpose:** Enterprise SSO and OAuth2 integration
- **When to Use:** For companies with existing SSO infrastructure
- **Documentation:** https://pkg.go.dev/github.com/zitadel/oidc

### Password Hashing
- **Package:** golang.org/x/crypto/argon2
- **Version:** Standard library
- **Purpose:** Secure password hashing for user credentials
- **Algorithm:** Argon2id (memory-hard, resistant to GPU attacks)
- **Usage:** Hash on account creation, verify on login
- **Alternative:** bcrypt via golang.org/x/crypto/bcrypt (if Argon2 not available)

### Security Headers
- **Fiber Middleware:** github.com/gofiber/fiber/v2/middleware/helmet
- **Purpose:** Set secure HTTP headers
- **Headers:** X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, etc.

## Dependency Injection

### samber/do v2
- **Package:** github.com/samber/do/v2
- **Version:** v2.x (latest stable)
- **Purpose:** Lightweight DI container for managing dependencies
- **Key Features:**
  - Type-safe dependency resolution
  - Zero runtime cost
  - Interface-based injection
  - Lazy initialization
- **Usage Pattern:**
  ```go
  container := do.New()
  do.Provide(container, func(i *do.Injector) (UserRepository, error) {
      // Initialize and return
  })
  ```
- **Documentation:** https://github.com/samber/do

## Logging & Observability

### slog (Structured Logging)
- **Package:** log/slog
- **Version:** Built-in (Go 1.21+)
- **Purpose:** Structured, context-aware logging
- **Key Features:**
  - JSON output (machine-readable)
  - Context propagation
  - Multiple log levels
  - Structured fields
- **Usage:** Logs all operations with context
- **Format:** JSON lines (one log per line)
- **Integration:** Works with log aggregation systems (ELK, Splunk, etc.)

### Health Check Patterns
- **Implementation:** Custom HTTP endpoints (/health, /ready)
- **Purpose:** Kubernetes liveness and readiness probes
- **Endpoints:**
  - `/health/live` - Is service running?
  - `/health/ready` - Is service ready for traffic?

## Configuration Management

### Viper (Config Loading)
- **Package:** github.com/spf13/viper
- **Version:** v1.x (latest stable)
- **Purpose:** Load configuration from environment variables
- **Key Features:**
  - Multiple config sources (env vars, files, flags)
  - Type conversion and validation
  - Default values
  - Environment variable precedence
- **Usage Pattern:**
  ```go
  viper.BindEnv("database.host", "DB_HOST")
  dbHost := viper.GetString("database.host")
  ```
- **Configuration Sources:** Environment variables with APP_ prefix

## Testing & Mocking

### mockgen (Mock Generation)
- **Package:** github.com/uber-go/mock
- **Version:** v0.x (latest stable)
- **Installation:** `go install github.com/uber-go/mock/mockgen@latest`
- **Purpose:** Automatic mock generation from interfaces
- **Key Features:**
  - Code generation via //go:generate comments
  - Type-safe mocks
  - Match/expect patterns
- **Usage:**
  ```go
  //go:generate mockgen -destination=../mocks/mock_user_repository.go . UserRepository
  type UserRepository interface { ... }
  ```
- **Documentation:** https://github.com/uber-go/mock

### go-sqlmock (Database Mocking)
- **Package:** github.com/DATA-DOG/go-sqlmock
- **Version:** v1.x (latest stable)
- **Purpose:** Mock SQL database for repository testing
- **Key Features:**
  - Mock sql.DB interface
  - Query expectation matching
  - Result simulation
  - Error injection
- **Usage Pattern:**
  ```go
  db, mock, _ := sqlmock.New()
  mock.ExpectQuery("SELECT ...").WillReturnRows(...)
  ```
- **Documentation:** https://github.com/DATA-DOG/go-sqlmock

### Testing Framework
- **Package:** testing (standard library)
- **Test Pattern:** *_test.go files in same directory
- **Table-Driven Tests:** Recommended for multiple test cases
- **Test Coverage:** Run with `go test -cover ./...`

## Code Quality & Linting

### golangci-lint
- **Package:** github.com/golangci/golangci-lint
- **Version:** v2.x (latest stable)
- **Installation:** `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **Configuration:** `.golangci.yml` in project root
- **Linters Enabled:**
  - gofmt - Code formatting
  - govet - Vet analysis
  - errcheck - Unchecked errors
  - staticcheck - Static analysis
  - gosimple - Code simplification
  - unused - Unused code detection
  - ineffassign - Ineffective assignments
  - misspell - Spelling errors
  - typecheck - Type checking
- **Usage:** `golangci-lint run ./...`
- **Zero-Error Policy:** Template must pass with zero errors/warnings

### go fmt
- **Package:** Standard library (cmd/go)
- **Purpose:** Code formatting (Go standard style)
- **Usage:** Applied by golangci-lint, can run manually with `go fmt ./...`

## Utilities & Helper Libraries

### samber/lo (Helper Functions)
- **Package:** github.com/samber/lo
- **Version:** v1.x (latest stable)
- **Purpose:** Synchronous helper functions (similar to lodash)
- **Key Functions:**
  - Map, Filter, Reduce
  - Find, Contains, Unique
  - Chunk, Drop, Take
- **Usage:** Simplifies collection operations
- **Documentation:** https://github.com/samber/lo

### samber/ro (Stream/Event Patterns)
- **Package:** github.com/samber/ro
- **Version:** v1.x (latest stable)
- **Purpose:** Event-driven stream patterns
- **Key Features:**
  - Reactive/functional programming patterns
  - Observable streams
  - Back-pressure handling
- **When to Use:** For event-driven features or reactive workflows
- **Documentation:** https://github.com/samber/ro

### UUID v7 Generation
- **Package:** github.com/google/uuid
- **Version:** Latest stable
- **Purpose:** Generate database-friendly unique IDs
- **Why v7:** Sortable by timestamp, better for database indexing
- **Usage:** `id := uuid.New().String()` (UUIDv7 when available)

## Containerization

### Docker
- **Base Image:** golang:1.25-alpine (for multi-stage builds)
- **Runtime Image:** alpine:latest (minimal)
- **Purpose:** Container image for deployment
- **Dockerfile:** Multi-stage build for small image size
- **Image Size Target:** <50MB for final image
- **Key Patterns:**
  - Build stage compiles binary
  - Runtime stage contains only compiled binary
  - Non-root user for security
  - Health check endpoint configured

### Docker Compose
- **Purpose:** Local development environment
- **Services:** API container + Database container
- **Configuration:** docker-compose.yml in root
- **Usage:** `docker-compose up` for development

## Build & Deploy Tools

### Makefile
- **Purpose:** Development task automation
- **Key Commands:**
  - `make build` - Compile binary
  - `make run` - Run locally
  - `make test` - Run tests with coverage
  - `make lint` - Run linters
  - `make migrate` - Run migrations
  - `make docker-build` - Build Docker image
- **Location:** Makefile in project root

### Go Build Tags
- **Purpose:** Conditional compilation
- **Usage:** `// +build <tag>` in files for platform-specific code
- **Example:** Windows-specific or database-specific implementations

## Environment & Dependencies

### Required Tools (Development)
- Go 1.25.3+
- Docker & Docker Compose (for containerization)
- make (for Makefile)
- golangci-lint v2 (for linting)
- golang-migrate (for migrations)
- sqlc (for code generation)

### Go Dependency Management
- **go mod tidy** - Clean up unused dependencies
- **go mod vendor** - Download dependencies locally
- **go mod download** - Pre-download dependencies
- **go list -m all** - List all dependencies

### Development Environment Setup
```bash
# Clone repository
git clone https://github.com/zercle/template-go-fiber.git
cd template-go-fiber

# Install dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/uber-go/mock/mockgen@latest

# Run application
make run

# Run tests
make test

# Run linter
make lint
```

## API Documentation

### Swagger/OpenAPI
- **Package:** github.com/swaggo/swag
- **Version:** Latest stable
- **Purpose:** Automatic API documentation from annotations
- **Key Features:**
  - Annotations in handler code
  - Generates Swagger/OpenAPI spec
  - Interactive Swagger UI at /docs
- **Annotations:**
  ```go
  // @Summary Create user
  // @Description Create a new user
  // @Accept json
  // @Produce json
  // @Param input body CreateUserRequest true "Create User Request"
  // @Success 201 {object} User
  // @Failure 400 {object} ErrorResponse
  // @Router /users [post]
  ```
- **Generation:** `swag init -g cmd/api/main.go`

## Performance Considerations

### Optimization References
- **Common Patterns:** https://goperf.dev/01-common-patterns/
- **Networking Patterns:** https://goperf.dev/02-networking/
- **Database:** Connection pooling, prepared statements
- **HTTP:** Keep-alive, compression, caching headers

### Benchmarking
- **Testing:** `go test -bench ./... -benchmem`
- **Profiling:** pprof in standard library
- **CPU Profiling:** `go tool pprof cpu.prof`
- **Memory Profiling:** `go tool pprof mem.prof`

## Response Format

### JSend Standard
- **Specification:** https://github.com/omniti-labs/jsend
- **Status Values:** success, fail, error
- **Structure:**
  ```json
  {
    "status": "success|fail|error",
    "data": {...},
    "message": "...",
    "code": "..."
  }
  ```
- **Usage:** All API responses follow JSend format

## Deployment & CI/CD

### GitHub Actions
- **CI Pipeline:** Automated testing and linting on push
- **CD Pipeline:** Docker image build and registry push
- **Key Workflows:**
  - test.yml - Run tests on every commit
  - lint.yml - Run linters
  - security.yml - Security scanning (optional)
  - deploy.yml - Build and push Docker image

### Kubernetes (Optional)
- **Deployment Model:** Stateless microservice
- **Health Checks:** /health/live and /health/ready endpoints
- **Graceful Shutdown:** Implemented in main.go
- **Horizontal Scaling:** Stateless design supports multiple replicas

## Version Constraints

### Go Dependencies (Go Modules)
All dependencies locked to specific versions in go.mod. Key constraints:
- Go 1.25.3+ (language version)
- Fiber v2.x
- sqlc v1.x
- golang-jwt/jwt v5.x
- samber/do v2.x

### Database Versions
- **MariaDB:** 11.0+ (or MySQL 5.7+)
- **PostgreSQL:** 18.0+
- **SQLite:** 3.37+

## Notes for Implementation

1. **Dependency Management:** Run `go mod tidy` before committing
2. **Test Coverage:** Aim for >80% coverage for handler/usecase/repository layers
3. **Code Generation:** Run `go generate ./...` after adding new interfaces
4. **Migration Scripts:** Always write up and down migrations
5. **Docker Images:** Regularly update base images for security patches
6. **Configuration:** All sensitive data via environment variables, never hardcoded
