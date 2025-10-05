# Technology Stack & Dependencies

## Core Technologies

### Go Runtime
- **Version**: Go 1.24.6 with toolchain 1.25.0
- **Features**: Latest Go features including generics, improved toolchain
- **Compilation**: Fast compilation with efficient binary generation
- **Concurrency**: Goroutines and channels for high-performance concurrency

### Web Framework
- **Framework**: Fiber v2 (github.com/gofiber/fiber/v2 v2.52.9)
- **Inspiration**: Express.js-inspired API design
- **Performance**: Built on fasthttp for high performance
- **Features**: Middleware support, routing, static file serving, WebSocket support

### Database Stack

#### Primary Database
- **Database**: PostgreSQL 18-alpine
- **Driver**: pgx v5 (github.com/jackc/pgx/v5 v5.6.0)
- **Features**: Connection pooling, prepared statements, batch operations
- **Migration Tool**: golang-migrate/migrate v4.19.0
- **Query Generation**: sqlc for type-safe SQL operations

#### Caching Layer
- **Cache**: Valkey 8-alpine (Redis-compatible)
- **Client**: go-redis/v9 (github.com/redis/go-redis/v9 v9.14.0)
- **Features**: Graceful fallback when cache unavailable
- **Use Cases**: Session storage, query result caching, rate limiting

## Dependency Injection & Architecture

### DI Framework
- **Framework**: Uber FX (go.uber.org/fx v1.24.0)
- **Features**: Constructor injection, lifecycle management, dependency graph
- **Benefits**: Loose coupling, testability, modular design
- **Pattern**: Interface-based dependency injection

### Architecture Pattern
- **Pattern**: Clean Architecture with Domain-Driven Design
- **Isolation**: Strict domain separation with zero cross-domain dependencies
- **Layers**: Entity → Repository → Usecase → Handler
- **Testing**: Mock-based testing with complete isolation

## Configuration Management

### Configuration Library
- **Library**: Viper (github.com/spf13/viper v1.21.0)
- **Features**: Environment variables, .env files, YAML support
- **Precedence**: Environment variables > .env file > defaults
- **Validation**: Automatic configuration validation and type conversion

### Configuration Structure
```go
type Config struct {
    Server   ServerConfig   // HTTP server configuration
    Database DatabaseConfig // PostgreSQL connection settings
    Redis    RedisConfig    // Cache configuration
    JWT      JWTConfig      // Authentication settings
}
```

## Authentication & Security

### Authentication
- **Library**: golang-jwt/jwt/v4 v4.5.2
- **Token Type**: JSON Web Tokens (JWT)
- **Signing Method**: HS256 with configurable secret
- **Password Hashing**: bcrypt from golang.org/x/crypto v0.42.0

### Security Features
- **Password Security**: bcrypt with DefaultCost (10)
- **Token Expiration**: Configurable (default 72 hours)
- **Middleware**: JWT validation middleware
- **Input Validation**: validator/v10 for request validation

## API Documentation & Testing

### Documentation Generation
- **Library**: swaggo/swag v1.16.6
- **UI**: fiber-swagger/v2 v2.31.1
- **Format**: OpenAPI 3.0 specification
- **Features**: Auto-generation from code comments

### Testing Framework
- **Testing**: Go's built-in testing package
- **Mocking**: 
  - go.uber.org/mock v0.6.0 (interface mocking)
  - github.com/golang/mock v1.6.0 (legacy support)
- **Assertions**: testify v1.11.1
- **Database Mocking**: DATA-DOG/go-sqlmock (planned)

## Development Tooling

### Code Generation
- **SQL Generation**: sqlc for type-safe database operations
- **Mock Generation**: go.uber.org/mock/mockgen for interface mocks
- **Documentation**: swag for API documentation
- **Migration**: golang-migrate for database schema management

### Development Tools
- **Hot Reloading**: Air (github.com/cosmtrek/air)
- **Linting**: golangci-lint with comprehensive rule set
- **Formatting**: Go's built-in fmt with goimports
- **Dependency Management**: Go modules with automatic tidying

### Build & Deployment
- **Containerization**: Docker with multi-stage builds
- **Base Image**: Alpine Linux for minimal footprint
- **Health Checks**: Built-in health check endpoints
- **Graceful Shutdown**: Proper signal handling and cleanup

## Logging & Monitoring

### Logging Framework
- **Library**: zerolog (github.com/rs/zerolog v1.34.0)
- **Features**: Structured logging, JSON output, log levels
- **Context**: Request ID tracing for distributed logging
- **Performance**: Zero-allocation logging for high performance

### Monitoring Integration
- **Health Checks**: /health and /ready endpoints
- **Metrics**: Structured logging for metrics extraction
- **Tracing**: Request ID middleware for distributed tracing
- **Error Tracking**: Structured error logging with context

## Data Validation & Serialization

### Request Validation
- **Library**: go-playground/validator/v10 v10.27.0
- **Features**: Custom validation rules, error messages
- **Integration**: Middleware integration with Fiber
- **Types**: Struct tag-based validation

### Data Serialization
- **JSON**: Fiber's built-in JSON handling
- **Response Format**: JSend standard for consistent API responses
- **Error Handling**: Centralized error response formatting
- **Content Types**: Automatic content negotiation

## Middleware Stack

