# Product

## Why this Project Exists

Building a robust e-commerce backend involves repetitive setup such as routing, persistence, database migrations, authentication, and testing scaffolding. This skeleton accelerates development by providing a ready-made boilerplate with common infrastructure and patterns.

## Problems It Solves

- Eliminates boilerplate setup for REST API endpoints using Go Fiber.
- Provides database migration tooling (golang-migrate) and SQL code generation (SQLC).
- Enforces a Clean Architecture structure with clear layer separation.
- Includes authentication and authorization best practices with JWT middleware.
- Integrates testing scaffolding with mocks (uber-go/mock) and SQL mocking (go-sqlmock).

## How It Should Work

Developers clone the repository, configure environment variables, run database migrations, and start implementing business logic in the usecase and repository layers. The application can be built and run locally or in Docker, with hot reloading provided by air for rapid development.

## User Experience Goals

- Provide a quick start to a working API with minimal configuration.
- Offer a consistent and intuitive code organization to simplify onboarding.
- Ensure a seamless development workflow with hot reload and Docker support.
- Include built-in testing tools to promote confidence and code quality.
- Facilitate integration with CI/CD pipelines for continuous delivery.