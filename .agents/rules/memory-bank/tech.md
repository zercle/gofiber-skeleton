# Technology Stack: Go Fiber Microservice Template

## Core Technology Stack

### Runtime & Framework
- **Go:** 1.25.3+ (Latest stable version with performance improvements)
- **Fiber:** v2.52.0+ (High-performance HTTP framework inspired by Express.js)
- **Architecture:** Clean Architecture with Domain-Driven Design

### Database Layer
- **Primary Database Options:**
  - **MariaDB:** 11.0+ (Small to medium systems - default choice)
  - **PostgreSQL:** 18+ (Large systems with advanced features)
  - **FerretDB:** 2.5+ (Document-based storage, MongoDB compatible)
- **Cache & Session Storage:**
  - **Valkey:** 9.0+ (Redis replacement for in-memory storage)
- **Database Tools:**
  - **sqlc:** v1.27.0+ (Type-safe SQL code generation)
  - **golang-migrate/migrate:** v4.17.0+ (Database migration management)

### Authentication & Security
- **Primary:** `golang-jwt/jwt` v5.2.0+ (JWT token management)
- **Alternative:** `zitadel/oidc` v2.3.0+ (OpenID Connect integration)
- **Security:** Rate limiting, input validation, CORS support

### Configuration & Environment
- **Configuration:** `spf13/viper` v1.18.2+ (Environment variable management)
- **Validation:** Built-in Go validation with custom validators
- **Environment Support:** Development, staging, production configurations

### Logging & Monitoring
- **Logging:** Go 1.21+ `slog` (Structured logging with context propagation)
- **Health Checks:** Custom health check endpoints
- **Metrics:** Prometheus-compatible hooks (optional integration)
- **Tracing:** OpenTelemetry integration points

### Testing & Quality
- **Mocking:** `uber-go/mock` v1.3.0+ (Interface mocking)
- **Database Testing:** `DATA-DOG/go-sqlmock` v1.5.0+ (Database mocking)
- **Linting:** `golangci-lint` v2.0+ (Comprehensive Go linting)
- **Code Generation:** `go:generate` annotations for mock generation

### Dependency Injection
- **Container:** `samber/do` v2.0.0+ (Type-safe dependency injection)
- **Utilities:**
  - `samber/lo` v1.39.0+ (Synchronous helpers for finite sequences)
  - `samber/ro` v1.4.0+ (Event-driven infinite data streams)

### Development Tools
- **Build Tool:** Go native build with Makefile orchestration
- **Containerization:** Docker with multi-stage builds
- **Documentation:** Swagger/OpenAPI 3.0 generation
- **Code Generation:** sqlc for database access, mockgen for testing

## Development Environment Setup

### Prerequisites
```bash
# Required tools versions
go version  # >= 1.25.3
docker --version  # >= 24.0.0
docker-compose --version  # >= 2.20.0
```

### Development Tools Installation
```bash
# Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.0.2
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install github.com/golang/mock/mockgen@v1.6.0
go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

# SQLC configuration
sqlc version  # Verify installation
```

### Project Initialization
```bash
# Clone and setup
git clone <repository-url>
cd template-go-fiber

# Install dependencies
go mod download
go mod verify

# Generate code
go generate ./...

# Run linting
golangci-lint run

# Setup database (local development)
docker-compose up -d mariadb
migrate -path sql/migrations -database "mysql://user:pass@tcp(localhost:3306)/dbname" up

# Run service
go run cmd/service/main.go
```

## Database Setup & Configuration

### MariaDB (Default)
```yaml
# docker-compose.yml snippet
services:
  mariadb:
    image: mariadb:11.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: template_db
      MYSQL_USER: app_user
      MYSQL_PASSWORD: app_password
    ports:
      - "3306:3306"
    volumes:
      - mariadb_data:/var/lib/mysql
```

### PostgreSQL (Alternative)
```yaml
# docker-compose.yml snippet
services:
  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_DB: template_db
      POSTGRES_USER: app_user
      POSTGRES_PASSWORD: app_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

### FerretDB (Document-based)
```yaml
# docker-compose.yml snippet
services:
  ferretdb:
    image: ghcr.io/FerretDB/FerretDB:2.5
    environment:
      FERRETDB_POSTGRESQL_URL: postgres://postgres:password@postgres:5432/ferretdb
    ports:
      - "27017:27017"
