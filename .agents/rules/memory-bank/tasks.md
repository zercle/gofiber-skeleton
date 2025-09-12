# Tasks: Repetitive Workflows

This document records repeatable procedures for common tasks in the Go Fiber Backend Mono-Repo Template.

Related docs:
- Brief: [.agents/rules/memory-bank/brief.md](.agents/rules/memory-bank/brief.md)
- Architecture: [.agents/rules/memory-bank/architecture.md](.agents/rules/memory-bank/architecture.md)
- Tech: [.agents/rules/memory-bank/tech.md](.agents/rules/memory-bank/tech.md)

## Initialize a New Project from this Template

Files to modify:
- [go.mod](go.mod:1)
- [config/config.yaml](config/config.yaml) (optional)
- [.env](.env) (not committed)

Steps:
1) Update module path in [go.mod](go.mod:1) to your repository
2) Create environment config: copy [.env.example](.env.example) to [.env](.env) and set DB_* variables
3) Install tools: Air, migrate, sqlc, swag, golangci-lint
4) Initialize git remotes, CI, and docker-compose if used

Notes:
- Keep secrets out of VCS; use environment variables
- DATABASE_URL is deprecated; construct DSNs from DB_* when needed

## Add a New Domain

Files to create/modify:
- [internal/domains/yourdomain/entities](internal/domains/)
- [internal/domains/yourdomain/models](internal/domains/)
- [internal/domains/yourdomain/repositories](internal/domains/)
- [internal/domains/yourdomain/usecases](internal/domains/)
- [internal/domains/yourdomain/handlers](internal/domains/)
- [internal/domains/yourdomain/routes](internal/domains/)
- [internal/domains/yourdomain/mocks](internal/domains/)   ← generated mocks live here
- [internal/domains/yourdomain/tests](internal/domains/)
- [internal/shared/container](internal/shared/)

Steps:
1) Define core entities and domain invariants under entities/
2) Define request/response DTOs under models/
3) Define repository interfaces in repositories/ (no concrete DB code here)
4) Implement usecase interfaces in usecases/ consuming repository interfaces
5) Implement HTTP handlers mapping DTOs to usecases with validation
6) Register routes under /api/v1 in routes/
7) Wire constructors in DI container so fx can provide handlers and usecases
8) Generate/update mocks under domain/mocks using mockgen (see Generate Domain Mocks)

Notes:
- Keep domain package dependencies pointing inward (no infra leakage)
- Prefer UUIDv7 for primary IDs

## Add a New Route in an Existing Domain

Files to modify:
- [internal/domains/yourdomain/models](internal/domains/)
- [internal/domains/yourdomain/handlers](internal/domains/)
- [internal/domains/yourdomain/routes](internal/domains/)
- [pkg/utils](pkg/utils/) if adding helpers

Steps:
1) Add or update request/response DTOs
2) Add handler method and validation
3) Register route path and HTTP method in routes
4) Update usecase or introduce a new one as needed
5) Add tests for handler and usecase behaviors

## Add Middleware

Files to modify/create:
- [internal/infrastructure/middleware](internal/infrastructure/)
- [cmd/server/main.go](cmd/server/main.go)

Steps:
1) Implement middleware as fiber.Handler in middleware/
2) Configure from [internal/infrastructure/config](internal/infrastructure/) when applicable
3) Register in [cmd/server/main.go](cmd/server/main.go) before routes with recommended ordering

## Configure Database and sqlc

Files to create/modify:
- [internal/infrastructure/database](internal/infrastructure/)
- [sqlc.yaml](sqlc.yaml)
- [db/queries](db/queries)

Policy:
- All SQL query files live under [db/queries](db/queries) (flat directory)
- All migrations live under [db/migrations](db/migrations)
- DATABASE_URL is deprecated; use DB_* variables

Steps:
1) Create connection factory using pgxpool with lifecycle hooks
2) Configure [sqlc.yaml](sqlc.yaml) to:
   - Input: [db/queries](db/queries)
   - Output package path: typically [internal/infrastructure/database/queries](internal/infrastructure/database/queries)
3) Place queries under [db/queries](db/queries) only
4) Run sqlc generate and commit generated code

Notes:
- Keep SQL centralized; domain code depends on interfaces, not SQL locations
- Provide sensible pool settings and timeouts

## Write and Run a Migration

Files to modify:
- [db/migrations](db/migrations)

Steps:
1) Create a new migration pair:
   ```bash
   migrate create -ext sql -dir ./db/migrations -seq add_example_table
   ```
2) Edit up/down SQL files with schema changes
3) Apply migrations using DSN from DB_* variables:
   ```bash
   export MIGRATE_URL=$(printf "postgres://%s:%s@%s:%s/%s?sslmode=%s" \
     "$DB_USER" "$DB_PASSWORD" "$DB_HOST" "$DB_PORT" "$DB_NAME" "$DB_SSLMODE")
   migrate -path ./db/migrations -database "$MIGRATE_URL" up
   ```

Notes:
- Version schema changes per feature and keep them small
- Ensure DB_SCHEMA is applied by the application (e.g., search_path)

## Generate Domain Mocks

Files to modify/create:
- [internal/domains/<domain>/repositories/interface.go](internal/domains/)
- [internal/domains/<domain>/mocks](internal/domains/)

