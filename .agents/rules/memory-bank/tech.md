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

## 4. Database & ORM/Query Builders

-   **PostgreSQL:** The relational database of choice for its robustness, features, and widespread adoption.
-   **`golang-migrate/migrate`:** A database migration tool for Go, enabling version-controlled schema evolution.
-   **`sqlc`:** Generates fully type-safe, idiomatic Go code from raw SQL queries, catching errors at compile-time and simplifying data access.

## 5. Authentication & Authorization

-   **`golang-jwt`:** A Go package for creating and verifying JSON Web Tokens (JWTs), used for stateless authentication.
-   **Bcrypt:** For secure hashing of user passwords.

## 6. Caching

-   **Valkey (or Redis):** An in-memory data structure store, used for caching and potentially other data-intensive operations. Integrated via `compose.yml`.

## 7. API Documentation

-   **`swaggo/swag`:** Automatically generates Swagger/OpenAPI 2.0 documentation from Go source code comments, providing an interactive API explorer.

## 8. Development & Tooling

-   **Air:** A live-reloading command-line tool for Go applications, significantly improving the development feedback loop by automatically rebuilding and restarting the server on file changes.
-   **Docker & Docker Compose:** For containerizing the application and its dependencies (PostgreSQL, Valkey), providing a consistent and isolated development environment.
-   **Makefile:** Contains helper commands for common development tasks (e.g., `make run`, `make test`, `make migrate`).
-   **`go generate`:** Used with `mockgen` and `sqlc` for code generation.

## 9. Testing & Mocking

-   **`go test`:** Go's built-in testing framework.
-   **`go.uber.org/mock/mockgen`:** A tool for generating mock implementations of Go interfaces, crucial for unit testing usecases logic in isolation. Used with `//go:generate` directives.
-   **`DATA-DOG/go-sqlmock`:** A library for mocking the SQL database driver, allowing for robust testing of data access logic without requiring a real database connection.

## 10. Observability

-   **Structured Logging:** Integration with a logging library (e.g., `zap` or `logrus`) for structured, context-rich logging.
-   **Distributed Tracing:** Planned integration with OpenTelemetry or similar for end-to-end request tracing.
-   **Metrics:** Planned integration with Prometheus/Grafana for application performance monitoring.

## 11. Code Quality & Linting

-   **`go fmt`:** Go's official code formatter.
-   **`golangci-lint`:** A fast Go linters aggregator, used to enforce code style and identify potential issues.

## 12. Project Structure

-   **Mono-repo:** A single repository containing multiple Go modules/services.
-   **Clean Architecture / Domain-Driven Design:** Logical separation of concerns into layers and domains (e.g., `cmd/`, `internal/app/`, `internal/domains/`, `internal/infrastructure/`, `internal/shared/`).