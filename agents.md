# Project Agent Guidelines

This document outlines the key aspects of the E-commerce Backend Boilerplate project, serving as a guide for agents working on the codebase.

## 1. Project Overview

This project provides a complete backend boilerplate for building a production-ready e-commerce management system using Go Fiber and Clean Architecture.

**Purpose:**
To offer a clean, maintainable, and scalable codebase implementing core e-commerce features. It reduces setup time for standardized backend architecture, ensures consistency, and provides built-in support for database migrations, testing, and containerization.

**Key Features:**
- RESTful endpoints for products, orders, and user authentication.
- Implementation of Clean Architecture layers for separation of concerns.
- Uses SQLC for type-safe database queries and golang-migrate for migrations.
- Integrates JWT-based authentication and middleware for secure access control.

**Project Structure (Clean Architecture):**
- **`cmd/server`**: Application entry point.
- **`internal/infrastructure`**: Holds infrastructure components (database, config, middleware, SQLC generated code).
- **`internal/domain`**: Defines interfaces for communication between different layers and core domain entities.
- **`internal/<domain>/handler`**: Manages HTTP requests and responses for specific domains (e.g., product, order, user).
- **`internal/<domain>/repository`**: Handles database interactions for specific domains.
- **`internal/<domain>/usecase`**: Contains business logic and orchestrates calls to repositories for specific domains.
- **`pkg`**: Stores shared packages (currently not explicitly used as per memory bank).
- **`compose.yml`**: Manages Docker services.
- **`Dockerfile`**: Builds the Docker image.
- **`db/migrations`**: Manages database migrations.
- **`db/queries`**: SQL query files used by SQLC.

## 2. Build and Test Commands

**Development Setup:**
1. Install Go (>=1.25) and set GOPATH.
2. Install Docker and Docker Compose.
3. Install `migrate` CLI.
4. Install `sqlc`: `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`.
5. Install `air`: `go install github.com/cosmtrek/air@latest`.
6. Install `golangci-lint`: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`.

**Common Commands:**
- Generate SQLC code: `sqlc generate`
- Run database migrations: `migrate -path db/migrations -database "$DATABASE_URL" up`
- Build the server executable: `go build -o bin/server cmd/server/main.go`
- Run services with Docker Compose: `docker compose up --build`
- Verify works (linting, testing): `go generate ./... && golangci-lint run --fix ./... && go clean -testcache && go test -v -race ./...`

## 3. Code Style Guidelines

The project adheres to Clean Architecture and SOLID Principles, emphasizing:
- **Modularity and Separation of Concerns**: Business logic is decoupled from framework-specific code.
- **Dependency Inversion**: Handlers and use cases depend on interfaces, not concrete implementations.
- **Composition over Inheritance**: Favor small, purpose-specific abstractions.
- **Error Handling**: Explicit checking and handling of errors/exceptions with context.
- **State Management**: Avoid global mutable state; use dependency injection for managing dependencies.
- **Concurrency**: Safe use of concurrency primitives; guard shared state with appropriate synchronization.
- **Go Idioms**: Adherence to standard Go formatting and best practices (enforced by `golangci-lint`).
- **Naming Conventions**: Consistent naming for packages, functions, and variables.

## 4. Testing Instructions

**General Principles:**
- All tests must use mocks for database access and other external services. No tests should interact with a live database.
- Use `go generate` to create mocks for interfaces (`uber-go/mock`).

**Test Types and Locations:**
- **Unit Tests**:
    - Located in the same package as the code they are testing.
    - File names follow the `_test.go` suffix (e.g., `product_handler_test.go`).
    - Focus on testing individual units of code in isolation.
- **Integration Tests**:
    - Located in a top-level `tests/` directory.
    - Mirror the structure of the `internal/` directory (e.g., `tests/integration/product_integration_test.go`).
    - Test the interaction between multiple components.

## 5. Security Considerations

- **Input Validation and Sanitization**: Rigorous validation on all inputs, especially from external sources, using `go-playground/validator`.
- **Parameterized SQL**: SQLC is used to generate type-safe queries, preventing SQL injection vulnerabilities.
- **Authentication**: JWT-based authentication for secure access control.
- **Middleware**: Custom middleware for JWT validation on authenticated API endpoints.
- **Secure Defaults**: Use secure defaults for authentication tokens and configuration settings.
