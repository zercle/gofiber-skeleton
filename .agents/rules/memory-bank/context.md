# Current Project Context: Go Fiber Template Repository

This document captures the current state and key contextual information of the Go Fiber Template Repository.

## 1. Project Status

The project is currently in an **initial setup/template development phase**. The core structure, foundational technologies, and essential tooling are being established to create a robust and reusable backend template.

## 2. Implemented Components & Features

Based on the `brief.md` and current file structure, the following are either implemented or in active development:

-   **Core Project Structure:** The directory layout (`cmd/`, `internal/app/`, `internal/domains/`, `internal/infrastructure/`, `internal/shared/`) is in place, reflecting the Clean Architecture principles.
-   **Configuration Management:** `internal/infrastructure/config/config.go` likely handles Viper-based configuration loading.
-   **Database Setup:**
    -   `db/migrations/` contains SQL migration files.
    -   `db/queries/` contains raw SQL queries for `sqlc`.
    -   `internal/infrastructure/database/queries/` contains generated Go code from `sqlc`.
    -   `internal/infrastructure/database/database.go` likely manages database connections.
-   **HTTP Server:** `internal/app/http/server.go` and `internal/shared/server/server.go` (along with `routes.go`) define the Fiber server setup and routing.
-   **Middleware:** `internal/infrastructure/middleware/` contains common middleware (e.g., `auth.go`, `logger.go`, `recover.go`).
-   **Dependency Injection:** `internal/app/providers/` and `internal/shared/di/` indicate the use of Uber fx for DI.
-   **JSend Responses:** `internal/shared/jsend/jsend.go` provides a standardized API response format.
-   **Validation:** `internal/infrastructure/validation/validator.go` suggests an input validation mechanism.
-   **Containerization:** `compose.yml` and `Dockerfile` are present for Docker-based development.
-   **Build/Dev Tools:** `Makefile` and `.air.toml` are configured for development tasks and hot-reloading.

## 3. Key Areas Under Development/Focus

-   **Domain Implementation:** While the structure for domains (`internal/domains/auth`, `internal/domains/posts`) exists, the full business logic, usecases (`usecases/`), and repository implementations (`repository/`) for these domains are likely still being fleshed out or serve as placeholders.
-   **Testing Strategy:** The framework for unit and repository testing is defined (mockgen, go-sqlmock), but comprehensive test coverage for all components will be an ongoing effort.
-   **Observability:** Basic logging (`internal/infrastructure/logging/`) might be present, but full distributed tracing (`internal/infrastructure/trace/`) and metrics integration are likely future enhancements.
-   **Authentication:** The `auth` domain is a "must-have feature," implying its implementation is a priority, including user registration, password hashing, and JWT handling.

## 4. Assumptions & Open Questions

-   **Database:** Assuming PostgreSQL is the primary database.
-   **Cache:** Assuming Valkey (or Redis) is the chosen caching solution.
-   **Error Handling:** A consistent, application-wide error handling strategy needs to be clearly defined and implemented across all layers.
-   **API Design:** Adherence to RESTful principles and JSend standard is expected.
-   **Security:** Input validation, secure password storage, and JWT best practices are assumed to be followed.

## 5. Next Steps

-   Continue populating the Memory Bank with detailed documentation for `tasks.md`.
-   Begin implementing specific domain logic and associated tests.
-   Refine and expand observability features.
-   Ensure comprehensive error handling.