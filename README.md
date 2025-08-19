# E-commerce Management System Backend Boilerplate 🛍️

This project provides a backend boilerplate for an e-commerce management system, built with Go Fiber, golang-migrate (Postgres), SQLC, JWT, uber-go/mock, DATA-DOG/go-sqlmock, air, and Docker. It follows Clean Architecture and SOLID Principles for a well-structured and maintainable codebase.

---

## 🏗️ Project Structure

The project is organized using a Clean Architecture approach:

*   **`cmd/server`**: Application entry point.
*   **`internal/infrastructure`**: Holds external infrastructure code (e.g., email, cache, messaging).
*   **`internal/domain`**: Defines interfaces for communication between different layers.
*   **`internal/handler`**: Manages HTTP requests and responses.
*   **`internal/repository`**: Handles database interactions.
*   **`internal/usecase`**: Contains business logic and calls the repository.
*   **`pkg`**: Stores shared packages.
*   **`compose.yml`**: Manages Docker services.
*   **`Dockerfile`**: Builds the Docker image.
*   **`migrations`**: Manages database migrations.

---

## 🎯 Features

### Product Management

- Design and Create Database Schema (`products` table)
- API Endpoints for CRUD operations on products.
- Mock Interfaces and Unit Tests for product functionality.

### Order Management

- Design and Create Database Schema (`orders` and `order_items` tables)
- API Endpoints for retrieving orders and updating order status.
- Unit Tests for order functionality.

### User Authentication & Authorization

- Design and Create Database Schema (`users` table)
- API Endpoints for user registration and login.
- JWT and Middleware for authentication and authorization.

### Customer Order Flow

- API Endpoint for creating new orders.
- Logic to reduce stock count when a product is purchased.

---

## 🛠️ Tools & Technologies

*   **Go Fiber**: A high-performance web framework.
*   **golang-migrate**: For database schema & query migration (Postgres).
*   **DATA-DOG/go-sqlmock**: For simulating SQL driver behavior in tests.
*   **uber-go/mock**: For mocking interfaces for tests.
*   **SQLC**: Generates Go code from SQL queries for increased safety and performance.
*   **JWT**: For managing authentication.
*   **air**: For hot reloading during development.
*   **Docker & Docker-Compose**: For environment management and multi-stage builds.
*   **Clean Architecture & SOLID Principles**: Design principles for creating well-structured and scalable code.

---

## 🚀 Getting Started

Instructions on how to set up and run the project will go here.

---

## API Endpoints

| Method | Path                        | Description                 | Authentication Required |
|--------|-----------------------------|-----------------------------|-------------------------|
| `POST`   | `/api/v1/register`          | Register a new user         | No                      |
| `POST`   | `/api/v1/login`             | Log in and get a JWT token  | No                      |
| `GET`    | `/api/v1/products`          | Retrieve all products       | No                      |
| `GET`    | `/api/v1/products/{id}`     | Retrieve a specific product | No                      |
| `POST`   | `/api/v1/products`          | Add a new product           | Yes                     |
| `PUT`    | `/api/v1/products/{id}`     | Update product information  | Yes                     |
| `DELETE` | `/api/v1/products/{id}`     | Delete a product            | Yes                     |
| `GET`    | `/api/v1/orders`            | Retrieve all orders         | Yes                     |
| `GET`    | `/api/v1/orders/{id}`       | Retrieve a specific order   | Yes                     |
| `PUT`    | `/api/v1/orders/{id}/status`| Update an order's status    | Yes                     |
| `POST`   | `/api/v1/orders/create`     | Create a new order          | Yes                     |

## Authentication Flow

### Registering a New User

To register a new user, send a `POST` request to `/api/v1/register` with the user's credentials in the request body.

### Logging In and Retrieving a JWT

After registration (or for existing users), send a `POST` request to `/api/v1/login` with the user's credentials. If successful, the API will return a JSON Web Token (JWT). This token is essential for accessing protected routes.

### Including the JWT in Requests

For all protected API endpoints, you must include the retrieved JWT in the `Authorization` header of your HTTP requests. The format should be `Authorization: Bearer <YOUR_JWT_TOKEN>`.

Example:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoiZXhhbXBsZXVzZXIifQ.signature
```

### Middleware Behavior and Protected Routes

The backend uses middleware to validate the JWT for all routes requiring authentication. If a request to a protected route does not include a valid JWT, the API will return an unauthorized error. Ensure your client-side application handles the storage and inclusion of this token for subsequent authenticated requests.