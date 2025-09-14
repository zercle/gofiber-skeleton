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

## 4. Next Steps

- Scaffold new feature folders: `internal/user` and `internal/post` with subdirectories for handler, usecase, repository, entity, model, and tests.
- Add migration scripts for `users` and `posts` tables under `db/migrations`.
- Generate SQL query definitions under `db/queries` and run `sqlc` to produce type-safe code.
- Implement user registration and login handlers with password hashing and JWT issuance.
- Implement post handlers for CRUD operations, secured behind authentication middleware.
- Update routing in `cmd/app/main.go` to wire feature routes.
- Verify work quality with:
  ```bash
  go generate ./... && golangci-lint run --fix ./... && go clean -testcache && go test -v -race ./...