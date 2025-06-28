# Go Fiber Monorepo Boilerplate

This is a Go Fiber monorepo boilerplate, suitable as a template repository, for a simple online shopping service demonstrating both REST and gRPC interfaces within a single application.

## Architectural Principles

*   Clean Architecture
*   SOLID Principles

## Core Technologies

*   **Framework:** Go Fiber
*   **ORM:** GORM (with `go-sqlite3` pure Go driver)
*   **Configuration:** Viper (YAML files, environment variable overrides)
*   **Authentication:** JWT
*   **Database Migrations:** Go-migrate

## Key Features & Structure

*   **Monorepo Structure:** Organized to support multiple internal modules (e.g., `user`, `product`, `order`) demonstrating Clean Architecture layers (domain, usecase, infrastructure, delivery).
*   **API Endpoints:** Example REST and gRPC endpoints for `User`, `Product`, and `Order` modules (to be implemented).
*   **Configuration:**
    *   `config/GO_ENV.yaml` for default settings.
    *   Support for runtime environment variable overrides.
*   **Database:**
    *   SQLite database.
    *   GORM models for `User`, `Product`, `Order`.
    *   Go-migrate setup for schema management.

## Getting Started

### Prerequisites

*   Go (1.24 or higher)
*   Docker (for containerization)
*   `migrate` CLI tool: `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### Setup

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd gofiber-boilerplate
    ```

2.  **Install Go modules:**
    ```bash
    go mod tidy
    ```

3.  **Run database migrations:**
    ```bash
    make migrate-up
    ```

### Running the Application

#### Locally

```bash
make run
```

#### With Docker

```bash
make docker-up
```

### Development

*   **Build:** `make build`
*   **Test:** `make test`
*   **Clean:** `make clean`
*   **Migrate Up:** `make migrate-up`
*   **Migrate Down:** `make migrate-down`

## Project Structure

```
. \
├── cmd/app             # Main application entry point
├── config              # Configuration files
├── database            # Database migrations
├── internal            # Internal modules (user, product, order) with Clean Architecture layers
│   ├── user
│   │   ├── domain      # Data structures, entities
│   │   ├── usecase     # Business logic, interfaces
│   │   ├── infrastructure # Repositories, external services
│   │   └── delivery    # API handlers (REST, gRPC)
│   ├── product
│   └── order
├── pkg                 # Reusable packages (config, database, auth)
├── api                 # Protobuf definitions for gRPC (to be implemented)
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```
