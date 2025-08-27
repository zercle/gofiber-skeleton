# E-commerce Management System Backend Boilerplate Instructions 🛍️

This guide provides step-by-step instructions for creating a **Backend Boilerplate** for an e-commerce management system. We'll use key technologies like **Go Fiber, golang-migrate (Postgres), SQLC, JWT, uber-go/mock, DATA-DOG/go-sqlmock, air, and Docker**, while applying **Clean Architecture and SOLID Principles** to design the system.

---

## 🏗️ Project Structure

The project will be organized using a Clean Architecture approach to make the code manageable and easy to maintain.

* **`cmd/server`**: Application entry point.
* **`internal/infrastructure`**: Holds the infrastructure for projects.
* **`internal/domain`**: Defines interfaces for communication between different layers.
* **`internal/handler`**: Manages HTTP requests and responses.
* **`internal/repository`**: Handles database interactions.
* **`internal/usecase`**: Contains business logic and calls the repository.
* **`pkg`**: Stores shared packages.
* **`migrations`**: Manages database migrations.
* **`queries`**: Stores SQL query files.
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

1.  **Design and Create Database Schema**:
    * Create a `products` table with fields like **`id`, `name`, `description`, `price`, `stock`, `image_url`, and `created_at`**.
    * Use **golang-migrate** to handle migrations.
2.  **Create API Endpoints**:
    * **`POST /api/v1/products`**: To add a new product.
    * **`PUT /api/v1/products/{id}`**: To update product information.
    * **`DELETE /api/v1/products/{id}`**: To delete a product.
    * **`GET /api/v1/products`**: To retrieve all products.
    * **`GET /api/v1/products/{id}`**: To retrieve a specific product.
3.  **Mock Interfaces**:
    * Mocks live alongside their interfaces in a mock subpackage.
4.  **Write Unit Tests**:
    * Test the functionality of each API endpoint and related business logic with mock interfaces.

### Story 2: Order Management

---

### Backend (Go)

1.  **Design and Create Database Schema**:
    * Create `orders` and `order_items` tables.
2.  **Create API Endpoints**:
    * **`GET /api/v1/orders`**: To retrieve all orders.
    * **`GET /api/v1/orders/{id}`**: To retrieve a specific order.
    * **`PUT /api/v1/orders/{id}/status`**: To update an order's status.
3.  **Write Unit Tests**:
    * Test all API and business logic functionality.

### Story 3: User Authentication & Authorization

---

### Backend (Go)

1.  **Design and Create Database Schema**:
    * Create a `users` table with fields like **`id`, `username`, `password_hash`, and `role`**.
2.  **Create API Endpoints**:
    * **`POST /api/v1/register`**: For user registration.
    * **`POST /api/v1/login`**: For user login.
3.  **Implement JWT and Middleware**:
    * Create **middleware** to validate the JWT for all authenticated API endpoints.

### Story 4: Customer Order Flow

---

### Backend (Go)

1.  **Create API Endpoint**:
    * **`POST /api/v1/orders/create`**: For customers to create a new order.
2.  **Update Logic**:
    * Modify the `POST /api/v1/products` and `PUT /api/v1/products/{id}` APIs to reduce stock count when a product is purchased.

---

### Story 5: Transaction Control Demonstration

---

#### Backend (Go)

1.  **Configure SQLC**:
    * Enable `emit_methods_with_db_argument: true` to pass the DB argument for transaction control.
2.  **Improve Model Types**:
    * Replace generated `null` types with `github.com/guregu/null/v6` for robust null handling.
3.  **Create API Endpoint**:
    * **`POST /api/v1/demo/transaction`**: Demonstrate full transaction flow (begin, operations, commit/rollback).
4.  **Write Tests**:
    * Unit tests for repository transaction methods and endpoint using mocks.

### Story 6: Complex Query Demonstration

---

#### Backend (Go)

1.  **Write Advanced SQL**:
    * Add a join query in `queries/demo_join.sql` combining orders, order_items, and products.
2.  **Generate SQLC Code**:
    * Run `sqlc generate` and use generated methods in `internal/infrastructure/sqlc`.
3.  **Create API Endpoint**:
    * **`GET /api/v1/demo/joined`**: Return joined data via SQLC-generated query.
4.  **Write Tests**:
    * Integration tests validating complex join results.

---

## 🛠️ Tools & Technologies

* **Go Fiber**: A high-performance web framework.
* **golang-migrate**: For database chema & query migration (Postgres).
* **DATA-DOG/go-sqlmock**: For simulate any sql driver behavior in tests, without needing a real database connection.
* **uber-go/mock**: For mocking interface for tests.
* **SQLC**: Generates Go code from SQL queries for increased safety and performance.
* **JWT**: For managing authentication.
* **air**: For hot reloading during development.
* **Docker & Docker-Compose**: For environment management and multi-stage builds.
* **Clean Architecture & SOLID Principles**: Design principles for creating well-structured and scalable code.
* **gofiber/swagger**: For API documentation.