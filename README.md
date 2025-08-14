# E-commerce Management System Backend

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![GoFiber](https://img.shields.io/badge/GoFiber-v2.52.0-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql)
![Valkey](https://img.shields.io/badge/Valkey-8-DC382D?style=for-the-badge&logo=valkey)
![Docker](https://img.shields.io/badge/Docker-2596be?style=for-the-badge&logo=docker)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

---

## 🚀 Overview

This repository provides a production-ready, high-performance **E-commerce Management System Backend** built with Go, adhering strictly to **Clean Architecture** and **SOLID principles**. It serves as a robust and scalable boilerplate for developing modern e-commerce services, emphasizing maintainability, testability, and independence from external frameworks.

## ✨ Features

*   **User Authentication & Authorization**: Secure JWT-based authentication with role-based access control (Admin/Customer).
*   **Product Management**: Complete CRUD operations for products with inventory tracking.
*   **Order Management**: Full order lifecycle management with status tracking.
*   **Inventory Control**: Automatic stock management and validation.
*   **Customer Order Flow**: Seamless order creation and processing.
*   **Role-Based Access**: Different permissions for administrators and customers.
*   **Database Management**: Integrated with PostgreSQL, `sqlc`, and `golang-migrate`.
*   **API Documentation**: Automated Swagger UI for interactive API exploration.
*   **Containerization**: Docker and Docker Compose for easy setup and deployment.
*   **CI/CD Ready**: Includes GitHub Actions workflows for automated testing and linting.

## 📐 Architecture

This project follows **Clean Architecture**, promoting a clear separation of concerns:

*   **Domain Layer**: Core business entities (Product, Order, User, OrderItem).
*   **Repository Layer**: Data access interfaces and implementations.
*   **Service Layer**: Business logic and use case implementations.
*   **Delivery Layer**: HTTP handlers and routing (Go Fiber).
*   **Infrastructure**: Database, caching, and external service integrations.

### 🏗️ Clean Architecture Benefits

- **Independence**: Business logic is independent of frameworks and databases
- **Testability**: Easy to unit test business logic in isolation
- **Maintainability**: Clear separation makes code easier to understand and modify
- **Scalability**: Modular design allows for easy scaling and feature additions

## 🛠️ Tech Stack

*   **Framework**: [Go Fiber](https://gofiber.io/) - High-performance web framework
*   **Database**: [PostgreSQL](https://www.postgresql.org/) - Robust relational database
*   **Caching**: [Valkey](https://valkey.io/) - High-performance Redis-compatible cache
*   **Database Tool**: [sqlc](https://sqlc.dev/) - Type-safe SQL code generation
*   **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) - Database schema management
*   **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management
*   **API Documentation**: [Swagger (swaggo)](https://github.com/swaggo/swag) - API documentation
*   **Authentication**: JWT with role-based access control
*   **Hot Reloading**: [Air](https://github.com/cosmtrek/air) - Development hot reloading
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
    git clone <repository-url>
    cd ecommerce-backend
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
    git clone <repository-url>
    cd ecommerce-backend
    ```

2.  **Install Go Dependencies:**
    ```bash
    go mod download
    ```

3.  **Set up Database (PostgreSQL):**
    *   Ensure PostgreSQL is installed and running.
    *   Create a new database for the application (e.g., `ecommerce_db`).
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

## 🚀 API Endpoints

The API provides comprehensive endpoints for e-commerce management. You can interact with it using tools like `curl`, Postman, or through the Swagger UI.

### 🔐 Authentication Endpoints

#### Register a New User
`POST /api/v1/users/register`
```json
{
  "username": "testuser",
  "password": "strongpassword",
  "role": "customer"
}
```

#### Login User
`POST /api/v1/users/login`
```json
{
  "username": "testuser",
  "password": "strongpassword"
}
```
*Response will include a `token` for subsequent authenticated requests.*

### 📦 Product Management

All product management endpoints require a valid JWT token in the `Authorization` header (e.g., `Bearer YOUR_JWT_TOKEN`).

#### Create a Product
`POST /api/v1/products`
```json
{
  "name": "Sample Product",
  "description": "A sample product description",
  "price": 29.99,
  "stock": 100,
  "image_url": "https://example.com/image.jpg"
}
```

#### Get All Products
`GET /api/v1/products`

#### Get a Single Product
`GET /api/v1/products/{id}`

#### Update a Product
`PUT /api/v1/products/{id}`
```json
{
  "name": "Updated Product Name",
  "description": "Updated description",
  "price": 39.99,
  "stock": 150,
  "image_url": "https://example.com/new-image.jpg"
}
```

#### Delete a Product
`DELETE /api/v1/products/{id}`

#### Update Product Stock
`PATCH /api/v1/products/{id}/stock?quantity=10`

### 🛒 Order Management

#### Create a New Order (Customer)
`POST /api/v1/orders/create`
```json
{
  "user_id": "uuid-here",
  "items": [
    {
      "product_id": "product-uuid-here",
      "quantity": 2
    }
  ],
  "shipping_address": "123 Main St, City, Country"
}
```

#### Get All Orders (Admin)
`GET /api/v1/orders`

#### Get a Specific Order
`GET /api/v1/orders/{id}`

#### Get User Orders
`GET /api/v1/orders/user/{userID}`

#### Update Order Status (Admin)
`PUT /api/v1/orders/{id}/status`
```json
{
  "status": "confirmed"
}
```

### 📊 Order Status Flow

Orders follow a specific status progression:
1. **Pending** → **Confirmed** or **Cancelled**
2. **Confirmed** → **Shipped**
3. **Shipped** → **Delivered**
4. **Delivered** or **Cancelled** (final states)

## 🗄️ Database Schema

### Core Tables

- **Users**: Authentication and role management
- **Products**: Product catalog with inventory tracking
- **Orders**: Order management with status tracking
- **Order Items**: Individual items within orders

### Key Features

- **UUID Primary Keys**: Secure and globally unique identifiers
- **Timestamps**: Automatic creation and update tracking
- **Foreign Key Constraints**: Referential integrity
- **Indexes**: Performance optimization for common queries
- **Check Constraints**: Data validation at database level

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
*   `make clean`: Cleans up build artifacts.

## 🧪 Testing

The project includes comprehensive testing setup:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **Mock Generation**: Automatic mock generation for testing
- **Test Coverage**: Track test coverage and quality

Run tests with:
```bash
make test
```

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Admin and customer permission levels
- **Password Hashing**: Bcrypt password encryption
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries via SQLC

## 📈 Performance Features

- **Database Indexing**: Optimized query performance
- **Connection Pooling**: Efficient database connection management
- **Caching Layer**: Redis/Valkey integration for performance
- **Efficient Queries**: Type-safe SQL with SQLC
- **Async Processing**: Non-blocking operations where appropriate

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
*   **`sqlc` or migration errors**: Ensure your `Makefile` commands are executed in the correct order (`sqlc-generate` before `run` or `dev` if schema changes).
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

## 🚀 Roadmap

### Phase 1: Core E-commerce (✅ Completed)
- [x] User authentication and authorization
- [x] Product management
- [x] Order management
- [x] Basic inventory control

### Phase 2: Enhanced Features (🔄 In Progress)
- [ ] Payment processing integration
- [ ] Shipping and delivery tracking
- [ ] Customer reviews and ratings
- [ ] Advanced inventory management

### Phase 3: Advanced Features (📋 Planned)
- [ ] Multi-tenant support
- [ ] Analytics and reporting
- [ ] Mobile app API
- [ ] Third-party integrations

## 🙌 Credits

*   **Go Fiber**: The fast, expressive, and zero-allocation web framework for Go.
*   **sqlc**: Generates type-safe Go code from SQL queries.
*   **golang-migrate**: Database migration handling.
*   **Viper**: For robust configuration management.
*   **swaggo**: Integrates Swagger UI for API documentation.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Troubleshooting](#-troubleshooting) section
2. Review the [API Documentation](#-api-endpoints)
3. Open an [Issue](../../issues) on GitHub
4. Check the [Swagger UI](http://localhost:8080/swagger/index.html) for API details

---

**Built with ❤️ using Clean Architecture and SOLID Principles**