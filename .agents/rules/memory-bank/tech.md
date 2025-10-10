# **Technology Stack: Go Fiber Production-Ready Template**

## **Core Technology Stack**

### **Primary Framework**
- **Go:** Version 1.25.0 (latest stable)
- **Web Framework:** Fiber v2 - Express.js-inspired, high-performance HTTP framework
- **Architecture:** Clean Architecture with Domain-Driven Design principles

### **Dependency Management**
- **Go Modules:** Native Go dependency management
- **Module Name:** github.com/zercle/gofiber-skeleton
- **Version Strategy:** Semantic versioning with go.mod pinning

## **Key Dependencies and Libraries**

### **Web Framework and Routing**
```go
// Core framework
"github.com/gofiber/fiber/v2" v2.52.4

// Middleware and utilities
"github.com/gofiber/fiber/v2/middleware/cors"
"github.com/gofiber/fiber/v2/middleware/logger"
"github.com/gofiber/fiber/v2/middleware/recover"
"github.com/gofiber/fiber/v2/middleware/limiter"
```

### **Dependency Injection**
```go
// Modern dependency injection with generics
"github.com/samber/do/v2" v2.19.0
```

### **Configuration Management**
```go
// Configuration and environment variables
"github.com/spf13/viper" v1.18.2
```

### **Database and ORM**
```go
// PostgreSQL driver
"github.com/lib/pq" v1.10.9

// Database migrations
"github.com/golang-migrate/migrate/v4"
"github.com/golang-migrate/migrate/v4/database/postgres"
"github.com/golang-migrate/migrate/v4/source/file"

// Type-safe SQL generation
"github.com/sqlc-dev/sqlc" v1.25.0
```

### **Authentication and Security**
```go
// JWT token handling
"github.com/golang-jwt/jwt/v5" v5.2.0

// Password hashing
"golang.org/x/crypto/bcrypt"
```

### **API Documentation**
```go
// Swagger/OpenAPI generation
"github.com/swaggo/swag" v1.16.3
"github.com/swaggo/fiber-swagger" v1.1.1
```

### **Testing and Mocking**
```go
// Testing framework (built-in)
"testing"

// Mock generation
"go.uber.org/mock/mockgen" v0.4.0

// Test assertions
"github.com/stretchr/testify/assert" v1.8.4
"github.com/stretchr/testify/suite" v1.8.4
```

### **Validation**
```go
// Input validation
"github.com/go-playground/validator/v10" v10.16.0
```

### **Utilities**
```go
// UUID generation
"github.com/google/uuid" v1.4.0

// Time handling
"github.com/golang-module/carbon/v2" v2.2.2

// Error handling
"github.com/pkg/errors" v0.9.1
```

## **Development Tooling**

### **Code Quality and Linting**
- **Tool:** golangci-lint v1.56.2
- **Configuration:** .golangci.yml with comprehensive rules
- **Integration:** Pre-commit hooks and CI/CD pipeline

### **Hot Reloading**
- **Tool:** Air v1.49.0
- **Configuration:** .air.toml for development file watching
- **Integration:** Makefile `make dev` command

### **Database Tools**
- **Migration Tool:** golang-migrate/migrate CLI
- **Query Generation:** sqlc CLI for Go code generation
- **Database Client:** psql or any PostgreSQL client

### **API Documentation**
- **Generation:** swag CLI for Swagger documentation
- **UI:** fiber-swagger for interactive API docs
- **Endpoint:** /swagger/ for development access

### **Containerization**
- **Container Runtime:** Docker v24.0+
- **Orchestration:** Docker Compose v2.0+
- **Base Images:** Alpine Linux for production efficiency

## **Database Technology**

### **Primary Database**
- **Database:** PostgreSQL 15+
- **Driver:** lib/pq (pure Go PostgreSQL driver)
- **Connection Pooling:** Built-in driver connection pooling
- **Migrations:** Version-controlled SQL migration files

### **Database Schema Management**
```sql
-- Migration files location: db/migrations/
-- Pattern: 001_initial_schema.up.sql, 001_initial_schema.down.sql
-- Tool: golang-migrate/migrate
```

