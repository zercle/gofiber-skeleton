# E-commerce Management System Backend Boilerplate Instructions 🛍️

This guide provides step-by-step instructions for creating a **Backend Boilerplate** for an e-commerce management system. We'll use key technologies like **Go Fiber, golang-migrate (Postgres), SQLC, JWT, uber-go/mock, DATA-DOG/go-sqlmock, air, and Docker**, while applying **Clean Architecture and SOLID Principles** to design the system.

---

## 🏗️ Project Structure

The project will be organized using a Clean Architecture approach to make the code manageable and easy to maintain.

* **`cmd/server`**: Application entry point.
* **`internal/infrastructure`**: Holds shared infrastructure components.
* **`internal/<domain>`**: Each domain package (e.g., `internal/product`, `internal/order`, `internal/user`) contains domain models and interfaces.
* **`internal/<domain>/handler`**: Manages HTTP requests and responses for its domain.
* **`internal/<domain>/usecase`**: Contains business logic (use cases) for its domain.
* **`internal/<domain>/repository`**: Handles domain-specific data persistence.
* **`pkg`**: Stores shared utility packages.
* **`db/migrations`**: Manages database migrations.
* **`db/queries`**: Stores SQL query files for code generation.
* **`docs`**: Contains API documentation and swagger definitions.
* **`configs`**: Holds configuration files.
* **`tests`**: Contains integration test suite.
* **`compose.yml`**: Manages Docker services.
* **`Dockerfile`**: Builds the Docker image.

---

## 🎯 Epic: MVP (Minimum Viable Product)

We will start by building the essential features required for the first version of the system.

### Story 1: Product Management

---

### Backend (Go)

1. **Design and Create Database Schema**:
   * Create a `products` table with fields like **`id`, `name`, `description`, `price`, `stock`, `image_url`, and `created_at`**.
   * Use **golang-migrate** to handle migrations.
2. **Create API Endpoints**:
   * **`POST /api/v1/products`**: To add a new product.
   * **`PUT /api/v1/products/{id}`**: To update product information.
   * **`DELETE /api/v1/products/{id}`**: To delete a product.
   * **`GET /api/v1/products`**: To retrieve all products.
   * **`GET /api/v1/products/{id}`**: To retrieve a specific product.
3. **Mock Interfaces**:
   * Mocks live alongside their interfaces in a mock subpackage.
4. **Write Unit Tests**:
   * Test the functionality of each API endpoint and related business logic with mock interfaces.

### Story 2: Order Management

---

### Backend (Go)

1. **Design and Create Database Schema**:
   * Create `orders` and `order_items` tables.
2. **Create API Endpoints**:
   * **`GET /api/v1/orders`**: To retrieve all orders.
   * **`GET /api/v1/orders/{id}`**: To retrieve a specific order.
   * **`PUT /api/v1/orders/{id}/status`**: To update an order's status.
3. **Write Unit Tests**:
   * Test all API and business logic functionality.

### Story 3: User Authentication & Authorization

---

### Backend (Go)

1. **Design and Create Database Schema**:
   * Create a `users` table with fields like **`id`, `username`, `password_hash`, and `role`**.
2. **Create API Endpoints**:
   * **`POST /api/v1/register`**: For user registration.
   * **`POST /api/v1/login`**: For user login.
3. **Implement JWT and Middleware**:
   * Create **middleware** to validate the JWT for all authenticated API endpoints.

### Story 4: Customer Order Flow

---

### Backend (Go)

1. **Create API Endpoint**:
   * **`POST /api/v1/orders/create`**: For customers to create a new order.
2. **Update Logic**:
   * Modify the `POST /api/v1/products` and `PUT /api/v1/products/{id}` APIs to reduce stock count when a product is purchased.

### Story 5: Multi-Stage Query Guidance

---

### Backend (Go)

1. **Demonstrate Complex Joins**:
   * Write advanced SQL joins combining `orders`, `order_items`, and `products`.
   * Place queries in `db/queries/order_product_join.sql`.
2. **Use SQLC for Multi-Stage Queries**:
   * Configure SQLC to generate methods for join queries.
   * Centralize generated code in `internal/infrastructure/sqlc`.
3. **Implement in Domain Repositories**:
   * Include example multi-stage query usage in repository implementations.
   * Guide developers on structuring multi-step data retrieval with transactions if needed.
4. **Write Integration Tests**:
   * Validate multi-stage query results in `tests/integration`.

---

## 🛠️ Tools & Technologies

* **Go Fiber**: A high-performance web framework.
* **golang-migrate**: For database schema & query migration (Postgres).
* **DATA-DOG/go-sqlmock**: To simulate SQL driver behavior in tests, without needing a real database connection.
* **uber-go/mock**: For mocking interfaces in tests.
* **SQLC**: Generates Go code from SQL queries for increased safety and performance.
* **JWT**: For managing authentication.
* **air**: For hot reloading during development.
* **Docker & Docker-Compose**: For environment management and multi-stage builds.
* **Clean Architecture & SOLID Principles**: Design principles for creating well-structured and scalable code.
* **gofiber/swagger**: For API documentation.