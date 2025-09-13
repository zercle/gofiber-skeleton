# Tech

This document captures the technologies, environment, tooling, and operational practices used in the Go Fiber Backend Mono-Repo Template.

Related docs:
- Brief: [.agents/rules/memory-bank/brief.md](.agents/rules/memory-bank/brief.md)
- Product: [.agents/rules/memory-bank/product.md](.agents/rules/memory-bank/product.md)
- Architecture: [.agents/rules/memory-bank/architecture.md](.agents/rules/memory-bank/architecture.md)
- Tasks: [.agents/rules/memory-bank/tasks.md](.agents/rules/memory-bank/tasks.md)
## Simplified Layered Architecture

Refer to the Clean Architecture diagram in [architecture.md](.agents/rules/memory-bank/architecture.md:82) for the presentation, application, domain, infrastructure, and shared layer boundaries.

## Languages and Runtime

- Go toolchain: Go 1.24.6+ (align [go.mod](go.mod:1) accordingly)
- Target: Linux-amd64 by default; cross-compilation supported via Go

## Frameworks and Libraries

- HTTP server: gofiber/fiber v2
- Dependency Injection and lifecycle: uber-go/fx
- Configuration: spf13/viper (YAML config, .env, process env)
- Database:
  - Driver and pool: jackc/pgx v5 with pgxpool
  - Migrations: golang-migrate/migrate
  - Query generation: sqlc for type-safe bindings
- JSON response format: omniti-labs/jsend
- Authentication: golang-jwt/jwt v5 (JWT with custom claims)
- Validation: go-playground/validator v10
- Testing: Go testing + stretchr/testify, go-sqlmock
- Developer experience:
  - Hot reload: cosmtrek/air
  - API docs: swaggo/swag

## Repository Layout

See the structure documented in the brief and architecture docs. Key planned entry points:
- Server entrypoint: [cmd/server/main.go](cmd/server/main.go)
- Migration runner: [cmd/migrate/main.go](cmd/migrate/main.go)
- Domain code under: [internal/domains/](internal/domains/)
- Infrastructure under: [internal/infrastructure/](internal/infrastructure/)
- Shared types and container under: [internal/shared/](internal/shared/)
- Utilities under: [pkg/utils/](pkg/utils/) (optional)
- Migrations: [db/migrations/](db/migrations/)
- SQL queries for sqlc: [db/queries/](db/queries/)
- API docs output: [docs/](docs/)

Note: The db directory is the single source of truth for all SQL migrations and query files. Legacy paths such as [migrations](migrations/) and [internal/queries](internal/queries) are deprecated.

## Quick Start

- Initialize (for new repos):
  - go mod init your-module
- Setup environment:
  - cp .env.example .env
  - Set DB_* variables; DATABASE_URL is deprecated
- Run migrations:
  - Use the migrate CLI with a DSN constructed from DB_* variables.
  - Example:
    ```bash
    export MIGRATE_URL=$(printf "postgres://%s:%s@%s:%s/%s?sslmode=%s" \
      "$DB_USER" "$DB_PASSWORD" "$DB_HOST" "$DB_PORT" "$DB_NAME" "$DB_SSLMODE")
    migrate -path ./db/migrations -database "$MIGRATE_URL" up
    ```
- Start server:
  - go run [cmd/server/main.go](cmd/server/main.go)
  - Or with hot reload: air

## Development Environment

Prerequisites:
- Go 1.24.6+ (align with [go.mod](go.mod:1))
- Git
- Optional for local services: Docker and Docker Compose
- Optional: PostgreSQL 17 and Valkey 8 running locally

Recommended global tools:
- Air: hot reload
  - Install: go install github.com/cosmtrek/air@latest
- golang-migrate CLI
  - Install: refer to https://github.com/golang-migrate/migrate
- sqlc
  - Install: refer to https://docs.sqlc.dev/en/latest/overview/install.html
- swag CLI
  - Install: go install github.com/swaggo/swag/cmd/swag@latest
- golangci-lint
  - Install: refer to https://golangci-lint.run/usage/install/

## Setup and Workflow

Initial setup:
1) Clone repository
2) Ensure Go version aligns with [go.mod](go.mod:1)
3) Download dependencies:
   - go mod download
4) Configure environment:
   - Create .env from [.env.example](.env.example) and set DB_* variables
5) Start local DB/Valkey (optional) using docker-compose or native services

Run in development:
- Simple run:
  - go run ./cmd/server
- With hot reload:
  - air

Run migrations (examples):
- Up to latest:
  - export MIGRATE_URL=$(printf "postgres://%s:%s@%s:%s/%s?sslmode=%s" "$DB_USER" "$DB_PASSWORD" "$DB_HOST" "$DB_PORT" "$DB_NAME" "$DB_SSLMODE")
  - migrate -path ./db/migrations -database "$MIGRATE_URL" up
- Rollback one:
  - migrate -path ./db/migrations -database "$MIGRATE_URL" down 1
- Create new migration pair:
  - migrate create -ext sql -dir ./db/migrations -seq add_example_table

Generate SQL code with sqlc:
- sqlc generate

## Configuration

Configuration is loaded with the following precedence:
1) Internal defaults
2) Configuration files: config.yaml in working directory or ./config
3) Environment variables (explicit bindings)

