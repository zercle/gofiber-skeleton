**Generate a Go Fiber Monorepo Boilerplate**

**Project Goal:**
Create a Go Fiber monorepo boilerplate, suitable as a template repository, for a simple online shopping service demonstrating both REST and gRPC interfaces within a single application.

**Architectural Principles:**
*   Clean Architecture
*   SOLID Principles

**Core Technologies:**
*   **Framework:** Go Fiber
*   **ORM:** GORM (with `go-sqlite3` pure Go driver)
*   **Configuration:** Viper (YAML files, environment variable overrides)
*   **Authentication:** JWT
*   **Database Migrations:** Go-migrate

**Key Features & Structure:**
*   **Monorepo Structure:** Organize the project to support multiple internal modules (e.g., `user`, `product`, `order`) demonstrating Clean Architecture layers (domain, usecase, infrastructure, delivery).
*   **API Endpoints:** Implement example REST and gRPC endpoints for `User`, `Product`, and `Order` modules.
*   **Configuration:**
    *   `config/GO_ENV.yaml` for default settings.
    *   Support for runtime environment variable overrides.
*   **Database:**
    *   SQLite database.
    *   GORM models for `User`, `Product`, `Order`.
    *   Go-migrate setup for schema management.
*   **Mocks & Test:**
    *   Unit test in `tests/*` with `uber-go/mock` and `DATA-DOG/go-sqlmock`

**Tooling & DevOps:**
*   **Docker Compose:**
    *   Multi-stage Dockerfile for the Go application.
    *   Service definition for the Go application.
*   **Makefile:**
    *   Common commands (e.g., `build`, `run`, `test`, `migrate-up`, `migrate-down`, `docker-build`, `docker-up`).
*   **README.md:** Comprehensive setup and usage instructions.