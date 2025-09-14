# Current Project Context: Go Fiber Forum Backend (Feature-Based Re-init)

This document captures the context for reinitializing the Go Fiber backend template to a feature-based folder structure and implementing a simple blog post application with authentication.

## 1. Project Goals

- Migrate existing codebase to a feature-based organization under `project/`, grouping handlers, usecases, repositories, entities, models, tests, and configs by feature.
- Implement a simple blog post feature supporting create, read, update, and delete operations.
- Add user authentication using Bcrypt for password hashing and JWT for stateless sessions.
- Organize SQL migration and query definitions under `db/migrations` and `db/queries`.

## 2. Scope

- Features: `user` and `post`.
- Shared utilities under `pkg/` (e.g., logging).
- Configuration files under `configs/`.

## 3. Context

- Starting from Clean Architecture template with `cmd/`, `internal/`, and `pkg/`.
- Database: PostgreSQL (supported by `golang-migrate/migrate` and `sqlc`).
- Testing strategy: Interface mocks for repositories, unit tests for usecases and handlers, integration tests for database.

## 4. Current Development Focus & Next Actionable Steps

- User authentication and post CRUD features have been fully implemented and tested.
- Dependencies updated in `go.mod` and `go.sum`.
- Code quality checks (linting, vetting, testing) have passed successfully.

## 5. Next Actionable Steps

- Configure CI/CD pipeline for automated builds, testing, and linting.
- Generate API documentation (Swagger/OpenAPI) for User and Post endpoints.