### **Query Generation**
- **Tool:** sqlc for compile-time type safety
- **Input:** Raw SQL queries in db/queries/
- **Output:** Generated Go code in internal/domains/*/repository/sqlc/
- **Benefits:** Compile-time SQL validation, type safety

### **Caching Layer**
- **Cache:** Valkey (Redis-compatible)
- **Driver:** go-redis/redis/v9
- **Use Cases:** Session storage, query result caching
- **Configuration:** Docker Compose service definition

## **API and Documentation**

### **REST API Framework**
- **Router:** Fiber v2 built-in router
- **Middleware Stack:** CORS, logging, recovery, rate limiting
- **Request/Response:** JSON with proper content negotiation
- **Error Handling:** Consistent error response format

### **API Documentation**
- **Standard:** OpenAPI 3.0 (Swagger)
- **Generation:** Automatic from Go comments
- **Interactive UI:** Swagger UI integrated via fiber-swagger
- **Development Access**: http://localhost:3000/swagger/

### **Authentication**
- **Scheme:** JWT Bearer tokens
- **Algorithm:** HS256 for server-to-server, RS256 for external
- **Token Storage**: Authorization header
- **Refresh Strategy**: Optional refresh token implementation

## **Development Environment**

### **Required Development Tools**
```bash
# Go toolchain
go version 1.25.0+

# Container tools
docker --version
docker-compose --version

# Database tools
migrate --version
sqlc version

# Development tools
air --version
golangci-lint --version
swag --version
```

### **Development Dependencies**
```bash
# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install go.uber.org/mock/mockgen@latest
```

### **Environment Configuration**
```bash
# Development environment variables
GO_ENV=development
PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gofiber_skeleton
JWT_SECRET=development-secret
REDIS_HOST=localhost
REDIS_PORT=6379
```

## **Production Environment**

### **Container Configuration**
```dockerfile
# Multi-stage Docker build
# Stage 1: Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Stage 2: Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### **Production Environment Variables**
```bash
# Production configuration
GO_ENV=production
PORT=8080
DB_HOST=$DATABASE_URL
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
JWT_SECRET=$JWT_SECRET
REDIS_HOST=$REDIS_URL
REDIS_PORT=6379
```

### **Performance Considerations**
- **Connection Pooling:** Optimized database connection settings
- **Caching Strategy:** Redis for frequently accessed data
- **Compression:** Gzip middleware for response compression
- **Rate Limiting:** Built-in Fiber rate limiting middleware

## **Testing Framework**

### **Unit Testing**
- **Framework:** Go built-in testing package
- **Assertions:** Testify for readable assertions
- **Mocks:** Uber mock for interface mocking
- **Coverage:** Target 90%+ for business logic

### **Integration Testing**
- **Database:** Test containers or test database
- **API Testing:** HTTP endpoint testing with test suite
- **Environment:** Isolated test environment setup

### **Testing Commands**
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Generate mocks
make generate-mocks

# Run integration tests
make test-integration
```

## **Build and Deployment**

### **Build Process**
```bash
# Build for development
make build

# Build for production
make build-prod

# Cross-platform builds
make build-all
```

### **Deployment Configuration**
- **Container Registry:** Docker Hub or private registry
- **Orchestration:** Kubernetes or Docker Swarm
- **Health Checks:** Built-in Fiber health check endpoints
- **Monitoring:** Structured logging and metrics

### **CI/CD Pipeline Components**
- **Linting:** golangci-lint in CI pipeline
- **Testing:** Automated test execution
- **Security:** Dependency scanning and vulnerability checks
- **Build:** Multi-stage Docker builds
- **Deploy:** Automated deployment to staging/production

## **Monitoring and Observability**

### **Logging**
- **Library:** Fiber built-in logger with structured output
- **Levels:** Debug, Info, Warn, Error
- **Format:** JSON in production, human-readable in development
- **Destination:** Standard output (container-friendly)

### **Health Checks**
- **Endpoint:** /health for basic health status
- **Endpoint:** /health/ready for readiness probe
- **Endpoint:** /health/live for liveness probe
- **Database:** Connectivity and query validation

### **Metrics (Future Enhancement)**
- **Library:** Prometheus integration planned
- **Endpoints:** /metrics for Prometheus scraping
- **Custom Metrics:** Business-specific metrics collection
- **Alerting:** Integration with monitoring systems

## **Technical Constraints and Limitations**

### **Platform Requirements**
- **Go Version:** 1.25+ required for generics support
- **Operating System:** Linux, macOS, Windows
- **Architecture:** amd64, arm64 supported
- **Memory:** Minimum 512MB for development, 1GB+ for production

### **Database Requirements**
- **PostgreSQL:** Version 12+ required
- **Storage:** SSD recommended for performance
- **Connection:** Standard TCP connection
- **Extensions:** No special extensions required

### **Network Requirements**
- **Ports:** Configurable HTTP port (default: 3000)
- **Firewall:** Outbound connections for dependencies
- **SSL/TLS:** HTTPS termination at load balancer
- **DNS:** Standard DNS resolution required

## **Security Considerations**

### **Built-in Security Features**
- **Authentication:** JWT-based stateless authentication
- **Password Security:** bcrypt hashing with salt
- **Input Validation:** Go-playground validator
- **CORS:** Configurable cross-origin resource sharing
- **Rate Limiting:** Built-in Fiber middleware
- **Security Headers:** Standard HTTP security headers

### **Security Best Practices**
- **Environment Variables:** Sensitive data in environment
- **SQL Injection Prevention:** sqlc compile-time validation
- **XSS Prevention:** Input sanitization and output encoding
- **CSRF Protection:** Token-based CSRF protection (future)
- **Dependency Updates:** Regular security patch updates