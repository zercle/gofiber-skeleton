# Tech

## Technologies Used
- Go (>=1.24.6)
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

## Development Setup
1. Install Go (>=1.24.6) and set GOPATH.
2. Install Docker and Docker Compose.
3. Install migrate CLI: `brew install golang-migrate` or appropriate package manager.
4. Install SQLC: `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`.
5. Install Air: `go install github.com/cosmtrek/air@latest`.
6. Install golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`.

## Commands
```bash
sqlc generate
migrate -path migrations -database "$DATABASE_URL" up
go build -o bin/server cmd/server/main.go
docker compose up --build
```

## Configuration
- Use `DATABASE_URL` environment variable for Postgres connection.
- Use `JWT_SECRET` for signing JWT tokens.

## Technical Constraints
- Enforce Clean Architecture layering.
- Avoid global mutable state.
- Use dependency injection for all components.
- Rely on parameterized SQL via SQLC to prevent SQL injection.
- Validate inputs using Fiber middleware.
- Keep migrations aligned with SQLC queries and domain models.
- Place generated mocks in a mock subpackage within each owner package (e.g., `internal/repository/mock`, `internal/usecase/mock`).
- Treat SQLC-generated code as entity providers; repositories should orchestrate and aggregate data before returning to use cases.
- Do not edit mock or generated files manually; use go generate to regenerate mocks and other generated code.