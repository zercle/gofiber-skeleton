# Common Task Workflows: Go Fiber Template Repository

This document outlines common development tasks and their associated workflows within the Go Fiber Template Repository.

## 1. Adding a New Business Domain

**Task Name:** `add_new_domain`

**Description:** Create a new self-contained business domain within the `internal/domains/` directory, adhering to Clean Architecture principles.

**Files to Update/Create:**
-   `internal/domains/<new_domain>/usecases/interface.go` (Use case interfaces)
-   `internal/domains/<new_domain>/usecases/service.go` (Use case implementation)
-   `internal/domains/<new_domain>/entities/entity.go` (Domain entities)
-   `internal/domains/<new_domain>/models/models.go` (Request/response DTOs)
-   `internal/domains/<new_domain>/repositories/interface.go` (Repository interfaces)
-   `internal/domains/<new_domain>/repositories/repository.go` (Repository implementation)
-   `internal/domains/<new_domain>/handlers/handler.go` (HTTP handlers)
-   `internal/domains/<new_domain>/handlers/routes.go` (Fiber routes)
-   `internal/domains/<new_domain>/tests/` (Test files for usecases, repository, handlers)
-   `internal/app/providers/modules.go` (Add new domain's fx module)
-   `internal/shared/server/routes.go` (Register new domain's API routes)
-   `db/migrations/` (If new database tables are needed)
-   `db/queries/<new_domain>.sql` (If new SQL queries are needed)

**Step-by-step Process:**

1.  **Create Domain Directory:** `mkdir -p internal/domains/<new_domain>/{handlers,usecases,entities,models,repositories,tests}`
2.  **Define Entities:** Create `internal/domains/<new_domain>/entities/entity.go` with core domain structs.
3.  **Define Models:** Create `internal/domains/<new_domain>/models/models.go` for API request/response DTOs.
4.  **Define Repository Interface:** Create `internal/domains/<new_domain>/repositories/interface.go` for data access operations.
5.  **Implement Repository:** Create `internal/domains/<new_domain>/repositories/repository.go` with the concrete database implementation.
6.  **Define Use Case Interface:** Create `internal/domains/<new_domain>/usecases/interface.go` for business logic.
7.  **Implement Use Case:** Create `internal/domains/<new_domain>/usecases/service.go` with the business logic, depending on the repository interface.
8.  **Define API Handlers:** Create `internal/domains/<new_domain>/handlers/handler.go` to handle HTTP requests, depending on the use case interface.
9.  **Define API Routes:** Create `internal/domains/<new_domain>/handlers/routes.go` to register Fiber routes for the handlers.
10. **Integrate with DI:** Add the new domain's components (repositories, use case, handlers) as `fx.Provide` functions in `internal/app/providers/modules.go`.
11. **Register Routes:** In `internal/shared/server/routes.go`, call the new domain's `handlers.RegisterRoutes` function.
12. **Database (Optional):** If new tables are required, create migration files in `db/migrations/` and SQL queries in `db/queries/<new_domain>.sql`. Run `make migrate up` and `make sqlc`.
13. **Write Tests:** Create unit tests for `usecases`, `repositories`, and integration tests for `handlers` in `internal/domains/<new_domain>/tests/`. Use `mockgen` for interfaces.

**Key Considerations/Gotchas:**
-   Ensure strict adherence to interface-based dependencies for testability.
-   Use `sqlc` for all database interactions to maintain type safety.
-   Remember to run `go generate ./...` after creating new interfaces for `mockgen`.
-   Update `internal/app/providers/modules.go` and `internal/shared/server/routes.go` for proper integration.

## 2. Running Database Migrations

**Task Name:** `run_migrations`

**Description:** Apply or revert database schema changes using `golang-migrate/migrate`.

**Step-by-step Process:**

1.  **Create Migration Files:**
    -   `migrate create -ext sql -dir db/migrations -seq <migration_name>`
    -   This creates `NNN_<migration_name>.up.sql` and `NNN_<migration_name>.down.sql`.
2.  **Write SQL:** Populate the `.up.sql` with schema changes and `.down.sql` with rollback logic.
3.  **Apply Migrations:** `make migrate up`
4.  **Revert Last Migration:** `make migrate down`
5.  **Force Version (DANGEROUS):** `make migrate force version=<version_number>` (Use with extreme caution).

**Key Considerations/Gotchas:**
-   Always create both `up` and `down` SQL files.
-   Test migrations thoroughly in a development environment before applying to production.
-   Ensure `init-db.sql` is up-to-date if it's used for initial database setup.

## 3. Generating Type-Safe SQL Queries

**Task Name:** `generate_sqlc_queries`

**Description:** Generate Go code from SQL queries using `sqlc`.

**Step-by-step Process:**

1.  **Write SQL Queries:** Create or modify `.sql` files in `db/queries/`.
2.  **Run sqlc:** `make sqlc`
3.  **Review Generated Code:** Inspect `internal/infrastructure/database/queries/` for the generated Go files.

**Key Considerations/Gotchas:**
-   Ensure `sqlc.yaml` is correctly configured for your database and output paths.
-   `sqlc` will overwrite existing generated files, so do not manually edit them.
-   New SQL queries must be added to `db/queries/` for `sqlc` to pick them up.

## 4. Generating Mocks for Interfaces

**Task Name:** `generate_mocks`

**Description:** Generate mock implementations for Go interfaces using `go.uber.org/mock/mockgen` for testing purposes.

**Step-by-step Process:**

1.  **Define Interface:** Create an interface (e.g., `internal/domains/<domain>/usecases/interface.go`).
2.  **Add `//go:generate` Directive:** Add a comment like `//go:generate mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mocks` above the interface definition.
3.  **Run Go Generate:** `go generate ./...` or `make mocks` (if a Makefile target exists).
4.  **Use Mocks in Tests:** Import the generated mock package in your tests.

**Key Considerations/Gotchas:**
-   Ensure the `-package` flag in `mockgen` matches the desired package name for the mocks.
-   Mocks should be generated in a `mocks/` subdirectory within the package being mocked.
-   Always run `go generate ./...` after modifying an interface or adding a new one.

## 5. Running Tests

**Task Name:** `run_tests`

**Description:** Execute unit and integration tests for the project.

**Step-by-step Process:**

1.  **Run All Tests:** `make test`
2.  **Run Specific Package Tests:** `go test ./internal/domains/<domain>/usecases/...`
3.  **Run Specific Test Function:** `go test ./internal/domains/<domain>/usecases/... -run TestMyFunction`

**Key Considerations/Gotchas:**
-   Ensure test dependencies (e.g., mock implementations) are correctly generated.
-   Separate unit tests from integration tests for faster feedback.
-   Use `go-sqlmock` for repository tests to avoid real database dependencies.

## 6. Starting the Development Server

**Task Name:** `start_dev_server`

**Description:** Start the Go Fiber application with hot-reloading using Air.

**Step-by-step Process:**

1.  **Ensure Dependencies:** Make sure Docker and Go dependencies are installed.
2.  **Start Docker Compose:** `docker compose up -d postgres valkey` (if not already running).
3.  **Run Air:** `air` or `make run` (if a Makefile target exists).

**Key Considerations/Gotchas:**
-   Ensure `.air.toml` is correctly configured.
-   Check for port conflicts if other services are running.
-   Monitor Air's output for build errors or server restarts.