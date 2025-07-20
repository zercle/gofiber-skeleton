# Go URL Shortener Service: Clean Architecture Boilerplate

This document outlines the detailed requirements for building a high-performance URL shortening service in Go, adhering strictly to Clean Architecture and SOLID principles. The final output must be a complete, runnable, and well-documented boilerplate project ready for immediate development and deployment.

## 1. Core Functionality

*   **User Management:**
    *   Endpoint for user registration.
    *   Endpoint for user login (JWT-based authentication).
    *   Authenticated users can create, view, update, and delete their shortened URLs.
*   **URL Shortening:**
    *   Endpoint for authenticated users to create tiny URLs.
    *   Endpoint for guest users to create tiny URLs (without management capabilities post-creation).
    *   Support for custom short codes for authenticated users.
*   **Redirection:**
    *   Secure redirection from short URLs to their original long URLs. Implement safeguards against open redirect vulnerabilities.
*   **QR Code Generation:**
    *   Generate QR codes for all created short URLs, accessible to both registered users and guests.

## 2. Architectural Principles

Strictly adhere to the following principles to ensure a maintainable, flexible, and scalable application:

*   **Clean Architecture:** Ensure clear separation of concerns between layers (entities, use cases, repositories, delivery). The core business logic (use cases) must remain independent of external frameworks (Go Fiber) and database implementations.
*   **SOLID Principles:** Apply SOLID principles throughout the codebase for maintainability, flexibility, and extensibility.
*   **Decoupling:** Ensure strong decoupling between the application's core logic (use cases) and the web framework (Go Fiber) and data access mechanisms. This means:
    *   Use cases should define interfaces for their dependencies (e.g., `UserRepository`, `URLRepository`).
    *   The `delivery` layer (HTTP handlers) will be responsible for registering routes with the Go Fiber application instance, but the actual business logic will be delegated to the `usecases` layer.
    *   The `repository` layer will implement the interfaces defined by the `usecases` layer, abstracting database details.

## 3. Technical Stack & Integration

*   **Go Framework:** Utilize `Go Fiber` for building efficient and high-performance RESTful APIs. Fiber will be used primarily in the `delivery/http` layer to handle HTTP requests and responses, and to register routes.
*   **Database:** PostgreSQL (default assumption).
*   **Database Interactions:**
    *   **Primary Keys (PKs):** All primary keys across the application must use **UUIDv7**.
        *   **UUIDv7 Generation Example (Go):**
            ```go
            package main

            import (
            	"fmt"
            	"log"

            	"github.com/gofrs/uuid" // Recommended library for UUIDv7
            )

            func main() {
            	// Generate a UUIDv7
            	u7, err := uuid.NewV7()
            	if err != nil {
            		log.Fatalf("failed to generate UUIDv7: %v", err)
            	}
            	fmt.Printf("Generated UUIDv7: %s\n", u7.String())

            	// In PostgreSQL, store as a UUID type. Example SQL:
            	// CREATE TABLE users (
            	//     id UUID PRIMARY KEY, -- Application generates UUIDv7
            	//     name TEXT NOT NULL
            	// );
            }
            ```
    *   `sqlc`: Generate type-safe Go code from SQL queries for robust and efficient database interactions. Configure `sqlc.yaml` and include example SQL schema and queries in `db/queries/`.
    *   `golang-migrate`: Integrate for managing database schema migrations and versioning. Define migration files in `db/migrations/`.
*   **Configuration Management:**
    *   `Viper`: Implement for loading application configurations from `configs/app.yaml`.
    *   Allow configuration overrides via environment variables.
*   **API Documentation:**
    *   `Swagger`: Generate comprehensive Swagger (OpenAPI) documentation for all API endpoints using `swaggo/swag`. Place the generated specification files within the `api/` directory.
*   **Development Workflow:**
    *   `Air`: Configure for live reloading during development to enhance developer productivity.
*   **Containerization:**
    *   `Dockerfile`: Provide a multi-stage `Dockerfile` for creating optimized production-ready Docker images. Use `mirror.gcr.io` for all base image names (e.g., `mirror.gcr.io/library/golang:1.22-alpine`, `mirror.gcr.io/library/alpine:latest`).
    *   `compose.yaml`: Include a `compose.yaml` file for orchestrating local development services (e.g., PostgreSQL, Valkey, application service).
*   **Authentication:**
    *   `JWT (JSON Web Tokens)`: Implement JWT-based authentication for secure user access and session management.
*   **Caching/Storage:**
    *   `Valkey (Redis fork)`: Integrate Valkey for efficient storage of URL mappings, rate limiting, or session management.

## 4. API Response Format

All JSON responses from API endpoints must adhere to the `omniti-labs/jsend` specification for consistent and predictable API communication (e.g., `success`, `fail`, `error` statuses with `data` or `message` fields).

## 5. Project Structure

Implement the following Clean Architecture-compliant directory structure:

```
.
├── cmd/                  # Main application entry points
│   └── api/              # Main API service entry (initializes Fiber, injects dependencies)
├── internal/             # Core application logic
│   ├── delivery/         # API handlers, controllers (Fiber specific, registers routes)
│   │   └── http/
│   ├── usecases/         # Business logic, application-specific rules, defines repository interfaces
│   ├── repository/       # Database implementations of repository interfaces
│   ├── entities/         # Core business objects/models
│   └── configs/          # Internal configuration structures/loaders
├── pkg/                  # Reusable utilities, common libraries
├── api/                  # API specifications (e.g., Swagger/OpenAPI YAML)
├── configs/              # External application configuration files (e.g., app.yaml)
├── db/                   # Database schema, migrations, and sqlc queries
│   ├── migrations/
│   └── queries/
├── tests/                # Integration tests
├── mocks/                # Generated mocks for testing
├── .github/              # GitHub Actions workflows
│   └── workflows/
├── Dockerfile            # Docker build instructions
├── compose.yaml          # Docker Compose file for local development
├── go.mod                # Go module file
├── go.sum                # Go module checksums
├── Makefile              # Common development commands (optional but recommended)
└── README.md             # Project documentation
```

## 6. Testing Strategy

*   **Unit Tests:**
    *   Place unit tests (`_test.go` suffix) alongside their respective source files within the `internal/` directories.
    *   Utilize `DATA-DOG/go-sqlmock` for mocking database interactions when testing the `repository` layer, allowing tests to run without a real database connection.
    *   Generate and utilize `uber-go/mock` mocks within the `mocks/` directory for other interfaces (e.g., use case dependencies, external services) to facilitate effective unit testing and isolate dependencies.
*   **Integration Tests:** Place all comprehensive integration tests within the dedicated `tests/` directory. These should cover end-to-end flows, potentially using the `compose.yaml` setup for a full environment.

## 7. CI/CD (GitHub Actions)

Set up GitHub Actions workflows within `.github/workflows/` for:

*   Automated `go test` execution on push/pull request.
*   Automated `golangci-lint` checks for code quality and style.