Steps:
1) Ensure repository interfaces are defined under the domain's repositories/
2) Generate mocks using go.uber.org/mock/mockgen:
   ```bash
   mockgen -source ./internal/domains/posts/repositories/interface.go \
     -destination ./internal/domains/posts/mocks/mock_repository.go \
     -package mocks
   ```
3) Use mocks in domain unit tests with stretchr/testify
4) For DB-layer tests, prefer DATA-DOG/go-sqlmock to simulate pgx-level interactions without real DB

Notes:
- Regenerate mocks whenever interfaces change
- Keep mocks in the domain-owned mocks/ package

## Implement JWT Issuance

Files to modify:
- [internal/infrastructure/config](internal/infrastructure/)
- [internal/domains/auth/usecases](internal/domains/)
- [internal/domains/auth/handlers](internal/domains/)

Steps:
1) Add jwt.secret and jwt.expires_in to config with env bindings
2) Implement token generation in auth usecase
3) Expose login endpoint to issue tokens and set claims
4) Add auth middleware to validate Bearer tokens

## Generate Swagger Docs

Files to modify:
- [cmd/server/main.go](cmd/server/main.go)
- [docs](docs/)

Steps:
1) Annotate handlers with Swagger comments
2) Generate docs:
   ```bash
   swag init -g cmd/server/main.go -o docs
   ```
3) Serve docs via Fiber route if desired

## Dockerize and Compose

Files to create/modify:
- [Dockerfile](Dockerfile)
- [compose.yml](compose.yml)

Steps:
1) Write multi-stage Dockerfile for server
2) Define services for app, Postgres, and Redis in compose.yml
3) Pass environment variables via compose for local dev, using DB_* keys
4) Add healthchecks and proper stop signals

## Set up CI Pipeline

Files to create/modify:
- [.github/workflows/ci.yml](.github/workflows/ci.yml)

Steps:
1) fmt and lint
2) test with race detector
3) build
4) optional docker build and push
5) run migrations using DB_* variables to build DSN

Notes:
- Cache Go modules and build artifacts to speed up CI

## Troubleshooting and Tips

- Ensure local Go toolchain matches [go.mod](go.mod:1)
- Use DB_* environment variables; avoid relying on DATABASE_URL
- Keep handlers thin and business logic in usecases
- Prefer interfaces for repositories to enable testing
## Migrate from Redis to Valkey 8

Files to modify:
- [.env.example](.env.example): replace REDIS_URL with VALKEY_URL
- [internal/infrastructure/config/config.go](internal/infrastructure/config/config.go): bind VALKEY_URL instead of REDIS_URL
- [compose.yml](compose.yml): replace redis service with valkey (image: valkey/valkey:8 or valkey/valkey:8-alpine), keep port 6379
- [docs and Memory Bank](.agents/rules/memory-bank/): ensure references use Valkey and VALKEY_URL
- Application code: if using a Redis-compatible client (e.g., github.com/redis/go-redis/v9), no code changes are typically required; update config/env usage names

Steps:
1) Environment variables
   - Replace REDIS_URL with VALKEY_URL across configs and scripts
   - Keep redis:// URL scheme; Valkey is protocol-compatible
2) Configuration loader
   - Update Viper/env bindings in [internal/infrastructure/config/config.go](internal/infrastructure/config/config.go) to read VALKEY_URL
   - Propagate new config field to any cache initializers
3) Docker Compose
   - Replace redis service with valkey:
     ```yaml
     valkey:
       image: valkey/valkey:8-alpine
       ports:
         - 6379:6379
       healthcheck:
         test: [CMD, valkey-cli, ping]
         interval: 5s
         timeout: 3s
         retries: 5
     ```
   - Update app service env from REDIS_URL to VALKEY_URL
4) Local run (optional)
   - Quick run: docker run --rm -p 6379:6379 valkey/valkey:8-alpine
5) Validation
   - From app: use existing Redis client with VALKEY_URL, call PING to verify connectivity
   - Smoke test rate limiting, session storage, and caching behaviors
6) Documentation
   - Replace textual references of Redis with Valkey where appropriate
   - Note that REDIS_URL is deprecated; VALKEY_URL is canonical

Gotchas:
- Redis Modules: Valkey does not support Redis Enterprise/Modules; ensure no module-dependent features are used
- Lua scripts: Most Lua scripts work; review for module calls
- Security: If using TLS/ACLs, confirm connection options are preserved with Valkey
- Observability: Update any dashboards/names referencing Redis to Valkey

Notes:
- Client compatibility: Valkey maintains Redis protocol compatibility; common Redis clients (including go-redis) work without code changes
- Version alignment: Adopt Valkey v8 for local and CI parity

## Refactor Code Structure for Simplicity and SOLID

Files to modify:
- .agents/rules/memory-bank/architecture.md
- .agents/rules/memory-bank/context.md
- .agents/rules/memory-bank/tasks.md
- .agents/rules/memory-bank/tech.md

Steps:
1) Define desired simplified layout and principles in memory bank.
2) Update code organization (packages, folder structure) per layout.
3) Adjust DI container registrations to reflect new module boundaries.
4) Update import paths and go.mod requires as needed.
5) Run build, tests, and validate behavior.
6) Update architecture diagram in memory bank.
7) Commit and push changes.