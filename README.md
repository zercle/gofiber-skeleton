# Go Fiber Skeleton

A production-ready Go backend template built with **Fiber v2**, featuring **Clean Architecture**, **Domain-Driven Design**, and comprehensive development tooling.

## 🚀 Features

- **Clean Architecture** - Domain-isolated, testable, and maintainable codebase
- **Fiber v2 Framework** - High-performance Express.js-inspired web framework
- **Uber fx** - Powerful dependency injection framework
- **Viper Configuration** - Flexible config management (env vars, .env files, YAML)
- **Type-Safe Database** - sqlc generates fully type-safe Go code from SQL
- **Database Migrations** - golang-migrate for versioned schema management
- **Authentication** - JWT-based auth with bcrypt password hashing
- **Redis Integration** - Ready-to-use caching layer with go-redis
- **Comprehensive Testing** - Unit tests with mocks (mockgen + go-sqlmock)
- **API Documentation** - Auto-generated Swagger/OpenAPI docs
- **Hot Reloading** - Air for fast development cycles
- **Production-Ready** - Docker, health checks, structured logging, graceful shutdown
- **Security** - CORS, rate limiting, security headers, input validation

## 📋 Prerequisites

- **Go** 1.24.6 or higher
- **Docker** and **Docker Compose**
- **Make** (for running commands)

## 🏃 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/zercle/gofiber-skeleton.git
cd gofiber-skeleton
```

### 2. Install Development Tools

```bash
make install-tools
```

This installs:
- Air (hot-reloading)
- sqlc (SQL code generation)
- swag (Swagger docs)
- mockgen (mock generation)
- golangci-lint (linting)
- migrate (database migrations)

### 3. Set Up Environment Variables

```bash
cp .env.example .env
```

Edit `.env` if needed. Default values work with Docker Compose.

### 4. Start Services with Docker Compose

```bash
docker-compose up -d
```

This starts:
- PostgreSQL database (port 5432)
- Redis/Valkey cache (port 6379)
- Application server (port 8080)

### 5. Run Database Migrations

```bash
make migrate-up
```

### 6. Start Development Server

```bash
make dev
```

The API will be available at `http://localhost:8080`.

## 📚 API Documentation

Access interactive Swagger UI at: `http://localhost:8080/swagger/`

## 🏗️ Project Structure

```
gofiber-skeleton/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point with fx setup
├── internal/
│   ├── cache/
│   │   └── redis.go             # Redis client wrapper
│   ├── config/
│   │   └── config.go            # Viper-based configuration
│   ├── database/
│   │   └── migrate.go           # Migration helpers
│   ├── db/                      # sqlc-generated code
│   ├── errors/
│   │   └── errors.go            # Custom error types
│   ├── logger/
│   │   └── logger.go            # Zerolog structured logger
│   ├── middleware/
│   │   ├── cors.go              # CORS middleware
│   │   ├── logging.go           # Request logging
│   │   ├── rate_limit.go        # Rate limiting
│   │   ├── request_id.go        # Request ID generation
│   │   └── security.go          # Security headers
│   ├── response/
│   │   └── jsend.go             # JSend response format
│   ├── server/
│   │   └── router.go            # Route registration
│   ├── user/                    # User domain (example)
│   │   ├── entity/
│   │   ├── repository/
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── tests/
│   └── post/                    # Post domain (example)
│       ├── entity/
│       ├── repository/
│       ├── usecase/
│       ├── handler/
│       └── tests/
├── pkg/
│   └── validator/               # Request validation utilities
├── db/
│   ├── migrations/              # SQL migration files
│   └── queries/                 # sqlc query definitions
├── docs/                        # Swagger-generated documentation
│   └── ADDING_NEW_DOMAIN.md     # Guide for adding new domains
├── .air.toml                    # Air configuration
├── .env.example                 # Environment variables template
├── compose.yml                  # Docker Compose configuration
├── Dockerfile                   # Multi-stage production build
├── Makefile                     # Development commands
├── sqlc.yaml                    # sqlc configuration
└── go.mod                       # Go module definition
```

## 🛠️ Development Commands

