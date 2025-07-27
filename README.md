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

*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/)
*   [Go (1.24+)](https://golang.org/doc/install) (if running without Docker)

### Running with Docker Compose (Recommended)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/zercle/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Configure:**
    ```bash
    cp example.env .env
    ```

3.  **Run:**
    ```bash
    docker-compose up -d --build
    ```

4.  **Access the API:**
    *   **API Server**: `http://localhost:8080`
    *   **Swagger UI**: `http://localhost:8080/swagger/index.html`

## ⚙️ Makefile Commands

*   `make setup`: Installs `air` and `swag`.
*   `make run`: Runs the application locally.
*   `make dev`: Runs the application with hot-reloading.
*   `make build`: Builds the Docker image.
*   `make compose-up`: Starts all services with Docker Compose.
*   `make compose-down`: Stops all services.
*   `make sqlc-generate`: Generates Go code from SQL.
*   `make migrate-up`: Applies database migrations.
*   `make migrate-down`: Rolls back the last migration.
*   `make test`: Runs all tests.
*   `make lint`: Runs the linter.
*   `make clean`: Cleans up build artifacts.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.