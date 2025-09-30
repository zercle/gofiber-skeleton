# Technology Stack: Go Fiber Template Repository

This document outlines the core technologies and tools utilized in the Go Fiber Template Repository.

## 1. Backend Framework & Language

-   **Go (Golang):** The primary programming language for its performance, concurrency features, and strong typing.
-   **Fiber v2:** A high-performance, Express.js-inspired web framework for Go, used for building RESTful APIs.

## 2. Dependency Management & Injection

-   **Go Modules:** Standard Go dependency management system.
-   **Uber's fx:** A dependency injection framework for Go, built on the concept of "application as a function graph," promoting modularity and testability.

## 3. Configuration Management

-   **Viper:** A complete Go configuration solution that supports various formats (JSON, TOML, YAML, HCL, INI, envfile) and environment variables, with clear precedence rules.
-   **`dotenv`:** For loading environment variables from `.env` files during local development.

## 4. Database & Query Builders

-   **PostgreSQL 18:** The relational database for production data storage with full ACID compliance.
-   **`golang-migrate/migrate`:** A database migration tool for Go, enabling version-controlled schema evolution.
-   **`sqlc`:** ✅ **ACTIVE** - Generates fully type-safe, idiomatic Go code from raw SQL queries, catching errors at compile-time and simplifying data access. Replaces ORM for better performance and type safety.

## 5. Authentication & Authorization

-   **`golang-jwt`:** A Go package for creating and verifying JSON Web Tokens (JWTs), used for stateless authentication.
-   **Bcrypt:** For secure hashing of user passwords.

## 6. Caching

-   **Valkey 8:** ✅ **CONFIGURED** - Redis-compatible in-memory data structure store for caching, rate limiting storage, and session management. Integrated via `docker-compose.yml` with persistent data volume.

## 7. API Documentation

-   **`swaggo/swag`:** ✅ **ACTIVE** - Automatically generates Swagger/OpenAPI 2.0 documentation from Go source code comments. All endpoints documented with request/response schemas. Access at `/swagger/*`.

## 8. Development & Tooling

-   **Air:** A live-reloading command-line tool for Go applications, significantly improving the development feedback loop by automatically rebuilding and restarting the server on file changes.
-   **Docker & Docker Compose:** For containerizing the application and its dependencies (PostgreSQL, Valkey), providing a consistent and isolated development environment.
-   **Makefile:** Contains helper commands for common development tasks (e.g., `make run`, `make test`, `make migrate`).
-   **`go generate`:** Used with `mockgen` and `sqlc` for code generation.

## 9. Testing & Mocking

-   **`go test`:** Go's built-in testing framework.
-   **`go.uber.org/mock/mockgen`:** ✅ **CONFIGURED** - Generates mock implementations of Go interfaces via `//go:generate` directives for isolated unit testing.
-   **`DATA-DOG/go-sqlmock`:** Library for mocking SQL database driver for repository testing without real database.
-
## 10. Request Validation

-   **`go-playground/validator/v10`:** ✅ **ACTIVE** - Comprehensive struct and field validation with tags (required, email, min, max, etc.). Integrated in all handlers with JSend error responses.

## 11. Observability & Monitoring

-   **`rs/zerolog`:** ✅ **ACTIVE** - High-performance structured logging with request context (request_id, method, path, status, duration, IP, user_agent).
-   **Request ID Middleware:** ✅ **ACTIVE** - UUID-based request tracing across distributed systems via X-Request-ID header.
-   **Distributed Tracing:** Planned integration with OpenTelemetry for end-to-end request tracing.
-   **Metrics:** Planned integration with Prometheus/Grafana for application performance monitoring.

## 12. Code Quality & Linting

-   **`go fmt`:** Go's official code formatter.
-   **`golangci-lint`:** ✅ **ACTIVE** - Fast Go linters aggregator enforcing code style. Integrated in CI/CD pipeline.

## 13. Production Middleware

-   **Recovery Middleware:** ✅ **ACTIVE** - Gracefully handles panics and prevents server crashes.
-   **Rate Limiting:** ✅ **ACTIVE** - 100 req/min for API endpoints, 5 req/min for auth (brute force protection).
-   **Graceful Shutdown:** ✅ **ACTIVE** - SIGTERM/SIGINT handling with 30s timeout for zero-downtime deployments.

## 14. Project Structure

-   **Mono-repo:** Single repository with modular feature-based architecture.
-   **Clean Architecture / Domain-Driven Design:** ✅ **IMPLEMENTED** - Strict separation into `cmd/`, `internal/`, `pkg/`, and `db/` with feature-based organization (`internal/user/`, `internal/post/`).

## 15. API Response Standards

-   **JSend Specification:** ✅ **ACTIVE** - All API responses follow JSend format with `status`, `data`, `message`, and `code` fields for consistency and client-side error handling.