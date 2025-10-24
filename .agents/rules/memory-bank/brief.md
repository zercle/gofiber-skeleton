# Project Brief: Go Fiber Microservice Template

## Project Identity
**Project Name:** template-go-fiber
**Type:** Go Fiber microservice Template
**Architecture:** Clean Architecture with Domain Driven Design
**Purpose:** A lightweight and performant foundation for building microservices with clear domain boundaries and container-native deployment.

## Executive Summary
A production-ready Go microservice template using the Fiber framework and Clean Architecture. It is designed for simplicity, performance, and ease of use, making it ideal for developing focused microservices with clear domain boundaries and container-native deployment patterns.

## Core Objectives
1. **High Performance**: Leverage Fiber's speed and Express.js-like API for low-latency microservices.
2. **Simplicity**: Offer a straightforward structure that is easy to understand and deploy.
3. **Production-Ready**: Include essentials like logging, configuration, and error handling for microservices.
4. **Developer Experience**: Ensure a minimal learning curve and quick setup for service development.
5. **Testability**: Structure the code to facilitate unit and integration testing.
6. **Container-Native**: Designed specifically for container deployment and orchestration.

## Project Scope

### In Scope
**Core Infrastructure:** Fiber framework, Clean Architecture with DDD (Handler, Usecase, Repository, Domain), Database integration (sqlc + golang-migrate), JWT Authentication, Swagger Docs, Environment-based Config, Structured Logging, Centralized Error Handling, Health Checks, Graceful Shutdown, Interface Mocking.
**Microservice Features:** Service discovery ready, Circuit breaker patterns, Distributed tracing hooks, Request/Response transformation, External service integration patterns.
**Development Tools:** Docker support, Makefile commands, Linting (golangci-lint v2), sqlc code generation, Database migrations.
**Documentation:** README, API (Swagger), Architecture Decision Records (ADRs), Example domain implementation.

## Technical Requirements

### Architecture Principles
**Clean Architecture with Domain-Driven Design:**
- **Handler Layer:** Manages HTTP requests/responses, DTOs, and Swagger annotations.
- **Usecase Layer:** Contains business logic and orchestrates domain operations.
- **Repository Layer:** Handles database interactions and transactions via sqlc generated code.
- **Domain Layer:** Defines interfaces, entities, value objects, and domain contracts.
- **Mock Generation:** All interfaces must include `//go:generate` annotations for mockgen.
**Microservice Architecture:** Focused service boundaries with clear API contracts.
**Dependency Injection:** Use `samber/do v2` for DI container management across all layers.

### Technology Stack
**Core Framework:** Go 1.25+, Fiber v2+.
**Database Options:**
- MariaDB 11+ (small-medium systems, default)
- PostgreSQL 18+ (large systems)
- FerretDB 2.5+ (document-based storage)
- Valkey 9+ (in-memory key-value, Redis replacement)
**Database Tools:** `sqlc` for type-safe queries, `golang-migrate/migrate` for migrations.
**Authentication:** `golang-jwt/jwt` (primary), `zitadel/oidc` (alternative).
**Configuration:** Environment variables with validation via `spf13/viper`.
**Logging:** `slog` (structured logging with context propagation).
**Utilities:**
- `samber/do v2` (Dependency Injection)
- `samber/lo` (Synchronous helpers for finite sequences)
- `samber/ro` (Event-driven infinite data streams)
**Testing:** `uber-go/mock` (interface mocking), `DATA-DOG/go-sqlmock` (DB mocking).
**Code Quality:** `golangci-lint v2`, UUIDv7 for database-friendly unique IDs.

