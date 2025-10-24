# Architecture Document: Go Fiber Microservice Template

## System Overview

This template implements Clean Architecture with Domain-Driven Design principles, creating a modular monolith that can easily transition to microservices. The architecture emphasizes clear separation of concerns, testability, and maintainability.

## Architectural Principles

### Core Principles
1. **Dependency Inversion:** High-level modules don't depend on low-level modules; both depend on abstractions
2. **Single Responsibility:** Each component has one reason to change
3. **Open/Closed:** Components are open for extension, closed for modification
4. **Interface Segregation:** Clients shouldn't depend on interfaces they don't use
5. **Domain Centric:** Business logic lives in the domain layer, free of infrastructure concerns

### Additional Principles
- **SQL-First Development:** Database schema and queries drive data access layer design
- **Container-Native Design:** Architecture optimized for container deployment and scaling
- **Microservice-Ready:** Clear service boundaries with well-defined API contracts
- **Test-Driven Infrastructure:** Comprehensive testing capabilities across all layers

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   HTTP Handlers │  │   Middleware    │  │   Swagger   │ │
│  │   (Fiber)       │  │   (Auth, CORS)  │  │   Docs      │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Usecases      │  │   DTOs          │  │   Validators │ │
│  │   (Business     │  │   (Request/     │  │   (Input     │ │
│  │   Logic)        │  │   Response)     │  │   Validation)│ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Entities      │  │   Value Objects │  │   Domain    │ │
│  │   (Core Models) │  │   (Business     │  │   Interfaces│ │
│  │                 │  │   Rules)        │  │   (Contracts)│ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Repositories  │  │   External      │  │   Database  │ │
│  │   (Data Access) │  │   Services      │  │   (SQLc)    │ │
│  │   (SQLc +       │  │   (HTTP Clients)│  │   (Generated│ │
│  │   Migrations)   │  │                 │  │   Code)     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure Deep Dive

### Application Entry Points
```
cmd/service/
├── main.go                    # Service bootstrap and application lifecycle
├── config.go                  # Configuration loading and validation
├── server.go                  # Fiber server setup and middleware
└── app.go                     # Dependency injection container setup
```

### Core Application Logic
```
internal/
├── handlers/                  # HTTP Presentation Layer
│   ├── user/                  # Domain-specific handlers
│   │   ├── user_handler.go    # User-related HTTP endpoints
│   │   ├── user_handler_test.go
│   │   └── user_mock.go       # Generated mock (go:generate)
│   ├── middleware/            # Cross-cutting HTTP concerns
│   │   ├── auth.go            # JWT authentication middleware
│   │   ├── cors.go            # CORS handling
│   │   ├── logging.go         # Request/response logging
│   │   └── recovery.go        # Panic recovery
│   └── common/                # Shared handler utilities
│       ├── response.go        # Standardized response formatting
│       └── validation.go      # Input validation helpers
├── usecases/                  # Application Business Logic Layer
│   ├── user/                  # Domain-specific business logic
│   │   ├── user_usecase.go    # User business operations
│   │   ├── user_usecase_test.go
│   │   └── user_mock.go       # Generated mock
│   └── common/                # Shared usecase utilities
│       ├── transaction.go     # Transaction management
│       └── validation.go      # Business validation rules
├── repositories/              # Data Access Layer
│   ├── user/                  # Domain-specific data access
│   │   ├── user_repository.go # User database operations
│   │   ├── user_repository_test.go
│   │   └── user_mock.go       # Generated mock
│   └── common/                # Shared repository utilities
│       ├── transaction.go     # Transaction handling via sqlc
│       └── base_repository.go # Common repository patterns
├── domains/                   # Domain Interface Definitions
│   ├── user/                  # Domain-specific contracts
│   │   ├── user.go            # User domain interfaces
│   │   └── user_entities.go   # User domain entities
│   └── common/                # Shared domain concepts
│       ├── base.go            # Base interfaces and types
│       └── errors.go          # Domain-specific error types
├── models/                    # Domain Entities and Value Objects
│   ├── user.go                # User entity definition
│   ├── user_validation.go     # User business rules
│   └── common/                # Shared models
│       ├── base.go            # Base model behaviors
│       └── pagination.go      # Pagination patterns
└── infrastructure/            # External Infrastructure
    ├── database/              # Database setup and management
    │   ├── connection.go      # Database connection management
    │   ├── migrations.go      # Migration runner
    │   └── health.go          # Database health checks
    ├── sqlc/                  # Generated SQL code
    │   ├── user.sql.go        # Generated user queries
    │   └── models.go          # Generated data models
    ├── config/                # Configuration Management
    │   ├── config.go          # Configuration structs
    │   ├── env.go             # Environment variable loading
    │   └── validation.go      # Config validation
    └── errors/                # Error Handling Infrastructure
        ├── errors.go          # Error type definitions
        ├── http_errors.go     # HTTP error mapping
        └── error_handlers.go  # Error handling middleware
```

### Shared Utilities
```
pkg/                           # Public API for external consumption
├── response/                  # Standardized response formats
│   ├── response.go            # Response struct definitions
│   ├── error_response.go      # Error response formatting
│   └── pagination.go          # Paginated response helpers
├── validation/                # Input validation utilities
│   ├── validator.go           # Validation framework setup
│   └── rules.go               # Common validation rules
├── logger/                    # Structured logging utilities
│   ├── logger.go              # Logger configuration
│   ├── context.go             # Context-aware logging
│   └── middleware.go          # HTTP logging middleware
└── utils/                     # General utility functions
    ├── crypto.go              # Cryptographic helpers
    ├── time.go                # Time utilities
    └── strings.go             # String manipulation helpers
```

