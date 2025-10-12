# **Technology Stack Documentation**

## **Core Technologies**

### **Go Language & Runtime**
- **Version**: Go 1.25.0
- **Module**: `github.com/zercle/gofiber-skeleton`
- **Features**: Generics, goroutines, garbage collection
- **Benefits**: Performance, concurrency, type safety

### **Web Framework: Fiber v2**
- **Package**: `github.com/gofiber/fiber/v2`
- **Architecture**: Express.js-inspired, high-performance
- **Features**: 
  - Zero-allocation HTTP router
  - Middleware support
  - WebSocket support
  - Static file serving
  - Template engine integration
- **Performance**: 10x faster than Express.js

## **Dependency Injection**

### **Samber's Do**
- **Package**: `github.com/samber/do`
- **Type**: Generics-based DI container
- **Features**:
  - Type-safe dependency resolution
  - Automatic dependency graph resolution
  - Lifecycle management
  - Interface-based design
- **Benefits**: Compile-time safety, clean architecture

## **Database Technologies**

### **Database: PostgreSQL**
- **Version**: 18-alpine (Docker)
- **Features**:
  - ACID compliance
  - JSON support
  - Full-text search
  - Connection pooling
- **Connection**: `pgx` driver for Go

### **Database Migration: golang-migrate**
- **Package**: `github.com/golang-migrate/migrate`
- **Features**:
  - Version-controlled migrations
  - Rollback capabilities
  - Multiple database support
  - CLI tool integration
- **Migration Files**: `db/migrations/`

### **SQL Generation: sqlc**
- **Package**: `github.com/sqlc-dev/sqlc`
- **Features**:
  - Compile-time SQL validation
  - Type-safe query generation
  - Auto-generated Go code
  - IDE support with autocomplete
- **Query Files**: `db/queries/`
- **Generated Code**: `internal/domains/*/repository/`

## **Caching Technology**

### **Cache: Valkey (Redis-compatible)**
- **Image**: `valkey/valkey:latest`
- **Features**:
  - Redis-compatible API
  - In-memory data structure store
  - Pub/sub capabilities
  - Persistence options
- **Go Client**: `github.com/redis/go-redis/v9`

## **Authentication & Security**

### **JWT: golang-jwt**
- **Package**: `github.com/golang-jwt/jwt/v5`
- **Features**:
  - Token generation and validation
  - Multiple signing algorithms
  - Claims customization
  - Middleware integration

### **Password Hashing: Argon2id**
- **Package**: `golang.org/x/crypto/argon2`
- **Features**:
  - Memory-hard hashing function
  - Resistance to GPU/ASIC attacks
  - Configurable time, memory, and parallelism
  - Recommended by OWASP for password hashing
  - Built-in Go library

## **Configuration Management**

### **Viper**
- **Package**: `github.com/spf13/viper`
- **Features**:
  - Multiple configuration sources
  - Environment variable support
  - .env file support
  - YAML/JSON/TOML support
  - Live configuration reloading
- **Precedence**: Environment > .env > defaults

## **API Documentation**

### **Swagger: swaggo/swag**
- **Package**: `github.com/swaggo/swag`
- **Features**:
  - Automatic OpenAPI specification
  - Interactive UI documentation
  - Code-based documentation
  - Schema generation
- **Output**: `docs/swagger.json`, `docs/swagger.html`
- **Web UI**: `/swagger/` endpoint

## **Development Tools**

### **Hot Reloading: Air**
- **Package**: `github.com/cosmtrek/air`
- **Features**:
  - Automatic server restart
  - File watching
  - Build optimization
  - Configuration via `.air.toml`

### **Code Quality: golangci-lint**
- **Package**: `github.com/golangci/golangci-lint`
- **Features**:
  - 50+ linters integrated
  - Fast parallel execution
  - Configuration via `.golangci.yml`
  - CI/CD integration

## **Testing Framework**

### **Mock Generation: uber-go/mock**
- **Package**: `go.uber.org/mock/mockgen`
- **Features**:
  - Interface-based mock generation
  - `//go:generate` integration
  - Customizable mock behavior
  - Test isolation support
  - Comprehensive mocking for unit tests

### **SQL Mocking: go-sqlmock**
- **Package**: `github.com/DATA-DOG/go-sqlmock`
- **Features**:
  - SQL driver mocking
  - Query expectation matching
  - Transaction simulation
  - Database-independent testing
  - Primary tool for database layer testing

## **Containerization**

### **Docker**
- **Base Images**: 
  - `golang:1.25-alpine` (build stage)
  - `alpine:latest` (runtime stage)
- **Features**:
  - Multi-stage builds
  - Minimal runtime image
  - Security scanning
  - Layer optimization

### **Docker Compose**
- **File**: `compose.yml`
- **Services**:
  - PostgreSQL database
  - Valkey cache
  - Application container
