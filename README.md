# Go Fiber Monorepo Boilerplate

This is a Go Fiber monorepo boilerplate, suitable as a template repository, for a simple online shopping service demonstrating both REST and gRPC interfaces within a single application.

## Architectural Principles

*   Clean Architecture
*   SOLID Principles

## Core Technologies

*   **Framework:** Go Fiber
*   **ORM:** GORM (with `modernc.org/sqlite` pure Go driver)
*   **Configuration:** Viper (YAML files, environment variable overrides)
*   **Authentication:** JWT
*   **Database Migrations:** Go-migrate

## Key Features & Structure

*   **Monorepo Structure:** Organized to support multiple internal modules (e.g., `user`, `product`, `order`) demonstrating Clean Architecture layers (domain, usecase, infrastructure, delivery).
*   **API Endpoints:** Example REST and gRPC endpoints for `User`, `Product`, and `Order` modules.
*   **JSend Responses:** All REST API responses adhere to the JSend specification.
*   **Swagger Documentation:** Integrated Swagger for interactive API documentation.
*   **Configuration:**
    *   `configs/local.yaml` for default settings.
    *   Support for runtime environment variable overrides.
*   **Database:**
    *   SQLite database.
    *   GORM models for `User`, `Product`, `Order` using UUIDv7 as primary keys.
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
    cd gofiber-skeleton
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

### Authentication (JWT)

This boilerplate includes JWT-based authentication.

1.  **Set JWT Secret Key:** Ensure the `JWT_SECRET_KEY` environment variable is set (e.g., in your `.env` file or directly in your shell).
    ```bash
    export JWT_SECRET_KEY="your_super_secret_jwt_key"
    ```

2.  **Access Protected Routes:**
    *   The `/protected` endpoint is an example of a protected route.
    *   To access it, you need a valid JWT. You can generate one using the `/auth/login` (or similar) endpoint once implemented, or manually for testing purposes.

    Example using `curl` (replace `YOUR_JWT_TOKEN` with an actual token):
    ```bash
    curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/protected
    ```

### Development

*   **Build:** `make build`
*   **Test:** `make test`
*   **Clean:** `make clean`
*   **Migrate Up:** `make migrate-up`
*   **Migrate Down:** `make migrate-down`
*   **Generate Swagger Docs:** `make swagger`

### Accessing Swagger UI

Once the application is running, you can access the Swagger UI at `http://localhost:<APP_PORT>/swagger/index.html` (replace `<APP_PORT>` with the port your application is running on, default is 3000).

## Project Structure

. \
├── cmd/app             # Main application entry point
├── config              # Configuration files
├── database            # Database migrations
├── internal            # Internal modules (user, product, order) with Clean Architecture layers
│   ├── infra           # Infrastructure components (app, database, auth, config, jsend, types, middleware)
│   │   ├── app         # Fiber app setup
│   │   ├── auth        # JWT authentication
│   │   ├── config      # Configuration loading
│   │   ├── database    # Database connection and migrations
│   │   ├── jsend       # JSend response formatting
│   │   ├── middleware  # Fiber middleware
│   │   └── types       # Custom types (e.g., UUIDv7)
│   ├── user
│   │   ├── domain      # Data structures, entities
│   │   ├── usecase     # Business logic, interfaces
│   │   ├── infrastructure # Repositories, external services
│   │   └── delivery    # API handlers (REST, gRPC)
│   ├── product
│   └── order
├── api                 # Protobuf definitions for gRPC (to be implemented)
├── Dockerfile
├── compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md


