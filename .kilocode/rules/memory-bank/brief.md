# E-commerce Management System Backend Boilerplate Instructions 🛍️

This guide provides step-by-step instructions for creating a **Backend Boilerplate** for an e-commerce management system. We'll use key technologies like **Go Fiber, golang-migrate (Postgres), SQLC, JWT, uber-go/mock, DATA-DOG/go-sqlmock, air, and Docker**, while applying **Clean Architecture and SOLID Principles** to design the system.

---

## 🏗️ Project Structure

The project will be organized using a Clean Architecture approach to make the code manageable and easy to maintain.

* **`cmd/server`**: Application entry point.
* **`internal/infrastruture`**: Holds the infrastruture for projects.
* **`internal/domain`**: Defines interfaces for communication between different layers.
* **`internal/handler`**: Manages HTTP requests and responses.
* **`internal/repository`**: Handles database interactions.
* **`internal/usecase`**: Contains business logic and calls the repository.
* **`pkg`**: Stores shared packages.
* **`compose.yml`**: Manages Docker services.
* **`Dockerfile`**: Builds the Docker image.
* **`migrations`**: Manages database migrations.

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
    * Add go generate annotation for gomock to mock interfaces in its own mock package.
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