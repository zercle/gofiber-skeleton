# Go URL Shortener Service: Clean Architecture Boilerplate

This project implements a high-performance URL shortening service in Go, adhering strictly to Clean Architecture and SOLID principles. It serves as a robust boilerplate for building scalable and maintainable Go applications.

## Features

*   **User Management:** Secure user registration and JWT-based authentication.
*   **URL Shortening:** Create tiny URLs with support for custom short codes for authenticated users, and guest URL creation.
*   **Redirection:** Secure redirection from short URLs to their original long URLs.
*   **QR Code Generation:** Generate QR codes for all created short URLs.

## Architecture

Built with a strong emphasis on Clean Architecture and SOLID principles, ensuring clear separation of concerns, maintainability, flexibility, and extensibility. The core business logic (use cases) is decoupled from external frameworks (Go Fiber) and database implementations.

## Technical Stack

*   **Go Framework:** [Go Fiber](https://gofiber.io/) for building efficient RESTful APIs.
*   **Database:** PostgreSQL (default).
*   **Database Interactions:**
    *   [UUIDv7](https://github.com/gofrs/uuid) for all primary keys.
    *   [sqlc](https://sqlc.dev/) for generating type-safe Go code from SQL queries.
    *   [golang-migrate](https://github.com/golang-migrate/migrate) for database schema migrations.
*   **Configuration Management:** [Viper](https://github.com/spf13/viper) for flexible application configuration.
*   **Authentication:** JWT (JSON Web Tokens).
*   **Caching/Storage:** [Valkey](https://valkey.io/) (Redis fork) for efficient data caching.
*   **Development Workflow:** [Air](https://github.com/cosmtrek/air) for live reloading.
*   **Containerization:** Multi-stage Dockerfile and Docker Compose for local development and deployment.

## Project Structure

The project follows a Clean Architecture-compliant directory structure:

```
.
├── cmd/                  # Main application entry points
│   └── api/              # Main API service entry
├── internal/             # Core application logic (delivery, usecases, repository, entities, configs)
├── pkg/                  # Reusable utilities, common libraries
├── api/                  # API specifications (e.g., Swagger/OpenAPI YAML)
├── configs/              # External application configuration files
├── db/                   # Database schema, migrations, and sqlc queries
├── tests/                # Integration tests
├── mocks/                # Generated mocks for testing
├── .github/              # GitHub Actions workflows
├── Dockerfile            # Docker build instructions
├── compose.yaml          # Docker Compose file for local development
├── go.mod                # Go module file
├── go.sum                # Go module checksums
├── Makefile              # Common development commands
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites

*   Go (1.22 or higher)
*   Docker and Docker Compose
*   `golang-migrate` CLI tool
*   `sqlc` CLI tool
*   `swag` CLI tool
*   `air` CLI tool

### Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Start Docker Compose services:**

    ```bash
    docker-compose up -d
    ```

3.  **Run database migrations:**

    ```bash
    make migrate-up
    ```

4.  **Generate sqlc code:**

    ```bash
    make sqlc
    ```

5.  **Generate Swagger documentation:**

    ```bash
    make swag
    ```

6.  **Run the application (development mode with Air):**

    ```bash
    make dev
    ```

    The API will be accessible at `http://localhost:8080`.

### Running Tests

```bash
make test
```

## CI/CD

GitHub Actions workflows are configured for:

*   Automated `go test` execution on push/pull request.
*   Automated `golangci-lint` checks for code quality and style.