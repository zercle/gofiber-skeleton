# Go Fiber URL Shortener

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![GoFiber](https://img.shields.io/badge/GoFiber-v2.52.0-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql)
![Valkey](https://img.shields.io/badge/Valkey-8-DC382D?style=for-the-badge&logo=valkey)
![Docker](https://img.shields.io/badge/Docker-2596be?style=for-the-badge&logo=docker)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

---

## 🚀 Overview

This repository provides a production-ready, high-performance URL shortener service built with Go, adhering strictly to **Clean Architecture** and **SOLID principles**. It serves as a robust and scalable boilerplate for developing modern web services, emphasizing maintainability, testability, and independence from external frameworks.

## ✨ Features

*   **User Authentication**: Secure JWT-based authentication.
*   **URL Management**: Create, list, update, and delete short URLs.
*   **Custom Short Codes**: Authenticated users can create custom short codes.
*   **Redirection**: Fast and secure redirection from short codes to original URLs.
<!-- *   **QR Code Generation**: Automatically generates QR codes for each shortened URL. -->
*   **Caching**: Utilizes Redis for caching frequently accessed URLs.
*   **Database Management**: Integrated with PostgreSQL, `sqlc`, and `golang-migrate`.
*   **API Documentation**: Automated Swagger UI for interactive API exploration.
*   **Containerization**: Docker and Docker Compose for easy setup and deployment.
*   **CI/CD Ready**: Includes GitHub Actions workflows for automated testing and linting.

## 📐 Architecture

This project follows **Clean Architecture**, promoting a clear separation of concerns.

*   **Entities**: Core domain models.
*   **Use Cases**: Business logic and repository interfaces.
*   **Repositories**: Database and cache implementations.
*   **Delivery (HTTP)**: HTTP handlers and routing (Go Fiber).

## 🛠️ Tech Stack

*   **Framework**: [Go Fiber](https://gofiber.io/)
*   **Database**: [PostgreSQL](https://www.postgresql.org/)
*   **Caching**: [Valkey](https://valkey.io/)
*   **Database Tool**: [sqlc](https://sqlc.dev/)
*   **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
*   **Configuration**: [Viper](https://github.com/spf13/viper)
*   **API Documentation**: [Swagger (swaggo)](https://github.com/swaggo/swag)
*   **Authentication**: JWT
*   **Hot Reloading**: [Air](https://github.com/cosmtrek/air)
*   **Containerization**: [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
*   **Testing**: [testify](https://github.com/stretchr/testify) & [mock](https://pkg.go.dev/go.uber.org/mock)

## 🏁 Getting Started

### Prerequisites

*   [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/) (Recommended for easiest setup)
*   [Go (1.24+)](https://golang.org/doc/install) (Required if running without Docker)
*   [PostgreSQL](https://www.postgresql.org/download/) (Required if running without Docker)
*   [Valkey](https://valkey.io/docs/install/) (or Redis, required if running without Docker)

### Running with Docker Compose (Recommended)

This is the quickest way to get the service up and running.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/zercle/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Configure Environment Variables:**
    Copy the example environment file and modify it as needed.
    ```bash
    cp example.env .env
    ```
    Refer to the [Configuration](#-configuration) section for details on environment variables.

3.  **Start Services:**
    ```bash
    docker-compose up -d --build
    ```
    This command builds the Docker images, sets up the PostgreSQL and Valkey containers, and starts the Go Fiber application.

4.  **Access the API:**
    Once all services are up, you can access:
    *   **API Server**: `http://localhost:8080`
    *   **Swagger UI**: `http://localhost:8080/swagger/index.html` (for interactive API documentation)

### Running Locally (Without Docker)

If you prefer to run the application directly on your machine, follow these steps:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/zercle/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Install Go Dependencies:**
    ```bash
    go mod download
    ```

3.  **Set up Database (PostgreSQL):**
    *   Ensure PostgreSQL is installed and running.
    *   Create a new database for the application (e.g., `url_shortener_db`).
    *   Update your `.env` file with the correct database connection string.

4.  **Set up Caching (Valkey/Redis):**
    *   Ensure Valkey or Redis is installed and running.
    *   Update your `.env` file with the correct cache connection string.

5.  **Configure Environment Variables:**
    Copy the example environment file and modify it according to your local setup.
    ```bash
    cp example.env .env
    ```
    Refer to the [Configuration](#-configuration) section for details.

6.  **Run Database Migrations:**
    ```bash
    make migrate-up
    ```
    This will apply all necessary database schemas.

7.  **Generate SQLC Code:**
    ```bash
    make sqlc-generate
    ```
    This generates Go code from your SQL queries.

8.  **Run the Application:**
    ```bash
    make run
    ```
    Or, for hot-reloading during development:
    ```bash
    make dev
    ```

## ⚙️ Configuration

The application uses environment variables for configuration, managed via `Viper`. A `example.env` file is provided as a template. Copy it to `.env` and adjust the values as needed.

Key environment variables:

*   `APP_PORT`: The port on which the Go Fiber application will listen (e.g., `8080`).
*   `DATABASE_URL`: Connection string for PostgreSQL (e.g., `postgresql://user:password@localhost:5432/dbname?sslmode=disable`).
*   `REDIS_ADDR`: Address for Valkey/Redis cache (e.g., `localhost:6379`).
*   `JWT_SECRET`: Secret key for JWT authentication. **Change this to a strong, random value in production.**
*   `JWT_EXPIRE_MINUTES`: Expiration time for JWT tokens in minutes (e.g., `60`).
*   `BASE_URL`: The base URL for the shortened links (e.g., `http://localhost:8080`). This is used when generating the full short URL.

## 🚀 Usage Examples

The API provides endpoints for user authentication, URL shortening, and management. You can interact with it using tools like `curl`, Postman, or through the Swagger UI.

### User Authentication

#### Register a New User
`POST /api/v1/register`
```json
{
  "username": "testuser",
  "password": "strongpassword"
}
```

#### Login User
`POST /api/v1/login`
```json
{
  "username": "testuser",
  "password": "strongpassword"
}
```
*Response will include a `token` for subsequent authenticated requests.*

### URL Management

All URL management endpoints require a valid JWT token in the `Authorization` header (e.g., `Bearer YOUR_JWT_TOKEN`).

#### Create a Short URL
`POST /api/v1/urls`
```json
{
  "original_url": "https://www.example.com/very/long/url/that/needs/shortening",
  "custom_code": "my-short-code" (optional)
}
```
*If `custom_code` is not provided, a random one will be generated.*

#### Get All Short URLs for User
`GET /api/v1/urls`

#### Get a Single Short URL Details
`GET /api/v1/urls/{id}`

#### Update a Short URL
`PUT /api/v1/urls/{id}`
```json
{
  "original_url": "https://www.newexample.com",
  "custom_code": "updated-code"
}
```

#### Delete a Short URL
`DELETE /api/v1/urls/{id}`

### Redirection

To redirect to an original URL, simply access the short code:
`GET /s/{short_code}`
Example: `http://localhost:8080/s/my-short-code`

## ⚙️ Makefile Commands

*   `make setup`: Installs `air` and `swag` tools.
*   `make run`: Runs the application locally.
*   `make dev`: Runs the application with hot-reloading (using `air`).
*   `make build`: Builds the Docker image for the application.
*   `make compose-up`: Starts all services (app, db, cache) using Docker Compose in detached mode.
*   `make compose-down`: Stops and removes all services started by Docker Compose.
*   `make sqlc-generate`: Generates Go code from SQL queries using `sqlc`.
*   `make migrate-up`: Applies pending database migrations.
*   `make migrate-down`: Rolls back the last applied database migration.
*   `make test`: Runs all unit and integration tests.
*   `make lint`: Runs the Go linter (`golangci-lint`).
*   `make clean`: Cleans up build artifacts (e.g., `api/docs.go`, `api/swagger.json`, `api/swagger.yaml`).

## ⚠️ Troubleshooting

*   **"Port already in use" error**: If you're running locally and get this error, another process is likely using `APP_PORT` (default `8080`). Either change `APP_PORT` in your `.env` file or stop the conflicting process.
*   **Database connection issues**:
    *   Ensure PostgreSQL is running and accessible from your application.
    *   Check `DATABASE_URL` in your `.env` file for correct hostname, port, username, password, and database name.
    *   If running with Docker Compose, ensure the database service is healthy (`docker-compose ps`).
*   **Valkey/Redis connection issues**:
    *   Ensure Valkey/Redis is running and accessible.
    *   Check `REDIS_ADDR` in your `.env` file.
    *   If running with Docker Compose, ensure the cache service is healthy.
*   **`sqlc` or migration errors**: Ensure your `Makefile` commands are executed in the correct order (`sqlc-generate` before `run` or `dev` if schema changes). Also, ensure your database is clean for migrations if you're experiencing issues.
*   **JWT token issues**: Ensure `JWT_SECRET` is set and consistent between your application and any clients. Check token expiration times.

## 🤝 Contributing

We welcome contributions to this project! To contribute:

1.  **Fork** the repository.
2.  **Create a new branch** for your feature or bug fix (`git checkout -b feature/your-feature-name`).
3.  **Make your changes**.
4.  **Write tests** for your changes.
5.  **Ensure tests pass** (`make test`).
6.  **Run the linter** (`make lint`) and fix any issues.
7.  **Commit your changes** (`git commit -m "feat: Add new feature"`).
8.  **Push to your fork** (`git push origin feature/your-feature-name`).
9.  **Open a Pull Request** to the `main` branch of this repository.

Please ensure your code adheres to the existing architectural patterns and Go best practices.

## 🙌 Credits

*   **Go Fiber**: The fast, expressive, and zero-allocation web framework for Go.
*   **sqlc**: Generates type-safe Go code from SQL queries.
*   **golang-migrate**: Database migration handling.
*   **Viper**: For robust configuration management.
*   **swaggo**: Integrates Swagger UI for API documentation.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.