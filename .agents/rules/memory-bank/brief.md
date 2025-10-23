# Project Brief: Go Fiber Microservice Template

## Project Identity
**Project Name:** template-go-fiber
**Type:** Go Fiber Microservice Template
**Architecture:** Layered Architecture
**Purpose:** A lightweight and performant foundation for building simple microservices.

## Executive Summary
A production-ready Go backend template using the Fiber framework and a classic Layered Architecture. It is designed for simplicity, performance, and ease of use, making it ideal for developing focused microservices with industry-standard tooling and best practices.

## Core Objectives
1.  **High Performance**: Leverage Fiber's speed for low-latency services.
2.  **Simplicity**: Offer a straightforward structure that is easy to understand and extend.
3.  **Production-Ready**: Include essentials like logging, configuration, error handling, and monitoring.
4.  **Developer Experience**: Ensure a minimal learning curve and quick setup.
5.  **Testability**: Structure the code to facilitate comprehensive unit and integration testing.
6.  **Code Quality**: Enforce high standards through automated linting and testing.

## Project Scope

### In Scope
**Core Infrastructure:** Fiber framework, Layered Architecture (Handler, Usecase, Repository), Database integration, JWT Authentication, Swagger Docs, Environment-based Config, Structured Logging, Centralized Error Handling, Health Checks, Graceful Shutdown, [JSend Response Format](https://github.com/omniti-labs/jsend).
**Development Tools:** Docker support, Basic Makefile commands, Linting, Interface mocking, CI/CD pipeline.
**Documentation:** README, API (Swagger), Example usecase, Architecture Decision Records (ADRs).

## Technical Requirements

### Architecture Principles
**Layered Architecture:**
-   **Handlers:** Manages HTTP requests/responses and data transfer objects (DTOs).
-   **Usecases:** Contains the core business logic.
-   **Repository:** Handles all database interactions via sqlc generated code.
**Dependency Injection:** Use `samber/do v2` to manage dependencies between layers.
**Interface-Driven Design:** All interfaces must include `//go:generate` annotations for automatic mock generation.

### Technology Stack
**Core Framework:** Go 1.25+, Fiber v2+.
**Database:** MariaDB 11+ (default for small/medium systems), PostgreSQL 18+ (for large systems), with flexibility for others.
**Database Access:** `sqlc` for type-safe queries, `golang-migrate/migrate` for migrations.
**Authentication:** `golang-jwt/jwt` (primary), `zitadel/oidc` (alternative).
**Configuration:** Environment variables, loaded via `spf13/viper`.
**Logging:** `slog` (structured logging with context propagation).
**Testing:** `uber-go/mock` (interface mocking), `DATA-DOG/go-sqlmock` (database mocking).
**Utilities:** `samber/do v2` (DI), `samber/lo` (synchronous helpers), `samber/ro` (event-driven streams).
**Code Quality:** `golangci-lint v2`, UUIDv7 for database-friendly unique IDs.

### Project Structure
```
template-go-fiber/
├── cmd/api/                  # Application entry point (main.go)
├── internal/                 # Private application code
│   ├── handler/              # HTTP handlers (with Swagger annotations)
│   ├── usecase/              # Business logic usecases
│   ├── repository/           # Data access layer (handles transactions)
│   ├── middleware/           # HTTP middleware (auth, logging, etc.)
│   ├── config/               # Configuration loading and validation
│   └── domain/               # Core domain models/structs and interfaces
├── pkg/                      # Shared utilities (e.g., response, errors)
├── sql/                      # SQL source files
│   ├── queries/              # SQL queries for sqlc code generation
│   └── migrations/           # Database migration files
├── docs/                     # Swagger and other documentation
├── .agents/                  # Agent rules and memory bank
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Makefile                  # Common development tasks
├── Dockerfile                # Container image definition
├── .env.example              # Environment variable template
└── README.md                 # Project documentation
```

## Key Features
1.  **RESTful API Foundation**: Basic routing, request validation, and [JSend-formatted JSON responses](https://github.com/omniti-labs/jsend).
2.  **Authentication**: JWT-based middleware for securing endpoints.
3.  **Database Integration**: `sqlc` for type-safe query generation and `golang-migrate` for schema versioning.
4.  **Configuration Management**: Simple, environment-variable-driven configuration with validation.
5.  **API Documentation**: Automated Swagger/OpenAPI generation from handler annotations.
6.  **Containerization**: Ready-to-use Dockerfile for building and deploying the service.
7.  **Monitoring**: Application metrics and health check endpoints.
8.  **CI/CD Pipeline**: Automated testing, linting, security scanning, and deployment.
9.  **Performance Optimization**: Follows patterns from goperf.dev guidelines.
10. **Security**: Rate limiting, input validation/sanitization, secure headers, HTTPS support.

## Constraints & Requirements

### Code Quality
1.  Must pass `golangci-lint v2` with zero errors.
2.  All interfaces must include `//go:generate` annotations for mock generation.
3.  Maintain consistent code style and documentation standards.

### Testing
1.  All usecases must have unit tests with high coverage.
2.  All repositories must have tests using `go-sqlmock`.
3.  Use `uber-go/mock` for interface mocking in tests.
4.  Integration tests for critical API endpoints.

### Security
1.  Implement rate limiting on public endpoints.
2.  Add comprehensive input validation and sanitization.
3.  Use secure headers and HTTPS for external communications.
4.  Follow OWASP best practices for API security.

### Performance
1.  Follow optimization patterns from https://goperf.dev/01-common-patterns/
2.  Follow networking patterns from https://goperf.dev/02-networking/

### Deployment
1.  All services must support container environments.
2.  Implement graceful shutdown mechanisms.
3.  Services should be stateless to support horizontal scaling.

### CI/CD Pipeline
1.  Automated testing on every commit.
2.  Code quality checks with `golangci-lint v2`.
3.  Security vulnerability scanning.
4.  Container image building and registry push.

### Documentation
1.  Maintain comprehensive API documentation via Swagger.
2.  Include Architecture Decision Records (ADRs) for key design choices.
3.  Provide clear setup and contribution guidelines.