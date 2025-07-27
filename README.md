# Go Fiber URL Shortener - Clean Architecture Boilerplate

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![GoFiber](https://img.shields.io/badge/GoFiber-v2.52.0-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7.2-DC382D?style=for-the-badge&logo=redis)
![Docker](https://img.shields.io/badge/Docker-2596be?style=for-the-badge&logo=docker)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## 🚀 Overview

This repository provides a production-ready, high-performance URL shortener service built with Go, adhering strictly to **Clean Architecture** and **SOLID principles**. It serves as a robust and scalable boilerplate for developing modern web services, emphasizing maintainability, testability, and independence from external frameworks.

## ✨ Key Features

- **User Authentication**: Secure JWT-based authentication for user registration and login.
- **URL Management**:
    - Create short URLs for authenticated and guest users.
    - Authenticated users can specify custom short codes.
    - Comprehensive CRUD operations (List, Update, Delete) for user-owned URLs.
- **Efficient Redirection**: Fast and secure redirection from short codes to original URLs.
- **QR Code Generation**: Automatically generates QR codes for each shortened URL, enhancing usability.
- **Robust Caching**: Utilizes Redis for caching frequently accessed URLs, significantly improving read performance.
- **Database Management**: Integrated with PostgreSQL for reliable data storage, managed with `sqlc` for type-safe queries and `golang-migrate` for schema evolution.
- **Structured Logging**: Implements structured logging for better observability and debugging.
- **API Documentation**: Automated Swagger UI for interactive API exploration and testing.
- **Containerization**: Docker and Docker Compose for easy setup and deployment across environments.
- **CI/CD Ready**: Includes GitHub Actions workflows for automated testing and linting.

## 📐 Architecture

This project strictly follows **Clean Architecture**, promoting a clear separation of concerns and ensuring the application remains decoupled, testable, and independent of external dependencies.

- **Entities**: Defines the core domain models and business entities.
- **Use Cases**: Encapsulates the application's business logic and orchestrates interactions between entities and repositories. These are independent of external frameworks.
- **Repositories**: Provides abstract interfaces for data access operations, allowing different database implementations to be swapped without affecting business logic.
- **Delivery (HTTP)**: Handles incoming HTTP requests, translates them into use case inputs, and formats use case outputs into HTTP responses using the Go Fiber framework.

```
.
├── cmd/api/                # API entry point: Fiber initialization, Dependency Injection setup
├── internal/
│   ├── delivery/http/      # HTTP handlers (Fiber), route registration, middleware
│   ├── usecases/           # Business logic, application-specific rules, repository interfaces
│   ├── repository/         # Database and cache implementations of repository interfaces
│   ├── entities/           # Core domain models and data structures
│   └── configs/            # Internal application configuration loaders
├── pkg/                    # Reusable utilities, common libraries, helper functions
├── api/                    # OpenAPI/Swagger documentation definitions
├── configs/                # External configuration files (e.g., app.yaml, environment-specific)
├── db/
│   ├── migrations/         # Database schema migration scripts (up/down)
│   └── queries/            # SQL queries used by sqlc to generate type-safe Go code
├── tests/                  # Comprehensive integration and end-to-end tests
├── mocks/                  # Generated mock implementations for interfaces (used in unit tests)
├── .github/workflows/      # GitHub Actions workflows for CI/CD
├── Dockerfile              # Multi-stage Docker build definition for the application
├── compose.yaml            # Docker Compose configuration for local development environment (app, db, cache)
├── go.mod, go.sum          # Go module dependency management files
├── Makefile                # Collection of common development and deployment commands
└── README.md               # Project documentation and guide
```

## 🛠️ Tech Stack

- **Framework**: [Go Fiber](https://gofiber.io/) - An Express.js inspired web framework built on top of Fasthttp.
- **Database**: [PostgreSQL](https://www.postgresql.org/) - Powerful, open-source relational database.
- **Caching**: [Valkey (Redis Fork)](https://valkey.io/) - In-memory data structure store, used for URL caching.
- **ORM/Database Tool**: [sqlc](https://sqlc.dev/) - Generates type-safe Go code from SQL queries.
- **Database Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) - Database migration handling.
- **Configuration Management**: [Viper](https://github.com/spf13/viper) - Go configuration solution.
- **API Documentation**: [Swagger (swaggo)](https://github.com/swaggo/swag) - Automatically generates Swagger API documentation.
- **Authentication**: JWT (JSON Web Tokens) - For secure API authentication.
- **Hot Reloading**: [Air](https://github.com/cosmtrek/air) - Live-reloading for Go applications during development.
- **Containerization**: [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/) - For containerized development and deployment.
- **Testing**: [stretchr/testify](https://github.com/stretchr/testify) and [go.uber.org/mock](https://pkg.go.dev/go.uber.org/mock) - Comprehensive testing utilities and mock generation.

## 🏁 Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed on your system:

-   [Docker](https://docs.docker.com/get-docker/)
-   [Docker Compose](https://docs.docker.com/compose/install/)
-   [Go (1.22+)](https://golang.org/doc/install) (if running without Docker)

### Installation & Running Locally (with Docker Compose)

This is the recommended way to get started quickly.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Configure Environment Variables:**
    Copy the example environment file and update it with your desired settings.
    ```bash
    cp example.env .env
    # Open .env in your editor and modify variables like DATABASE_URL, REDIS_ADDR, etc.
    ```

3.  **Start the Application Services:**
    This command will build the Docker images, start the application, PostgreSQL database, and Valkey (Redis) cache containers, and automatically run database migrations.
    ```bash
    docker-compose up -d --build
    ```
    To view logs:
    ```bash
    docker-compose logs -f
    ```

4.  **Access the Application:**
    The API server will be available at `http://localhost:8080`.

### Running Locally (without Docker)

If you prefer to run the Go application directly on your host machine:

1.  **Ensure PostgreSQL and Redis are running and accessible** (e.g., via Docker Compose or local installations). Update your `.env` file with the correct connection strings.
2.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Run database migrations:**
    ```bash
    make migrate-up
    ```
4.  **Run the application:**
    ```bash
    make run
    # For hot-reloading during development:
    # make dev
    ```

## 🚀 API Usage

### API Documentation (Swagger UI)

Full, interactive API documentation is available via Swagger UI, allowing you to explore endpoints and test requests directly from your browser:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Example API Requests

Here are some common API interactions using `curl`:

#### 1. Register a New User

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"securepassword"}' \
  http://localhost:8080/api/users/register
```

#### 2. Log In to Get a JWT

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"securepassword"}' \
  http://localhost:8080/api/users/login
# Save the 'token' from the response for authenticated requests
```

#### 3. Create a Short URL (Authenticated)

```bash
TOKEN="YOUR_JWT_TOKEN_HERE" # Replace with the token obtained from login
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url":"https://www.google.com","short_code":"mygoogle"}' \
  http://localhost:8080/api/urls
```

#### 4. Redirect to Original URL

```bash
curl -L http://localhost:8080/mygoogle # Use -L to follow redirects
```

#### 5. Get User's URLs

```bash
TOKEN="YOUR_JWT_TOKEN_HERE"
curl -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/urls
```

## ⚙️ Configuration

Application configuration is managed by `Viper`, which loads settings from `.env` files and environment variables.

-   **`.env`**: The primary configuration file (copy from `example.env`).
-   **Environment Variables**: Override `.env` settings by defining environment variables (e.g., `DATABASE_HOST=my-db-host`).

Key configuration parameters include database connection strings, Redis addresses, JWT secrets, and server port.

## 🧪 Testing

The project includes unit and integration tests to ensure code quality and functionality.

-   **Unit Tests**: Located alongside their respective packages (e.g., `internal/repository/url_repository_test.go`). They use `go.uber.org/mock` for mocking dependencies, allowing isolated testing of business logic.
-   **Integration Tests**: Found in the `tests/` directory (e.g., `tests/integration_test.go`). These tests interact with a running database and Redis instance, providing end-to-end verification of API endpoints and data flows.

To run all tests:

```bash
go test ./...
```

## 📜 Makefile Commands

For convenience, a `Makefile` is provided with common development and operational commands:

-   `make setup`: Installs `air` and `swag` tools.
-   `make run`: Runs the Go application locally (requires local DB/Redis).
-   `make dev`: Runs the Go application locally with hot-reloading using `Air`.
-   `make build`: Builds the Docker image for the application.
-   `make compose-up`: Starts all services (app, db, redis) using Docker Compose in detached mode.
-   `make compose-down`: Stops and removes all services defined in `compose.yaml`.
-   `make sqlc-generate`: Generates Go code from SQL queries using `sqlc`.
-   `make migrate-up`: Applies all pending database migrations.
-   `make migrate-down`: Rolls back the last applied database migration.
-   `make test`: Runs all unit and integration tests.
-   `make lint`: Runs `golangci-lint` to check code quality and style.
-   `make clean`: Cleans up build artifacts and temporary files.

## 🚀 CI/CD

A basic CI/CD pipeline is configured using **GitHub Actions** in `.github/workflows/ci.yml`. This workflow automatically:

-   Runs unit and integration tests on every push and pull request to the `main` branch.
-   Performs linting checks to maintain code quality and style.

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature-name`).
3.  Make your changes.
4.  Write comprehensive tests for your changes.
5.  Ensure all tests pass (`make test`).
6.  Ensure your code passes linting checks (`make lint`).
7.  Commit your changes (`git commit -m 'feat: Add new feature X'`).
8.  Push to the branch (`git push origin feature/your-feature-name`).
9.  Open a Pull Request.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.