## Database Architecture

### SQL-First Approach
1. **Schema Definition:** Database schema defined in `/sql/migrations`
2. **Query Development:** SQL queries written in `/sql/queries` with sqlc annotations
3. **Code Generation:** `sqlc generate` creates type-safe Go code
4. **Type Safety:** Generated code provides compile-time query validation

### Migration Strategy
```
sql/
├── migrations/                # Database version control
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   └── schema.sql             # Current schema reference
└── queries/                   # SQL queries for sqlc
    ├── user.sql               # User-related queries
    └── common.sql             # Common query patterns
```

### Transaction Management
- Repository layer handles transaction boundaries
- Use `samber/do` for dependency injection of transaction contexts
- Support for distributed transactions across multiple domains

## Dependency Injection Architecture

### Container Structure
```
app/
├── container.go               # DI container definition
├── providers/                 # Dependency providers
│   ├── database_provider.go   # Database connection provider
│   ├── config_provider.go     # Configuration provider
│   ├── logger_provider.go     # Logger provider
│   └── repository_provider.go # Repository layer providers
└── wire/                      # Wire-generated injection code
    └── wire_gen.go
```

### Provider Patterns
- **Singleton:** Database connections, configuration, logger
- **Scoped:** HTTP handlers, use cases per request
- **Transient:** Value objects, DTOs

## Security Architecture

### Authentication & Authorization
```
security/
├── jwt/                       # JWT implementation
│   ├── token.go               # Token generation/validation
│   ├── claims.go              # Custom claim definitions
│   └── middleware.go          # JWT middleware
├── auth/                      # Authentication interfaces
│   ├── auth_service.go        # Authentication business logic
│   └── auth_provider.go       # External auth integration
└── rbac/                      # Role-based access control
    ├── permissions.go         # Permission definitions
    └── middleware.go          # RBAC middleware
```

### Security Layers
1. **Network Layer:** TLS/HTTPS, rate limiting
2. **Application Layer:** Input validation, authentication
3. **Domain Layer:** Business rule authorization
4. **Data Layer:** Row-level security, encryption

## Monitoring & Observability

### Health Check Architecture
```
health/
├── checker.go                 # Health check interface
├── database_health.go         # Database health check
├── external_service_health.go # External service health checks
└── middleware.go              # Health check HTTP endpoints
```

### Observability Stack
- **Logging:** Structured logging with correlation IDs
- **Metrics:** Prometheus-compatible metrics
- **Tracing:** OpenTelemetry integration hooks
- **Health Checks:** Liveness, readiness, startup probes

## Configuration Architecture

### Environment-Based Configuration
```
config/
├── types.go                   # Configuration type definitions
├── loader.go                  # Environment variable loading
├── validation.go              # Configuration validation
└── defaults.go                # Default values
```

### Configuration Hierarchy
1. **Environment Variables:** Runtime configuration
2. **Configuration Files:** Environment-specific configs
3. **Default Values:** Built-in sensible defaults

## Error Handling Architecture

### Structured Error Types
```
errors/
├── domain_errors.go           # Domain-specific error types
├── infrastructure_errors.go   # Infrastructure error types
├── http_errors.go             # HTTP error mapping
└── error_handlers.go          # Global error handlers
```

### Error Response Format
```json
{
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User not found",
    "details": {
      "user_id": "123",
      "timestamp": "2024-01-01T12:00:00Z"
    }
  }
}
```

## Testing Architecture

### Test Organization
```
tests/
├── integration/               # Integration tests
├── e2e/                       # End-to-end tests
├── fixtures/                  # Test data fixtures
├── mocks/                     # Generated mocks
└── testutils/                 # Testing utilities
```

### Testing Strategy
1. **Unit Tests:** Fast, isolated tests for business logic
2. **Integration Tests:** Database and external service integration
3. **Contract Tests:** API contract validation
4. **Performance Tests:** Load and stress testing

## Deployment Architecture

### Container Strategy
- **Multi-stage builds:** Optimized for production
- **Health checks:** Built-in container health checks
- **Configuration:** Environment-based runtime configuration
- **Security:** Non-root user, minimal attack surface

### Scaling Considerations
- **Horizontal Scaling:** Stateless service design
- **Database Scaling:** Connection pooling, read replicas
- **Caching:** Application-level caching strategies
- **Load Balancing:** Multiple instance support

## Architecture Decision Records (ADRs)

### ADR-001: Fiber Framework Selection
**Status:** Accepted
**Decision:** Use Fiber v2 as the HTTP framework
**Rationale:** 
- Express.js-like API reduces learning curve
- Superior performance characteristics
- Excellent middleware ecosystem
- Active community and maintenance

### ADR-002: SQL-First Database Access
**Status:** Accepted
**Decision:** Use sqlc with SQL-first approach
**Rationale:**
- Type safety at compile time
- Clear separation of data access logic
- Excellent database optimization visibility
- Migration-friendly approach

### ADR-003: Clean Architecture Implementation
**Status:** Accepted
**Decision:** Implement Clean Architecture with DDD principles
**Rationale:**
- Clear separation of concerns
- Excellent testability
- Framework-agnostic business logic
- Easy refactoring and maintenance

### ADR-004: Dependency Injection with samber/do
**Status:** Accepted
**Decision:** Use samber/do for dependency injection
**Rationale:**
- Type-safe dependency resolution
- Excellent performance characteristics
- Simple API without code generation
- Go-idiomatic approach

This architecture document provides the foundation for understanding the system structure, component relationships, and the reasoning behind key architectural decisions. It serves as a guide for developers working with the template and for future architectural evolution.