- **Features**:
  - Development environment
  - Service orchestration
  - Volume management
  - Network configuration

## **CI/CD Pipeline**

### **GitHub Actions**
- **Workflows**:
  - `ci.yml`: Continuous integration
  - `cd.yml`: Continuous deployment
  - `go-ci.yml`: Go-specific pipeline
- **Features**:
  - Automated testing
  - Code quality checks
  - Security scanning
  - Docker image building
  - Multi-environment deployment

### **Code Quality Tools**
- **Gosec**: Security vulnerability scanning
- **Trivy**: Container security scanning
- **Codecov**: Test coverage reporting

## **Monitoring & Observability**

### **Structured Logging**
- **Package**: `github.com/gofiber/fiber/v2/middleware/logger`
- **Features**:
  - Request/response logging
  - Structured log format
  - Configurable log levels
  - Output customization

### **Request Tracing**
- **Package**: Custom middleware
- **Features**:
  - Request ID generation
  - Request propagation
  - Performance tracking
  - Debugging support

### **Health Checks**
- **Package**: Custom health endpoints
- **Features**:
  - Database connectivity
  - Cache connectivity
  - Application status
  - Dependency health

## **Response Handling**

### **JSend Response Format**
- **Package**: Custom response utilities
- **Features**:
  - Standardized response format
  - Success/error handling
  - Data pagination
  - Consistent API responses

### **Error Handling**
- **Package**: Custom error middleware
- **Features**:
  - Centralized error handling
  - Error categorization
  - Stack trace management
  - User-friendly error messages

## **Performance Optimizations**

### **Connection Pooling**
- **Database**: `pgxpool` for PostgreSQL
- **Cache**: Redis connection pool
- **Features**:
  - Configurable pool size
  - Connection lifecycle management
  - Performance monitoring
  - Resource optimization

### **Middleware Stack**
- **CORS**: Cross-origin resource sharing
- **Rate Limiting**: Request rate control
- **Compression**: Gzip response compression
- **Recovery**: Panic recovery middleware

## **Development Environment**

### **Required Tools**
```bash
# Go runtime
go version 1.25.0+

# Development tools
make - build automation
docker - containerization
docker-compose - local development
air - hot reloading
sqlc - SQL code generation
swag - API documentation
golangci-lint - code quality
mockgen - mock generation
```

### **Development Commands**
```bash
make dev          # Start development server
make test         # Run test suite
make lint         # Code quality checks
make build        # Build production binary
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations
make sqlc         # Generate SQL code
make swag         # Generate API docs
make mocks        # Generate test mocks
```

## **Environment Configuration**

### **Required Environment Variables**
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gofiber_skeleton
DB_SSLMODE=disable

# Cache
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_PASSWORD=""
VALKEY_DB=0

# Authentication
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=3000
SERVER_ENV=development
```

### **Configuration Files**
- `.env`: Local development variables
- `config.yml`: Default configuration
- `docker-compose.yml`: Container orchestration
- `.air.toml`: Hot reload configuration
- `.golangci.yml`: Linting rules

## **Security Considerations**

### **Dependencies**
- Regular security updates via `go get -u`
- Vulnerability scanning with `gosec` and `Trivy`
- Dependency version pinning in `go.mod`
- SBOM generation for compliance

### **Runtime Security**
- Input validation and sanitization
- SQL injection prevention with sqlc
- XSS protection via Fiber middleware
- Rate limiting and DDoS protection
- Secure headers configuration

## **Performance Benchmarks**

### **Expected Performance**
- **Throughput**: 100,000+ requests/second
- **Latency**: < 1ms average response time
- **Memory**: < 100MB baseline memory usage
- **CPU**: Efficient goroutine utilization

### **Monitoring Metrics**
- Request rate and response times
- Database query performance
- Cache hit ratios
- Memory and CPU usage
- Error rates and patterns

## **Technology Rationale**

### **Why Go?**
- Performance and efficiency
- Strong typing and safety
- Excellent concurrency support
- Rich standard library
- Growing ecosystem

### **Why Fiber?**
- Express.js familiarity
- High performance
- Rich middleware ecosystem
- Comprehensive documentation
- Active development

### **Why Clean Architecture?**
- Testability and maintainability
- Domain-driven design
- Separation of concerns
- Framework independence
- Long-term sustainability

## **Future Technology Considerations**

### **Potential Enhancements**
- **GraphQL**: Alternative to REST APIs
- **gRPC**: High-performance RPC framework
- **Kubernetes**: Container orchestration
- **Prometheus**: Metrics collection
- **Jaeger**: Distributed tracing
- **Vault**: Secret management

### **Technology Debt Management**
- Regular dependency updates
- Technology lifecycle planning
- Performance monitoring and optimization
- Security audit scheduling
- Code review processes