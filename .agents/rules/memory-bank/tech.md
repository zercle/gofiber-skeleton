# Technology Stack

## Core Technologies

### Web Framework
- **Go Fiber v2**: Fast HTTP framework built on Fasthttp
- **Usage**: REST API development, routing, middleware
- **Benefits**: High performance, Express.js-like API, built-in utilities

### Dependency Injection
- **Uber FX**: Dependency injection framework for Go
- **Usage**: Application lifecycle management, dependency wiring
- **Benefits**: Type-safe DI, lifecycle hooks, testing support

### Configuration Management
- **Viper**: Configuration management library
- **Sources**: ./config/<env>.yaml files, .env files, environment variables
- **Hierarchy**: Environment-specific YAML → .env file → runtime environment variables (highest priority)
- **Usage**: Multi-environment configuration, runtime config changes

## Database Stack

### Database
- **PostgreSQL**: Primary relational database
- **Connection**: pgx driver for performance
- **Pooling**: Built-in connection pooling

### SQL Generation
- **sqlc**: Compile-time SQL query generation
- **Benefits**: Type-safe SQL, generated Go code, query validation
- **Usage**: Repository implementations, data access layer

### Migrations
- **golang-migrate/migrate**: Database migration tool
- **Features**: Version control, rollback support, multiple database support
- **Usage**: Schema evolution, database setup automation

## Development Tools

### Hot Reloading
- **Air**: Live reloading for Go applications
- **Configuration**: `air.toml` for custom watch patterns
- **Usage**: Development workflow, automatic rebuilds

### Code Generation
- **Uber Mock (go.uber.org/mock)**: Mock generation
- **Usage**: Interface mocking, unit testing, dependency isolation
- **Command**: `go generate ./...`

### Testing
- **Go Testing**: Standard Go testing framework
- **Testify**: Testing toolkit with assertions and mocks
- **go-sqlmock**: SQL driver mock for database testing

## Security & Authentication

### JWT Authentication
- **golang-jwt**: JSON Web Token implementation
- **Usage**: Stateless authentication, API security
- **Features**: Token generation, validation, middleware

### Validation
- **go-playground/validator**: Struct and field validation
- **Usage**: Input validation, business rule enforcement
- **Features**: Tag-based validation, custom validators

## API & Documentation

### API Response Format
- **omniti-labs/jsend**: Standardized JSON response format
- **Structure**: Consistent success/error response patterns
- **Usage**: API response standardization

### Documentation
- **swaggo/swag**: Swagger documentation generation
- **Usage**: API documentation from Go annotations
- **Output**: OpenAPI/Swagger JSON/YAML

## Unique Identifiers
- **UUIDv7**: Time-ordered UUIDs for database-friendly indexing
- **Benefits**: Index-friendly, time-sortable, globally unique
- **Usage**: Primary keys, entity identification

## Development Environment

### Go Version
- **Go 1.24.6**: Latest Go version for modern features
- **Modules**: Go modules for dependency management
- **Features**: Generics, improved performance, security updates

### CI/CD
- **GitHub Actions**: Automated testing and linting
- **Workflow**: 
  - Go mod download
  - Mock generation
  - golangci-lint for code quality
  - Test execution with race detection

### Code Quality
- **golangci-lint**: Comprehensive Go linting
- **gofmt**: Code formatting
- **go vet**: Static analysis

## Containerization

### Docker
- **Multi-stage builds**: Optimized production images
- **Development**: compose for local development
- **Services**: PostgreSQL, Redis (optional), application container

### Production
- **Minimal base images**: Alpine Linux for security and size
- **Health checks**: Built-in health endpoints
- **Configuration**: Environment-based configuration

## Configuration Loading Hierarchy

**Priority Order (highest to lowest):**
1. **Runtime Environment Variables** - Override everything
2. **.env File** - Environment-specific overrides
3. **./config/<env>.yaml** - Environment-specific YAML configuration
4. **Default Values** - Fallback defaults in code

### Configuration File Structure
- **./config/development.yaml** - Development environment settings
- **./config/staging.yaml** - Staging environment settings  
- **./config/production.yaml** - Production environment settings
- **.env** - Local environment overrides (not committed)

### Required Environment Variables
```bash
PORT=3000                    # Server port
ENV=development             # Environment (development/staging/production)
DATABASE_URL=postgres://... # PostgreSQL connection string
JWT_SECRET=your-secret      # JWT signing secret
JWT_EXPIRES_IN=24h          # Token expiration time
REDIS_URL=redis://...       # Redis connection (optional)
CORS_ORIGINS=*              # CORS allowed origins
```

### Configuration Files
- **./config/<env>.yaml**: Environment-specific YAML configuration files (e.g., config/development.yaml, config/production.yaml)
- **.env**: Environment-specific overrides loaded after YAML
- **Environment variables**: Runtime configuration with highest precedence

## Development Commands

### Core Commands
```bash
# Development
go run cmd/server/main.go    # Start server
air                          # Hot reload development

# Database
go run cmd/migrate/main.go   # Run migrations

# Testing
go test ./...                # Run all tests
go test -race ./...          # Run with race detection

# Code Quality
golangci-lint run           # Lint codebase
gofmt -s -w .               # Format code

# Documentation
swag init -g cmd/server/main.go -o docs  # Generate API docs

# Build
go build -o bin/server cmd/server/main.go  # Production build
```

### Mock Generation
```bash
go generate ./...           # Generate mocks for interfaces
mockgen -source=file.go -destination=mocks/mock.go
```

## Dependencies Status
- **Current State**: Basic go.mod with minimal dependencies
- **Needed**: All major dependencies must be added to go.mod
- **Version Management**: Use latest stable versions
- **Security**: Regular dependency updates via Dependabot

## Performance Considerations
- **Fiber**: High-performance HTTP framework
- **pgx**: Optimized PostgreSQL driver
- **Connection Pooling**: Database connection management
- **UUIDv7**: Index-friendly primary keys
- **Minimal Dependencies**: Reduced binary size and attack surface

## Development Constraints
- **Go Version**: Requires Go 1.21+ for modern features
- **Database**: PostgreSQL 12+ required
- **Memory**: Minimal memory footprint design
- **Startup Time**: Fast application startup
- **Binary Size**: Optimized for container deployment

## Testability Without Real Data Access
- Interface-driven design enables swapping real data sources with mocks or in-memory fakes during tests.
- SQL mocking validates repository logic without a live database.
- In-memory repositories provide deterministic, fast unit tests; optional test containers may be used for integration tests when desired.
- Dependency injection wires implementations per environment: memory, mock, or postgres.
- Configuration flags (e.g., DATA_SOURCE=memory|mock|postgres) control which implementation is used.
- Test helpers and factories create deterministic test data without relying on external services.

## Graceful Shutdown
- Application lifecycle management coordinates startup and shutdown; resources register start/stop hooks.
- HTTP server supports shutdown with context-based timeouts to drain in-flight requests.
- Database pools are explicitly closed; long-running queries respect context deadlines.
- Background workers listen on context cancellation and stop promptly.
- Shutdown timeout is configurable (default 15s); logs/metrics/traces are flushed before exit.