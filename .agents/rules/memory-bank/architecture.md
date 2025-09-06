# Architecture Overview

## System Architecture
**Domain-Driven Clean Architecture** with mono-repo structure. Each domain contains a complete Clean Architecture implementation following SOLID principles with clear separation of concerns and domain isolation.

## Project Structure

```
gofiber-skeleton/
├── cmd/                          # Application entry points
│   ├── server/                   # Main server application
│   └── migrate/                  # Database migration runner
├── internal/                     # Private application code
│   ├── domains/                  # Business domains (core business logic)
│   │   ├── auth/                 # Authentication domain
│   │   │   ├── entities/         # Core business entities
│   │   │   ├── usecases/         # Business logic and use cases
│   │   │   ├── repositories/     # Data access interfaces
│   │   │   ├── handlers/         # HTTP handlers (delivery layer)
│   │   │   ├── routes/           # Route definitions
│   │   │   ├── models/           # DTOs and request/response models
│   │   │   └── tests/            # Domain-specific tests
│   │   └── posts/                # Posts domain (example)
│   │       └── [same structure as auth]
│   ├── infrastructure/           # External concerns
│   │   ├── database/             # Database connections and configs
│   │   ├── middleware/           # HTTP middleware
│   │   └── config/               # Configuration management
│   └── shared/                   # Shared components across domains
│       ├── types/                # Common types and interfaces
│       └── container/            # Dependency injection container
├── pkg/                          # Public utilities (can be imported)
│   └── utils/                    # Utility functions
├── migrations/                   # Database migrations
├── docs/                         # API documentation (Swagger)
└── compose.yml            # Development environment
```

## Clean Architecture Layers

### 1. Entities Layer (`entities/`)
- **Purpose**: Core business entities and domain rules
- **Dependencies**: None (innermost layer)
- **Contains**: Business objects, domain interfaces, business rules
- **Example**: `User` entity, `Post` entity with business validation

### 2. Use Cases Layer (`usecases/`)
- **Purpose**: Application business logic and orchestration
- **Dependencies**: Entities only
- **Contains**: Use case interfaces and implementations, application services
- **Example**: `AuthenticateUser`, `CreatePost`, `GetUserPosts`

### 3. Interface Adapters (`repositories/`, `handlers/`, `models/`)
- **Purpose**: Convert data between use cases and external layers
- **Dependencies**: Use cases and entities
- **Contains**: Repository implementations, HTTP handlers, DTOs
- **Example**: `PostgreSQLUserRepository`, `AuthHandler`, `LoginRequest`

### 4. Frameworks & Drivers (`infrastructure/`)
- **Purpose**: External frameworks, databases, web frameworks
- **Dependencies**: All inner layers
- **Contains**: Database drivers, HTTP framework setup, external services
- **Example**: Fiber app setup, PostgreSQL connection, JWT middleware

## Key Architectural Patterns

### Dependency Injection
- **Framework**: Uber FX for dependency injection
- **Pattern**: Constructor injection with interfaces
- **Benefits**: Testability, loose coupling, easy mocking

### Repository Pattern
- **Interface**: Defined in entities/use cases layer
- **Implementation**: In infrastructure layer
- **Benefits**: Database abstraction, testability, swappable storage

### Clean API Design
- **Request/Response**: Separate DTOs in `models/` directory
- **Validation**: Input validation at handler level
- **Error Handling**: Consistent error responses using JSend format

### Domain Isolation
- **Principle**: Each domain is self-contained
- **Communication**: Domains communicate through shared interfaces
- **Benefits**: Independent development, clear boundaries, scalability

## Technology Stack Integration

### HTTP Layer (Go Fiber)
- **Entry Point**: `cmd/server/main.go`
- **Routing**: Domain-specific route groups
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)

### Database Layer
- **ORM**: SQL generation with sqlc
- **Migrations**: golang-migrate for version control
- **Connection**: Configurable database connections with pooling

### Configuration Management
- **Tool**: Viper for multi-source configuration  
- **Loading Hierarchy**: ./config/<env>.yaml → .env → runtime environment variables
- **Pattern**: Environment-specific YAML configs with .env overrides and runtime precedence
- **Structure**: Centralized config with domain-specific sections

### Testing Strategy
- **Unit Tests**: Each layer tested in isolation
- **Integration Tests**: Database integration tests
- **Mocking**: Uber mock for interface mocking
- **Test Data**: SQL mock for database testing

## Critical Code Flows

### Request Flow
1. **HTTP Request** → Fiber Router
2. **Router** → Domain Route Handler
3. **Handler** → Use Case (via dependency injection)
4. **Use Case** → Repository Interface
5. **Repository** → Database/External Service
6. **Response** flows back through layers

### Domain Addition Flow
1. Create domain directory structure
2. Define entities with business rules
3. Create repository interfaces
4. Implement use cases
5. Build HTTP handlers and routes
6. Register dependencies in DI container
7. Add tests for all layers

### Authentication Flow
1. Login request → Auth handler
2. Handler → AuthenticateUser use case
3. Use case → User repository for validation
4. Generate JWT token
5. Return authenticated response
6. Subsequent requests use JWT middleware

## Component Relationships

### Core Dependencies
- **Handlers** depend on **Use Cases**
- **Use Cases** depend on **Repository Interfaces** and **Entities**
- **Repository Implementations** depend on **Database Infrastructure**
- **DI Container** wires all dependencies

### Cross-Domain Communication
- **Shared Types**: Common interfaces and types in `shared/`
- **Event System**: For eventual domain event handling
- **Service Interfaces**: For cross-domain service calls

## Operational Requirements

### Testability Without Real Data Access
- Repository interfaces are the primary seam for substituting data access during tests; domain logic never depends on concrete databases.
- Provide in-memory and mock implementations of repositories for unit tests; use interface-based DI to inject them.
- Use SQL mocking for repository-level tests to validate SQL behavior without a live database.
- Configuration-driven wiring enables switching repository implementations (e.g., memory, mock, postgres) through the DI container.
- Contract tests ensure that mock/fake implementations conform to the same behaviors as real repositories.
- Test data is created via factories/builders; unit tests must not depend on migrations or a running database.
- External integrations are abstracted behind interfaces and can be replaced with no-op or deterministic stubs in tests.

### Graceful Shutdown
- The service must handle termination signals and begin a coordinated shutdown: stop accepting new requests, drain in-flight work, then release resources.
- HTTP server shutdown uses timeouts to allow in-flight requests to complete within a configurable window.
- DI lifecycle hooks ensure resources are closed in order: stop HTTP listener, cancel background workers, close database pools, flush logs/telemetry.
- Background workers and goroutines must observe context cancellation and exit promptly.
- Health probes: readiness becomes false at shutdown initiation; liveness remains true until final termination to support zero-downtime rollouts.
- Default shutdown timeout is 15s and is configurable per environment.