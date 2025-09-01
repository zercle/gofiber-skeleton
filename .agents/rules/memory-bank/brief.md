# Go Fiber Backend Boilerplate Template

## 1. Introduction

This document serves as a standardized template for initializing new projects using the Go Fiber Backend Boilerplate. The boilerplate is built on **Go Fiber**, following **Clean Architecture** and **SOLID principles** to provide a robust, scalable, and maintainable foundation for backend systems.

**Purpose**: To streamline the setup process for new projects, ensuring consistency and saving development time by providing a pre-configured architecture with essential features.

## 2. Core Principles & Architecture

The boilerplate follows a domain-driven Clean Architecture approach, which decouples business logic from infrastructure concerns.

### 2.1. Layered Architecture

-   **Presentation/API Layer**: (`cmd/server`, `internal/<domain>/handler/router.go`) - Handles HTTP requests, routing, and dependency injection setup.
-   **Handler Layer**: (`internal/<domain>/handler`) - Validates (go-playground/validator) request payloads, calls use cases, and formats HTTP responses (omniti-labs/jsend).
-   **Use Case/Service Layer**: (`internal/<domain>/usecase`) - Encapsulates all business rules and orchestrates repository calls.
-   **Domain Layer**: (`internal/<domain>`) - Contains domain models and interfaces (contracts) for repositories and use cases.
-   **Repository/Infrastructure Layer**: (`internal/<domain>/repository`, `internal/infrastructure/sqlc`) - Implements data persistence using SQLC-generated queries and manages database transactions.
-   **Shared Infrastructure**: (`internal/infrastructure`) - Contains shared components like database connections, configuration management (Viper), and middleware.

### 2.2. Key Design Patterns & Decisions

-   **Dependency Injection**: Uses `samber/do` to manage and inject dependencies, promoting loose coupling.
-   **SQLC for Type-Safe Queries**: Generates Go code from raw SQL queries, ensuring type safety and preventing SQL injection.
-   **Repository-Managed Transactions**: Repositories are responsible for all database transactions, ensuring atomicity.
-   **Viper for Configuration**: Manages configuration from `configs/<env>.yaml`, environment variables and `.env` files.
-   **gofiber/swagger**: Automatically generates OpenAPI documentation.
-   **JWT**: Use JWT for private endpoint authentication.

## 3. Getting Started: Project Initialization

Follow these steps to initialize a new project based on this template.

### 3.1. Prerequisites

Ensure the following tools are installed:
-   Go (>=1.25)
-   Docker and Docker Compose
-   `golang-migrate` CLI
-   `sqlc` CLI
-   `air` (for hot-reloading)
-   `golangci-lint` (for linting)

You can install the Go tools using:
```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 3.2. Initial Setup

1.  **Clone the Boilerplate**:
    ```bash
    git clone <repository-url> <new-project-name>
    cd <new-project-name>
    ```

2.  **Configure Environment**:
    -   Copy the example environment file: `cp .env.example .env`
    -   Update the `.env` file with your project-specific configurations (database credentials, app settings, etc.).

3.  **Initialize Go Module**:
    -   Update the module path in `go.mod` to match your new project's repository URL.
    -   Run `go mod tidy` to sync dependencies.

4.  **Run the Database**:
    -   Start the PostgreSQL container: `docker compose up -d db`

5.  **Apply Database Migrations**:
    -   Run the initial schema migrations: `migrate -path db/migrations -database "$DATABASE_URL" up`
    (Ensure `$DATABASE_URL` is set in your environment or `.env` file).

6.  **Generate SQLC Code**:
    -   Generate Go code from your SQL queries: `sqlc generate`

7.  **Run the Application**:
    -   For development with hot-reloading: `air`
    -   To build and run the binary: `go build -o bin/server cmd/server/main.go && ./bin/server`

## 4. Project Structure Overview

```
.
├── cmd/server/main.go      # Application entry point, DI setup
├── configs/                # (Legacy) Configuration files
├── db/
│   ├── migrations/         # Database migration files (*.sql)
│   └── queries/            # SQL files for SQLC code generation
├── docs/                   # Swagger/OpenAPI documentation
├── internal/
│   ├── infrastructure/     # Shared components (DB, config, middleware)
│   ├── <domain>module/     # Each business domain is a module
│   │   ├── <domain>.go     # Domain models and interfaces
│   │   ├── handler/        # HTTP handlers and router
│   │   ├── mock/           # Generated mocks for testing
│   │   ├── repository/     # Data persistence logic
│   │   └── usecase/        # Business logic
├── tests/                  # Integration tests
├── .env.example            # Example environment variables
├── compose.yml             # Docker Compose configuration
├── Dockerfile              # Docker build file
├── go.mod                  # Go module definition
├── Makefile                # Helper commands
└── sqlc.yaml               # SQLC configuration
```

## 5. How to Add a New Domain

This is a step-by-step guide to extending the boilerplate with a new feature.

1.  **Define the Database Schema**:
    -   Create a new migration file in `db/migrations/`:
        `migrate create -ext sql -dir db/migrations -seq create_<entities>_table`
    -   Define the `<entities>` table in the `.up.sql` file and the `DROP TABLE` command in the `.down.sql` file.
    -   Apply the migration: `migrate -path db/migrations -database "$DATABASE_URL" up`

2.  **Write SQL Queries**:
    -   Create `db/queries/<domain>.sql` and add your `CREATE`, `GET`, `UPDATE`, `DELETE` queries with SQLC annotations.

3.  **Generate SQLC Code**:
    -   Run `sqlc generate`. This will create `internal/infrastructure/sqlc/<domain>.sql.go`.

4.  **Create the Domain Module**:
    -   Create a new directory: `internal/<domain>module`
    -   **`internal/<domain>module/<domain>.go`**: Define the `<Entity>` struct and the `<Domain>Repository` and `<Domain>Usecase` interfaces.
    -   **`internal/<domain>module/repository/<domain>_repository.go`**: Implement the `<Domain>Repository` interface using the generated SQLC querier.
    -   **`internal/<domain>module/usecase/<domain>_usecase.go`**: Implement the `<Domain>Usecase` interface, containing the business logic.
    -   **`internal/<domain>module/handler/<domain>_handler.go`**: Create the Fiber handler functions.
    -   **`internal/<domain>module/handler/router.go`**: Define the routes for the domain module and register the handler.

5.  **Register the New Module**:
    -   In `cmd/server/main.go`, import your new module.
    -   In the `main` function, register the repository, usecase, and handler with the dependency injection container.
    -   Register the new router with the Fiber app.

6.  **Write Tests**:
    -   Generate mocks with `uber-go/mock` for your new interfaces: `go generate ./...`
    -   Write unit tests for the handler, usecase, and repository, using mocks to isolate dependencies.
    -   Write integration tests in the `tests/integration` directory if needed.
    -   Most test should be done without real data access

## 6. Development Workflow & Commands

-   **Run all checks (lint, test, generate)**:
    ```bash
    go generate ./... && golangci-lint run --fix ./... && go clean -testcache && go test -v -race ./...
    ```
-   **Generate Mocks**: `go generate ./...`
-   **Run Tests**: `go test -v -race ./...`
-   **Run Linter**: `golangci-lint run --fix ./...`
-   **Build Docker Image**: `docker compose build`
-   **Access API Docs**: `http://localhost:<APP_PORT>/swagger`

This template provides a solid foundation for building scalable and maintainable backend services. Refer to the existing modules for practical examples.