### Project Structure
```
template-go-fiber/
├── cmd/service/              # Microservice entry point (main.go)
├── internal/                 # Private application code
│   ├── handlers/             # HTTP handlers with Swagger annotations
│   ├── usecases/             # Business logic layer
│   ├── repositories/         # Data access via sqlc (handles transactions)
│   ├── domains/              # Domain interfaces and contracts
│   ├── models/               # Domain models and entities
│   ├── middleware/           # HTTP middleware (auth, cors, etc.)
│   ├── infrastructure/       # Project infrastructure
│   │   └── sqlc/             # Generated SQL code from sqlc
│   ├── config/               # Configuration loading and validation
│   └── errors/               # Custom error types and handling
├── pkg/                      # Shared utilities (response, validation, etc.)
├── sql/                      # SQL source files
│   ├── queries/              # SQL queries for sqlc code generation
│   └── migrations/           # Database migration files (golang-migrate)
├── docs/                     # Swagger/OpenAPI specs and ADRs
├── .agents/                  # Agent rules and memory bank
│   └── rules/
│       └── memory-bank/      # Project context and documentation
├── bin/                      # Binary releases
├── tmp/                      # Temporary files and artifacts
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Makefile                  # Development tasks (build, test, lint, migrate)
├── compose.yml               # Local development environment
├── Dockerfile                # Container definition
├── .env.example              # Environment variable template
└── README.md                 # Project documentation
```

## Key Features
1. **RESTful API Foundation**: Fiber-based routing, request validation, and JSON responses.
2. **Authentication**: JWT-based middleware for securing endpoints (with OIDC alternative).
3. **Type-Safe Database Access**: `sqlc` for generating type-safe Go code from SQL queries.
4. **Database Migrations**: `golang-migrate/migrate` for version-controlled schema changes.
5. **Configuration Management**: Environment-variable-driven with validation and defaults.
6. **API Documentation**: Automated Swagger/OpenAPI generation from handler annotations.
7. **Testing Infrastructure**: Mock generation via `go:generate`, sqlmock for DB tests.
8. **Containerization**: Production-ready Dockerfile optimized for microservice deployment.
9. **Performance Optimization**: Following goperf.dev patterns and Fiber's performance characteristics.
10. **Security**: Rate limiting, input validation/sanitization, secure headers, HTTPS support.
11. **Service Readiness**: Health check endpoints, graceful shutdown, signal handling.
12. **Distributed Systems Support**: Hooks for service discovery, circuit breakers, and tracing.

## Constraints & Requirements
1. **Code Quality:** Must pass `golangci-lint v2` with all checks enabled.
2. **Testing:** High test coverage across all layers (handlers, usecases, repositories).
3. **Mock Generation:** Every interface file must have `//go:generate` annotations.
4. **Swagger Documentation:** Every HTTP handler must include Swagger annotations.
5. **Security:** Rate limiting, input validation, sanitization, secure headers, HTTPS.
6. **Error Handling:** Structured error types with proper HTTP status mapping.
7. **Logging:** Structured logging with request correlation and context propagation.
8. **Monitoring:** Health checks, metrics endpoints, and distributed tracing support.
9. **Statelessness:** Service must be stateless for horizontal scaling.
10. **Container Support:** Optimized for container deployment with multi-stage builds.
11. **Graceful Shutdown:** Implement graceful shutdown with proper resource cleanup.
12. **Service Boundaries:** Clear API contracts and domain separation.
13. **Configuration:** Environment-based configuration with validation.
14. **Health Checks:** Multiple health endpoints (liveness, readiness, startup).

## Development Workflow
1. **SQL Development**: Write queries in `/sql/queries` with sqlc annotations.
2. **Code Generation**: Run `sqlc generate` to create Go code in `/internal/infrastructure/sqlc`.
3. **Database Migrations**: Use `migrate -path /sql/migrations -database <db_url> up`.
4. **Testing**: Generate mocks with `go generate ./...`, write tests with sqlmock.
5. **Container Development**: Use `docker-compose up` for local development environment.
6. **Service Testing**: Run integration tests in containerized environment.
7. **CI/CD Pipeline**: Automated testing, linting, security scanning, container builds.

## Documentation Requirements
1. **API Documentation**: Maintain Swagger/OpenAPI specs via code annotations.
2. **Architecture Decision Records (ADRs)**: Document significant design decisions.
3. **Service Documentation**: Clear API contracts and service boundaries.
4. **Deployment Guide**: Container deployment and orchestration instructions.
5. **README**: Comprehensive project overview with quick-start guide.

## Quality Assurance
1. Generate mocks and tests with `go generate ./...`
2. Lint with `golangci-lint run --fix ./...`
3. Run tests with `go clean -testcache && go test -v -race ./...`
4. Build container image with `docker build -t service-name .`
5. Test service health endpoints before deployment.