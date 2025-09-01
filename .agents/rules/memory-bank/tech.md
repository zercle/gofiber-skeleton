# Technologies Used

## 1. Core Framework & Language

-   **Go**: The primary programming language (version >=1.25).
-   **Go Fiber**: A high-performance web framework for Go, used as the foundation for the API.

## 2. Data & Persistence

-   **PostgreSQL**: The recommended relational database for data storage.
-   **SQLC**: A tool for generating type-safe Go code from SQL queries, ensuring query correctness at compile time.
-   **golang-migrate**: A CLI tool for managing database schema migrations.

## 3. Architecture & Design

-   **Clean Architecture**: The architectural pattern used to separate concerns and create a decoupled, maintainable system.
-   **SOLID Principles**: The design principles that guide the structure of the code.
-   **samber/do**: A dependency injection library for managing dependencies and promoting loose coupling.

## 4. Development & Tooling

-   **Docker & Docker Compose**: Used for containerizing the application and its dependencies, ensuring a consistent development environment.
-   **Air**: A live-reloading tool for Go applications, used to speed up development.
-   **golangci-lint**: A linter for Go code, used to enforce code quality and consistency.
-   **uber-go/mock**: A tool for generating mock implementations of interfaces for testing. Most test should be done without real data access.
-   **Viper**: A configuration management library for handling configuration from files, environment variables, and `.env` files.
-   **gofiber/swagger**: A library for automatically generating OpenAPI documentation from the code.
-   **golang-jwt/jwt**: A library for generating and validating JSON Web Tokens (JWTs).

## 5. CI/CD & Testing

-   **GitHub Actions**: The CI/CD platform used to automate testing and linting.
-   **Go Test**: The built-in testing framework for Go, used for unit and integration tests.