```bash
# Development
make dev              # Run with hot-reloading
make build            # Build binary
make run              # Build and run

# Code Quality
make fmt              # Format code
make lint             # Run linter
make test             # Run tests
make test-race        # Run tests with race detector
make test-coverage    # Generate coverage report

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback last migration
make migrate-create name=<name>  # Create new migration
make sqlc             # Generate type-safe DB code

# Documentation
make generate-docs    # Generate Swagger docs
make generate-mocks   # Generate test mocks

# Utilities
make clean            # Clean build artifacts
make ci               # Run full CI pipeline
make help             # Show all commands
```

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run with Race Detector

```bash
make test-race
```

### Generate Coverage Report

```bash
make test-coverage
```

This creates `coverage.html` in the project root.

### Test Structure

- **Unit Tests**: Test usecases with mocked repositories
- **Repository Tests**: Test data access with go-sqlmock
- **Integration Tests**: Test complete flows

Example test locations:
- `internal/post/tests/post_usecase_test.go`
- `internal/user/tests/user_repository_test.go`

## 📦 Adding a New Domain

Follow the comprehensive guide: [docs/ADDING_NEW_DOMAIN.md](docs/ADDING_NEW_DOMAIN.md)

Quick steps:
1. Create directory structure
2. Define entity
3. Create database migration
4. Write SQL queries
5. Generate sqlc code
6. Implement repository
7. Implement usecase
8. Implement handlers
9. Write tests
10. Register routes

## 🔑 Authentication

### Register a New User

```bash
POST /api/v1/auth/register
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

### Login

```bash
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Returns a JWT token valid for 72 hours.

### Protected Routes

Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 🐳 Docker Deployment

### Build Production Image

```bash
docker build -t gofiber-skeleton:latest .
```

### Run with Docker Compose

```bash
docker-compose up -d
```

### Environment Variables

Configure via environment variables in production:

```bash
SERVER_PORT=8080
SERVER_ENV=production
DATABASE_DSN=postgres://...
JWT_SECRET=your-secret-key
REDIS_ADDR=redis:6379
```

## 🔧 Configuration

Configuration is managed via **Viper** with the following priority:

1. Environment variables (highest priority)
2. `.env` file
3. Default values (lowest priority)

### Configuration Structure

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
}
```

See `.env.example` for all available options.

## 📊 Database Migrations

### Create a New Migration

```bash
make migrate-create name=add_users_table
```

This creates two files in `db/migrations/`:
- `NNNN_add_users_table.up.sql`
- `NNNN_add_users_table.down.sql`

### Run Migrations

```bash
make migrate-up
```

### Rollback Last Migration

```bash
make migrate-down
```

### Check Migration Version

```bash
make migrate-version
```

## 🔍 Logging

Structured JSON logging with **zerolog**:

```go
import "github.com/zercle/gofiber-skeleton/internal/logger"

log := logger.GetLogger()
log.Info().Msg("Something happened")
log.Error().Err(err).Msg("Something failed")
```

All HTTP requests are automatically logged with:
- Request ID
- Method and path
- Status code
- Duration
- Client IP

## 🛡️ Security Features

- **bcrypt** password hashing (cost factor 10)
- **JWT** authentication with expiration
- **CORS** protection
- **Rate limiting** (API and auth endpoints)
- **Security headers** (XSS, clickjacking, etc.)
- **Input validation** with go-playground/validator
- **Request ID** tracking for debugging

## 🚀 Performance

- **Connection pooling** for database
- **Redis caching** ready
- **Graceful shutdown** (30-second timeout)
- **Health checks** for liveness and readiness probes
- **Static binary** compilation for fast startup

## 📝 API Response Format

All API responses follow the **JSend** specification:

### Success

```json
{
  "status": "success",
  "data": { ... }
}
```

### Fail (Client Error)

```json
{
  "status": "fail",
  "data": {
    "field": "error message"
  }
}
```

### Error (Server Error)

```json
{
  "status": "error",
  "message": "Something went wrong",
  "code": 5000
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [Uber fx](https://uber-go.github.io/fx/) - Dependency injection
- [sqlc](https://sqlc.dev/) - SQL code generation
- [Viper](https://github.com/spf13/viper) - Configuration
- [golang-migrate](https://github.com/golang-migrate/migrate) - Migrations
- [zerolog](https://github.com/rs/zerolog) - Logging

## 📞 Support

For issues and questions:
- GitHub Issues: [https://github.com/zercle/gofiber-skeleton/issues](https://github.com/zercle/gofiber-skeleton/issues)

---

**Built with ❤️ using Go and Fiber**