Database configuration:
- Canonical DB_* variables:
  - DB_HOST: database hostname
  - DB_PORT: database port
  - DB_USER: database user
  - DB_PASSWORD: database password
  - DB_NAME: database name
  - DB_SCHEMA: default schema (used to set search_path)
  - DB_SSLMODE: sslmode parameter for Postgres
- Fallbacks:
  - DB_URL is recognized when DB_* are not provided.
  - If both DB_* and DB_URL are set, DB_* take precedence.
- Deprecated:
  - DATABASE_URL is deprecated. When a URL is required (e.g., CLI tools), compose a DSN at runtime from DB_*.

Valkey configuration:
- Canonical VALKEY_* variables:
  - VALKEY_HOST: host
  - VALKEY_PORT: port
  - VALKEY_PASSWORD: optional password
  - VALKEY_DB: database index
- Fallbacks:
  - VALKEY_URL is recognized when VALKEY_* are not provided.
  - If both VALKEY_* and VALKEY_URL are set, VALKEY_* take precedence.
- Deprecated:
  - REDIS_URL is deprecated.

Other common environment variables:
- PORT: server port (e.g., 3000)
- ENV: environment (development, staging, production)
- JWT_SECRET: secret used to sign tokens
- JWT_EXPIRES_IN: token expiration (e.g., 24h)
- CORS_ORIGINS: CSV of allowed origins (e.g., *)

Example .env:
```env
PORT=3000
ENV=development

# Database (canonical)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gofiber_skeleton
DB_SCHEMA=public
DB_SSLMODE=disable
# Optional fallback if DB_* is not provided
DB_URL=postgres://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable

# Auth
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Cache (canonical)
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_PASSWORD=
VALKEY_DB=0
# Optional fallback if VALKEY_* is not provided
VALKEY_URL=redis://localhost:6379

# CORS
CORS_ORIGINS=*
```

Migrations example:
```bash
# Prefer DB_* by constructing MIGRATE_URL
export MIGRATE_URL=$(printf "postgres://%s:%s@%s:%s/%s?sslmode=%s" \
  "$DB_USER" "$DB_PASSWORD" "$DB_HOST" "$DB_PORT" "$DB_NAME" "$DB_SSLMODE")
# If DB_URL is provided and DB_* are not, you can use:
# export MIGRATE_URL="$DB_URL"
migrate -path ./db/migrations -database "$MIGRATE_URL" up
```

## HTTP and Middleware

Recommended order:
1) Logger
2) Recover
3) CORS
4) Auth
5) Routes

Response format:
- Use jsend for consistent success, fail, and error responses
- Map validation errors to clear messages and appropriate HTTP status codes

## Database and Persistence

- PostgreSQL 17 with pgxpool for efficient pooling
- Versioned migrations using golang-migrate, stored under [db/migrations](db/migrations)
- sqlc generates typed code from SQL files located under [db/queries](db/queries)
- Generated code is typically emitted to [internal/infrastructure/database/queries](internal/infrastructure/database/queries) as configured in [sqlc.yaml](sqlc.yaml)
- Keep domain code decoupled from SQL by depending on interfaces; SQL lives centrally in db/queries

## Authentication

- Use JWT (HS256 or as configured) with custom claims
- Bearer token verification in middleware
- Inject user context (e.g., userID, email) for downstream handlers and usecases

## Observability

- Structured logging with environment-aware formats
- Panic recovery with stack traces enabled in development
- Add health and readiness endpoints for orchestration as part of server setup
- Future: Integrate metrics/tracing (OpenTelemetry) for unified observability

## Security Guidelines

- Never commit secrets; inject via environment variables or secret managers
- Rotate JWT secrets and credentials per environment
- Validate and sanitize all inputs strictly
- Apply least-privilege DB roles for application connections
- Enable TLS at the ingress or load balancer layer in production

## CI/CD Considerations

- Suggested pipeline stages:
  - fmt and lint
  - test
  - build
  - package (Docker)
  - deploy (staging, then production)
- Cache Go module downloads and build artifacts
- Run migrations as part of deploy or pre-deploy hooks using DB_* variables

## Docker and Containers

- Use multi-stage builds for small production images
- Example runtime env vars for the app container:
  - PORT, ENV
  - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SCHEMA, DB_SSLMODE
  - JWT_SECRET, JWT_EXPIRES_IN
  - VALKEY_URL, CORS_ORIGINS
- docker-compose can provision Postgres and Valkey for local development ([compose.yml](compose.yml))
- Volume persistence for Postgres/Valkey data
- Healthchecks and proper stop signals for graceful shutdown

## Known Gaps and Next Steps

- Consolidate any legacy migrations under [migrations](migrations/) to [db/migrations](db/migrations)
- Consolidate any SQL under [internal/queries](internal/queries) or elsewhere to [db/queries](db/queries)
- Ensure [sqlc.yaml](sqlc.yaml) points to [db/queries](db/queries) for inputs and desired output package path
- Update CI and local scripts to construct MIGRATE_URL from DB_* variables
- Review [internal/infrastructure/config/config.go](internal/infrastructure/config/config.go) to bind and use DB_* variables exclusively

## Tooling Commands Reference

Development:
```bash
go run ./cmd/server
air
```

Migrations:
```bash
export MIGRATE_URL=$(printf "postgres://%s:%s@%s:%s/%s?sslmode=%s" "$DB_USER" "$DB_PASSWORD" "$DB_HOST" "$DB_PORT" "$DB_NAME" "$DB_SSLMODE")
migrate -path ./db/migrations -database "$MIGRATE_URL" up