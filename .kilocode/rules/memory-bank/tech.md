# Tech

## Technologies Used
- Go (>=1.25)
- Fiber v2 (github.com/gofiber/fiber/v2)
- SQLC (github.com/kyleconroy/sqlc)
- golang-migrate (github.com/golang-migrate/migrate)
- PostgreSQL
- JWT (github.com/golang-jwt/jwt)
- go-sqlmock (github.com/DATA-DOG/go-sqlmock)
- gomock (github.com/golang/mock)
- Docker
- Docker Compose
- Air (github.com/cosmtrek/air)
- golangci-lint (github.com/golangci/golangci-lint)
- go-playground/validator (github.com/go-playground/validator/v10)
- Viper (github.com/spf13/viper)
- jsend (github.com/omniti-labs/jsend)
- uuidv7 (e.g., github.com/google/uuid) for index-friendly primary keys
- gofiber/swagger (github.com/gofiber/swagger): For API documentation.
- samber/do (github.com/samber/do): For application dependency injection
- github.com/guregu/null/v6: For improved null handling in SQLC models.

## Development Setup
1. Install Go (>=1.25) and set GOPATH.
2. Install Docker and Docker Compose.
3. Install golang-migrate CLI (`go install github.com/golang-migrate/migrate/cmd/migrate@latest`) or via your package manager.
4. Install SQLC (`go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`).
5. Install Air (`go install github.com/cosmtrek/air@latest`).
6. Install golangci-lint (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`).
7. Copy `.env.example` to `.env` and update environment variables; Viper will load configurations at startup.

## Commands
```bash
sqlc generate
migrate -path db/migrations -database "$DATABASE_URL" up
go build -o bin/server cmd/server/main.go
docker compose up --build
```

## Configuration
- Application configuration is managed by **Viper**, loading from environment variables and optionally from `.env` files.
- Configuration files (`configs/<ENV>.yaml`) are no longer required.
- The application validates required configuration fields at startup using `go-playground/validator` to ensure all necessary values are present.
- A shared `config` struct is used throughout the application to avoid global variables.

## Technical Constraints
- Enforce Clean Architecture layering.
- Avoid global mutable state.
- Use dependency injection for all components.
- Rely on parameterized SQL via SQLC to prevent SQL injection.
- Validate inputs using `go-playground/validator`.
- Keep migrations aligned with SQLC queries and domain models.
- Mocks live alongside their interfaces in a mock subpackage.
- Treat SQLC-generated code as entity providers; repositories should orchestrate and aggregate data before returning to use cases. The generated code is centralized in `internal/infrastructure/sqlc` and queries are located in the `db/queries` directory at the project root.
- Use UUIDv7 instead of UUIDv4 for index-friendly primary keys.
- Do not edit mock or generated files manually; use go generate to regenerate mocks and other generated code.
- Enable SQLC `emit_methods_with_db_argument: true` for repository-level transaction management.

## Testing
- **Unit Tests**: Should be located in the same package as the code they are testing, using the `_test.go` suffix.
- **Integration Tests**: Should be located in a top-level `tests/` directory, mirroring the structure of the `internal/` directory. All database interactions in integration tests must use go-sqlmock; no real database connections.
- **Mocks**: All database-related tests must use go-sqlmock for SQL interactions, and mocks for other external services. No tests should connect to a real database.
- **Verify Works:** Linting and testing with `go generate ./... && golangci-lint run --fix ./... && go clean -testcache && go test -v -race ./...` and fix any problems.