# Architecture

## System Architecture

The application follows a Clean Architecture pattern with distinct layers to promote testability, maintainability, and separation of concerns:

- **Presentation/API Layer**: Located in `cmd/server`, initializes the Fiber router, registers middleware, and defines routes.
- **Handler Layer**: In `internal/handler`, contains HTTP handlers that translate incoming requests into use case inputs and format responses.
- **Use Case/Application Layer**: In `internal/usecase`, implements business logic and orchestrates calls to repositories and other services.
- **Domain Layer**: In `internal/domain`, defines core domain models, value objects, and interfaces for repositories and services.
- **Repository/Infrastructure Layer**: 
  - `internal/repository`: Provides concrete implementations of repository interfaces for PostgreSQL using SQLC.
  - `internal/infrastructure`: Contains external infrastructure code (e.g., email, cache, messaging).

- **Shared Utilities**: `pkg` holds reusable utilities such as logging, configuration, and error handling.
- **Database Migrations**: The `migrations` directory contains SQL migration files managed by golang-migrate.
- **Docker and Compose**: The `Dockerfile` and `compose.yml` orchestrate container builds and service definitions.

## Key Technical Decisions

- **Clean Architecture**: Enforces clear separation between API, business logic, and data access.
- **Dependency Injection**: All dependencies are injected via constructors and passed explicitly.
- **SQLC**: Generates type-safe Go code from SQL queries to prevent runtime errors.
- **JWT Authentication**: Middleware handles JWT validation and authorization across protected routes.
- **Testing Strategy**: Use `go-sqlmock` for repository tests and `gomock` for mocking domain interfaces in use case tests.
- **Hot Reload**: `air` enables live code reloading during development for rapid feedback.