```

## Configuration Management

### Environment Variables
```bash
# .env.example - Copy to .env for local development
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_NAME=template_db
DB_USER=app_user
DB_PASSWORD=app_password
DB_DRIVER=mysql

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h

# Redis/Valkey Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

### Configuration Structure
```go
// config/types.go
type Config struct {
    Server   ServerConfig   `yaml:"server" validate:"required"`
    Database DatabaseConfig `yaml:"database" validate:"required"`
    JWT      JWTConfig      `yaml:"jwt" validate:"required"`
    Redis    RedisConfig    `yaml:"redis"`
    Logging  LoggingConfig  `yaml:"logging"`
    RateLimit RateLimitConfig `yaml:"rate_limit"`
}
```

## Build & Deployment

### Local Development
```bash
# Development build
make build-dev

# Run with hot reload (requires air)
make run-dev

# Run tests
make test

# Run integration tests
make test-integration

# Lint code
make lint

# Generate code (sqlc, mocks)
make generate
```

### Production Build
```bash
# Production build
make build-prod

# Build container image
make docker-build

# Run production container
make docker-run

# Health check
curl http://localhost:8080/health
```

### Container Configuration
```dockerfile
# Dockerfile - Multi-stage build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## Code Quality Standards

### Linting Configuration (golangci.yml)
```yaml
version: "2"
run:
  timeout: 5m
  relative-path-mode: gomod
  modules-download-mode: readonly
  allow-parallel-runners: true
  tests: true

linters-settings:
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace
```

### Testing Requirements
```bash
# Unit tests
go test -v -race -coverprofile=coverage.out ./...

# Integration tests
go test -v -tags=integration ./tests/integration/...

# Coverage report
go tool cover -html=coverage.out -o coverage.html

# Benchmark tests
go test -bench=. -benchmem ./...
```

## Performance Optimization

### Database Optimization
- Connection pooling with pgbouncer/proxy
- Query optimization with EXPLAIN ANALYZE
- Index strategy for common query patterns
- Read replica support for scaling reads

### Application Optimization
- Fiber middleware optimization
- Memory pool usage for high-frequency allocations
- Context propagation for request tracking
- Graceful shutdown with connection draining

### Caching Strategy
- Application-level caching for frequently accessed data
- Database query result caching
- HTTP response caching where appropriate
- Session state management with Redis/Valkey

## Security Best Practices

### Authentication
- JWT token rotation and refresh mechanisms
- Secure token storage (httpOnly cookies)
- Rate limiting on authentication endpoints
- Account lockout after failed attempts

### Data Protection
- Input validation and sanitization
- SQL injection prevention via sqlc
- XSS protection in API responses
- Secure headers middleware

### Infrastructure Security
- Container security scanning
- Non-root container execution
- Network segmentation
- Secrets management via environment variables

## Monitoring & Observability

### Logging Standards
```go
// Structured logging example
logger.Info("User created successfully",
    slog.String("user_id", userID),
    slog.String("email", email),
    slog.Time("created_at", time.Now()),
    slog.String("request_id", requestID),
)
```

### Health Check Endpoints
- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe
- `/health/startup` - Startup probe
- `/metrics` - Prometheus metrics (optional)

### Performance Monitoring
- Request latency tracking
- Error rate monitoring
- Database connection pool metrics
- Memory and CPU usage tracking

## Integration Points

### External Service Integration
- HTTP client with circuit breaker patterns
- Retry mechanisms with exponential backoff
- Service discovery hooks for Kubernetes/Consul
- Distributed tracing integration points

### Message Queue Integration (Future)
- Kafka integration hooks
- RabbitMQ support patterns
- Event sourcing infrastructure
- CQRS pattern implementation support

## Version Management & Updates

### Dependency Updates
- Regular security updates via `go get -u`
- Automated vulnerability scanning
- Semantic versioning for releases
- Backward compatibility considerations

### Database Migration Strategy
- Version-controlled migrations
- Rollback procedures
- Data validation post-migration
- Zero-downtime migration patterns

This technology stack provides a solid foundation for building production-ready microservices with excellent performance, maintainability, and developer experience. The stack balances modern Go practices with proven enterprise patterns.