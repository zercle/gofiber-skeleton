# System Architecture

## 1. Core Principles & Architecture

The boilerplate follows a domain-driven **Clean Architecture** approach, which decouples business logic from infrastructure concerns. This ensures the application is scalable, maintainable, and testable. The core principles are based on **SOLID** principles, promoting a separation of concerns and dependency inversion.

## 2. Layered Architecture

The architecture is divided into distinct layers, each with a specific responsibility:

-   **Presentation/API Layer**: (`cmd/server`, `internal/<domain>/handler/router.go`) - This is the outermost layer, responsible for handling HTTP requests, routing, and setting up dependency injection. It's the entry point of the application.
-   **Handler Layer**: (`internal/<domain>/handler`) - This layer is responsible for validating request payloads (using `go-playground/validator`), calling the appropriate use cases, and formatting the HTTP responses (using `omniti-labs/jsend`).
-   **Use Case/Service Layer**: (`internal/<domain>/usecase`) - This layer contains all the business rules and logic. It orchestrates calls to the repository layer to perform data operations.
-   **Domain Layer**: (`internal/<domain>`) - This is the core of the application, containing the domain models and the interfaces (contracts) for the repositories and use cases. It has no dependencies on any other layer.
-   **Repository/Infrastructure Layer**: (`internal/<domain>/repository`, `internal/infrastructure/sqlc`) - This layer is responsible for data persistence. It implements the repository interfaces defined in the domain layer, using SQLC-generated queries to interact with the database. It also manages database transactions.
-   **Shared Infrastructure**: (`internal/infrastructure`) - This layer contains shared components that are used across the application, such as database connections, configuration management (Viper), and middleware.

## 3. Key Design Patterns & Decisions

-   **Dependency Injection**: The project uses `samber/do` for dependency injection, which promotes loose coupling and makes the application easier to test and maintain.
-   **SQLC for Type-Safe Queries**: SQLC is used to generate type-safe Go code from raw SQL queries. This prevents SQL injection and ensures that the queries are correct at compile time.
-   **Repository-Managed Transactions**: All database transactions are managed by the repository layer, ensuring that the business logic remains clean and focused on the business rules.
-   **Viper for Configuration**: Viper is used for configuration management, allowing the application to be configured via files (`configs/<env>.yaml`), environment variables, and `.env` files.
-   **gofiber/swagger for API Documentation**: `gofiber/swagger` is used to automatically generate OpenAPI documentation from the code, making it easy to keep the documentation up-to-date.

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

A critical implementation path is adding a new domain to the application. This process involves the following steps:

1.  **Define the Database Schema**: Create a new migration file in `db/migrations/` and apply it.
2.  **Write SQL Queries**: Create a new `.sql` file in `db/queries/` with the necessary queries for the new domain.
3.  **Generate SQLC Code**: Run `sqlc generate` to generate the Go code for the new queries.
4.  **Create the Domain Module**: Create a new directory in `internal/` for the new domain, and define the domain models, interfaces, repository, use case, and handlers.
5.  **Register the New Module**: Register the new module's components with the dependency injection container in `cmd/server/main.go`.
6.  **Write Tests**: Write unit and integration tests for the new domain. Most test should be done without real data access.