### Core Middleware
- **Recovery**: Panic recovery with proper error logging
- **Request ID**: Unique request identification
- **Logging**: Structured request/response logging
- **CORS**: Cross-Origin Resource Sharing configuration
- **Security**: Security headers and best practices

### Rate Limiting
- **Implementation**: Custom rate limiting middleware
- **Storage**: Redis-based with in-memory fallback
- **Configuration**: Configurable limits per endpoint
- **Strategy**: Token bucket algorithm

## Database Schema & Migrations

### Migration Strategy
- **Tool**: golang-migrate/migrate
- **Versioning**: Sequential version numbering
- **Rollback**: Down migration support
- **Environment**: Separate migration sets per environment

### Schema Design
- **Primary Keys**: UUID with gen_random_uuid()
- **Audit Fields**: created_at, updated_at timestamps
- **Indexes**: Optimized indexes for common queries
- **Constraints**: Foreign key constraints for data integrity

## Performance Optimizations

### Database Optimizations
- **Connection Pooling**: Configurable pool settings
- **Query Optimization**: sqlc-generated efficient queries
- **Batch Operations**: Support for bulk operations
- **Indexing Strategy**: Optimized indexes for performance

### Application Optimizations
- **Memory Management**: Efficient memory usage patterns
- **Concurrency**: Goroutine pools for CPU-bound tasks
- **Caching**: Redis integration for frequently accessed data
- **Compression**: Gzip compression for responses

## Security Implementation

### Input Security
- **Validation**: Comprehensive input validation
- **Sanitization**: Input sanitization for XSS prevention
- **SQL Injection**: Parameterized queries via sqlc
- **File Upload**: Secure file handling with validation

### Authentication Security
- **Password Hashing**: bcrypt with appropriate cost factor
- **JWT Security**: Secure token generation and validation
- **Session Management**: Secure session handling
- **Rate Limiting**: Brute force attack prevention

### Infrastructure Security
- **HTTPS**: TLS configuration for production
- **Headers**: Security headers (HSTS, CSP, etc.)
- **Environment**: Environment-based secret management
- **Dependencies**: Regular security updates

## Development Environment

### Local Development
- **Containerization**: Docker Compose for local environment
- **Database**: PostgreSQL container with persistent storage
- **Cache**: Valkey container for Redis-compatible caching
- **Hot Reload**: Air for automatic server restarts

### Code Quality Tools
- **Linting**: golangci-lint with 50+ rules
- **Formatting**: Automatic code formatting
- **Testing**: Comprehensive test suite with coverage
- **Documentation**: Auto-generated API documentation

## Dependency Management

### Go Modules
- **Module**: github.com/zercle/gofiber-skeleton
- **Versioning**: Semantic versioning for releases
- **Updates**: Regular dependency updates
- **Vulnerability Scanning**: Automated security scanning

### Key Dependencies
```
github.com/gofiber/fiber/v2 v2.52.9      // Web framework
go.uber.org/fx v1.24.0                   // Dependency injection
github.com/spf13/viper v1.21.0           // Configuration
github.com/jackc/pgx/v5 v5.6.0           // PostgreSQL driver
github.com/golang-migrate/migrate/v4 v4.19.0 // Migrations
github.com/sqlc-dev/sqlc v0.0.0-202401... // SQL generation
github.com/golang-jwt/jwt/v4 v4.5.2      // JWT authentication
github.com/rs/zerolog v1.34.0            // Structured logging
github.com/swaggo/swag v1.16.6           // API documentation
go.uber.org/mock v0.6.0                  // Mock generation
```

## Build & Deployment Configuration

### Build Process
- **Compilation**: Go build with optimization flags
- **Containerization**: Multi-stage Docker builds
- **Artifact Management**: Versioned binary generation
- **Environment**: Environment-specific configurations

### Deployment Strategy
- **Containers**: Docker containerization
- **Orchestration**: Kubernetes-ready configuration
- **Health Checks**: Built-in health monitoring
- **Scaling**: Horizontal scaling support

## Monitoring & Observability

### Logging Strategy
- **Structured Logging**: JSON format for log aggregation
- **Log Levels**: Configurable logging levels
- **Context**: Request tracing with unique IDs
- **Performance**: Low-overhead logging implementation

### Health Monitoring
- **Liveness**: Basic application health checks
- **Readiness**: Database and cache connectivity checks
- **Metrics**: Performance metrics collection
- **Alerting**: Integration-ready alerting hooks

## Version Control & CI/CD

### Git Workflow
- **Branching**: Feature branch workflow
- **Commits**: Conventional commit messages
- **Tags**: Semantic versioning tags
- **Releases**: Automated release process

### Quality Gates
- **Testing**: All tests must pass
- **Linting**: Zero linting errors
- **Coverage**: Minimum test coverage requirements
- **Security**: Security vulnerability scanning

## Future Technology Considerations

### Potential Enhancements
- **GraphQL**: GraphQL integration for flexible APIs
- **gRPC**: gRPC for internal service communication
- **Message Queues**: Event-driven architecture support
- **Search Engine**: Elasticsearch integration for search

### Performance Improvements
- **Caching**: Advanced caching strategies
- **Database**: Read replicas for scaling
- **CDN**: Content delivery network integration
- **Load Balancing**: Advanced load balancing strategies