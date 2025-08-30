# E-commerce Backend Boilerplate

This project provides a robust and scalable backend boilerplate for building e-commerce applications using Go. It is designed with **Clean Architecture** and **SOLID Principles** to ensure maintainability, testability, and clear separation of concerns.

## ✨ Features

-   **Go Fiber**: High-performance web framework.
-   **Clean Architecture**: Structured layers for clear separation of business logic from infrastructure.
-   **PostgreSQL**: Reliable and powerful relational database.
-   **SQLC**: Generates type-safe Go code from raw SQL queries, improving development speed and reducing errors.
-   **golang-migrate**: Database migration management for seamless schema evolution.
-   **JWT Authentication**: Secure and stateless user authentication.
-   **Uber Go Mock**: Interface mocking for robust unit testing.
-   **DATA-DOG/go-sqlmock**: Database mocking for isolated repository testing.
-   **Viper**: Flexible configuration management (environment variables, YAML).
-   **Docker & Docker Compose**: Containerization for consistent development and deployment environments.
-   **Air**: Live-reloading for rapid development.
-   **Observability**: Structured logging, metrics, and tracing for better system insights.
-   **UUIDv7**: Index-friendly primary keys for optimized database performance.

## 🚀 Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.24 or higher)
-   [Docker](https://www.docker.com/get-started)
-   [Docker Compose](https://docs.docker.com/compose/install/)
-   [golang-migrate CLI](https://github.com/golang-migrate/migrate#installation)
-   [SQLC](https://docs.sqlc.dev/en/stable/overview/install.html)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/zercle/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Set up environment variables:**
    Copy the example environment file and modify it as needed.
    ```bash
    cp .env.example .env
    ```

3.  **Start Docker containers (PostgreSQL):**
    ```bash
    docker compose up -d postgres
    ```

4.  **Run database migrations:**
    ```bash
    migrate -path db/migrations -database "postgres://user:password@localhost:5432/ecommerce?sslmode=disable" up
    ```
    **Note**: Replace `user`, `password`, and `ecommerce` with your actual database credentials from `.env`.

5.  **Generate SQLC code:**
    This command generates Go code for database interactions based on your SQL queries.
    ```bash
    sqlc generate
    ```

6.  **Run the application with Air (for hot-reloading during development):**
    ```bash
    air
    ```
    The API will be available at `http://localhost:8080`.

### Building and Running without Air

1.  **Build the application:**
    ```bash
    go build -o bin/server cmd/server/main.go
    ```

2.  **Run the compiled application:**
    ```bash
    ./bin/server
    ```

## 🧪 Running Tests

### Unit Tests

Unit tests are located alongside the code they test.
```bash
go test -v ./internal/...
```

### Integration Tests

Integration tests are located in the `tests/integration` directory. Ensure your PostgreSQL container is running before running integration tests.
```bash
go test -v ./tests/integration/...
```

### All Tests

To run all tests, including code generation and linting:
```bash
go generate ./... && golangci-lint run --fix ./... && go clean -testcache && go test -v -race ./...
```

## 🐳 Docker

The project includes a `Dockerfile` for building a production-ready Docker image.

### Build Docker Image

```bash
docker build -t gofiber-ecommerce-backend .
```

### Run Docker Container

```bash
docker run -p 8080:8080 -d gofiber-ecommerce-backend
```

## 📂 Project Structure

```
.
├── cmd/                # Application entry points
│   └── server/         # Main server application
├── configs/            # Configuration files
├── internal/           # Internal application code (Clean Architecture layers)
│   ├── domain/         # Core business entities and interfaces
│   │   └── mock/       # Generated mocks for domain interfaces
│   ├── infrastructure/ # Shared infrastructure components (database, config, SQLC generated code, middleware)
│   ├── <domain>/       # Feature-specific modules (e.g., product, order, user)
│   │   ├── handler/    # HTTP handlers and routers
│   │   ├── repository/ # Database interaction implementations
│   │   └── usecase/    # Business logic and orchestration
├── db/migrations/         # Database migration files
├── db/queries/            # SQL query files for SQLC
├── tests/              # Integration tests
├── .env.example        # Example environment variables
├── compose.yml         # Docker Compose configuration
├── Dockerfile          # Docker build instructions
├── go.mod              # Go modules file
├── go.sum              # Go module checksums
├── Makefile            # Makefile for common commands
├── README.md           # Project README
└── sqlc.yaml           # SQLC configuration
```

## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.