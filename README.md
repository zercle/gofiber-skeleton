# Go Fiber URL Shortener - Clean Architecture Boilerplate

This repository provides a production-ready boilerplate for a high-performance URL shortener service built with Go. It strictly follows Clean Architecture and SOLID principles, making it a robust foundation for any web service.

## Features

- **User Management**: Endpoints for user registration and JWT-based login.
- **URL Shortening**: Authenticated and guest users can create short URLs. Authenticated users can also specify custom short codes.
- **URL Management**: Authenticated users can list, update, and delete their own shortened URLs.
- **Secure Redirection**: Handles short URL redirection securely.
- **QR Code Generation**: Provides QR codes for each short URL.

## Architecture

This project is built upon the principles of **Clean Architecture**, ensuring a clear separation of concerns between:

- **Entities**: Core domain models.
- **Use Cases**: Business logic and rules.
- **Repositories**: Data access interfaces.
- **Delivery**: HTTP handlers and routing (Fiber).

This structure makes the application decoupled, testable, and independent of frameworks and external dependencies.

## Tech Stack

- **Framework**: [Go Fiber](https://gofiber.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Database Access**: [sqlc](https://sqlc.dev/) for type-safe SQL query generation.
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **API Documentation**: [Swagger (swaggo)](https://github.com/swaggo/swag)
- **Authentication**: JWT
- **Hot Reloading**: [Air](https://github.com/cosmtrek/air)
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/api/                # API entry point: Fiber initialization, DI
├── internal/
│   ├── delivery/http/      # HTTP handlers (Fiber), route registration
│   ├── usecases/           # Business logic, repository interfaces
│   ├── repository/         # DB implementations of interfaces
│   ├── entities/           # Domain models
│   └── configs/            # Internal configuration loaders
├── pkg/                    # Utilities, libraries
├── api/                    # OpenAPI/Swagger docs
├── configs/                # Config files (e.g., app.yaml)
├── db/
│   ├── migrations/         # DB migrations
│   └── queries/            # sqlc queries
├── tests/                  # Integration tests
├── mocks/                  # Generated mocks
├── .github/workflows/      # GitHub Actions workflows
├── Dockerfile              # Multi-stage build
├── compose.yaml            # Docker Compose for dev
├── go.mod, go.sum          # Go modules
├── Makefile                # Common commands
└── README.md               # Project documentation
```

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/gofiber-skeleton.git
    cd gofiber-skeleton
    ```

2.  **Create a `.env` file from the `example.env` file and update the values:**
    ```sh
    cp example.env .env
    ```

3.  **Run the application using Docker Compose:**
    This command will build the Docker image, start the application, PostgreSQL, and Valkey containers, and run the database migrations.
    ```sh
    docker-compose up -d --build
    ```

4.  **The application will be available at `http://localhost:8080`**.

## API Usage

### API Documentation

Full, interactive API documentation is available via Swagger UI at:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Example Requests

#### Register a new user

```sh
curl -X POST -H "Content-Type: application/json" \
-d '{"username":"testuser","password":"password"}' \
http://localhost:8080/api/users/register
```

#### Log in to get a JWT

```sh
curl -X POST -H "Content-Type: application/json" \
-d '{"username":"testuser","password":"password"}' \
http://localhost:8080/api/users/login
```

#### Create a short URL (Authenticated)

```sh
TOKEN="your-jwt-token"

curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Bearer $TOKEN" \
-d '{"original_url":"https://www.google.com"}' \
http://localhost:8080/api/urls
```

## Configuration

Application configuration is managed by `Viper`. The base configuration is located in `.env`. You can override these settings with environment variables. For example, to change the database host, you can set the `DATABASE_HOST` environment variable.

## Testing

To run the unit tests, use the following command:

```sh
go test ./...
```

The tests use `gomock` for mocking repository interfaces, ensuring that the business logic is tested in isolation.

## Makefile Commands

The following commands are available in the `Makefile` for convenience:

- `make run`: Run the application locally.
- `make dev`: Run the application with hot-reloading using Air.
- `make build`: Build the Docker image.
- `make compose-up`: Start the services with Docker Compose.
- `make compose-down`: Stop the services.
- `make sqlc-generate`: Generate Go code from SQL queries.
- `make migrate-up`: Run database migrations.
- `make migrate-down`: Roll back database migrations.

## CI/CD

A basic CI/CD pipeline is configured in `.github/workflows/ci.yml`. It automatically runs tests and linting on every push and pull request to the `main` branch.

## License

This project is licensed under the MIT License.
