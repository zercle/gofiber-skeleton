# Product: Go Fiber Backend Mono-Repo Template

## Purpose
A production-ready starter to build HTTP APIs and services in Go using Fiber and Domain-Driven Clean Architecture, optimized for modular domains, explicit dependency injection, and clear boundaries.

## Problem Statement
Teams repeatedly bootstrap similar Go backends and re-solve cross-cutting concerns such as configuration, migrations, SQL generation, validation, authentication, observability, and testing. This template standardizes these decisions with sane defaults and a maintainable structure.

## Target Audience
- Backend engineers building greenfield APIs or internal services
- Teams seeking a consistent, testable, and scalable project skeleton
- OSS maintainers who want a sensible, opinionated baseline

## Value Proposition
- Domain isolation and clean interfaces enable safe change
- Batteries-included infrastructure: config, DI, DB, middlewares
- Testing-first approach with mocks and structured error handling
- Documentation and developer experience: hot reload and Swagger

## Success Criteria
- New domain feature can be added with minimal coupling
- Environment can be configured per deployment without code changes
- Database migrations are reproducible and versioned
- APIs are documented and testable locally and in CI

## Scope
- Provide the project structure, patterns, and tooling to ship APIs quickly
- Include auth, posts sample domains for reference in documentation
- Support Postgres by default via pgx and sqlc
- Include JWT-based authentication primitives

## Non-goals
- Full-featured CMS or admin UI
- Multi-database abstraction layer
- Vendor-specific cloud lock-in

## High-level Features
- Go Fiber v2 HTTP server with production middlewares
- Uber fx for dependency injection and lifecycle management
- Viper-based configuration with env overrides
- pgx connection pooling and migrations via golang-migrate
- sqlc for type-safe query generation
- jwt v5 for authentication and claims
- Validation via go-playground validator
- Testing with Go test, Testify, and go-sqlmock
- Hot reload via Air for local dev
- Swagger docs generation via swaggo

## Conventions

- Database configuration
  - Canonical: DB_* across local, CI, and containers
    - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SCHEMA, DB_SSLMODE
  - Fallback: DB_URL is accepted when DB_* are not provided; when tools need a URL form, construct from DB_* where available
  - Deprecated: DATABASE_URL
  - See [`tech.md`](.agents/rules/memory-bank/tech.md) for examples and CLI usage.

- Cache configuration
  - Canonical: VALKEY_* (VALKEY_HOST, VALKEY_PORT, VALKEY_PASSWORD, VALKEY_DB)
  - Fallback: VALKEY_URL is accepted when VALKEY_* are not provided
  - Deprecated: REDIS_URL

- SQL placement
  - Migrations live only under [`db/migrations/`](db/migrations/)
  - SQL query files for sqlc live only under [`db/queries/`](db/queries/)
  - Legacy locations such as [`migrations/`](migrations/) and [`internal/queries`](internal/queries) are deprecated

- Domain mocks and testing
  - Each domain owns a `mocks` package under its path: [`internal/domains/<domain>/mocks`](internal/domains/)
  - Generate mocks using go.uber.org/mock/mockgen from repository interfaces
  - Use DATA-DOG/go-sqlmock for DB-layer testing, avoiding real data access

- Source of truth
  - The [`db`](db) directory is the single source of truth for SQL schema and query files
  - Generated code locations and configuration are documented in [`sqlc.yaml`](sqlc.yaml) and reflected in [`tech.md`](.agents/rules/memory-bank/tech.md)

## Key Tools and Libraries

- HTTP Framework: Go Fiber v2
- Dependency Injection: go.uber.org/fx
- Configuration: spf13/viper (yaml, .env, runtime env vars)
- Migrations: golang-migrate/migrate
- SQL Generation: sqlc
- JSON Format: omniti-labs/jsend
- Mocking/Testing: stretchr/testify, go-sqlmock, go.uber.org/mock/mockgen
- Auth: golang-jwt/jwt
- UUID: UUIDv7 (index-friendly)
- Validation: go-playground/validator
- Hot Reload: Air
- Docs: swaggo/swag

See [`tech.md`](.agents/rules/memory-bank/tech.md) for installation, configuration, and workflow details, and [`cmd/server/main.go`](cmd/server/main.go) for the HTTP server wiring.

## Developer Experience Goals
- Single command to run dev server with hot reload
- Consistent error format using jsend
- Clear conventions for folder layout and naming
- Fast feedback through unit tests and lightweight integration tests

## API UX Principles
- Consistent resource naming and versioned routes under /api/v1
- Structured responses using jsend
- Clear validation failures and error codes
- Pagination and filtering patterns for list endpoints

## Observability and Operations
- Structured request logging with production format
- Panic recovery middleware in all environments
- Health and readiness endpoints for orchestration
- Configuration surfaced via environment variables for 12-factor compatibility

## Risks and Assumptions
- The repository currently includes documentation and scaffolding guidance; full reference implementation may be added incrementally.
- Go version pinned by module file may need alignment with your toolchain.

## Links
- Source module: https://github.com/zercle/gofiber-skeleton
- Memory Bank brief: [`brief.md`](.agents/rules/memory-bank/brief.md)