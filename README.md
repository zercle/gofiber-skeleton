# Go Fiber Backend Mono-Repo Template

This is a production-ready starter template for building HTTP APIs and services in Go using Fiber and Domain-Driven Clean Architecture.

## Features

- **Clean Architecture**: Modular domains with clear boundaries
- **Dependency Injection**: Powered by Uber fx
- **Configuration**: Viper with YAML, .env, and environment variable support
- **Database**: PostgreSQL with pgx and sqlc
- **Authentication**: JWT-based with custom claims
- **Middlewares**: Logger, Recover, CORS
- **Development**: Hot reload with Air, Swagger docs
- **Deployment**: Multi-stage Docker build

## Getting Started

1. Clone the repo
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Copy environment file:
   ```bash
   cp .env.example .env
   ```
4. Start local services:
   ```bash
   docker-compose up -d
   ```
5. Run migrations:
   ```bash
   go run cmd/migrate/main.go
   ```
6. Start server:
   ```bash
   # Development with hot reload
   air
   # Or without hot reload
   go run cmd/server/main.go
   ```

## Project Structure

See [`.agents/rules/memory-bank/brief.md`](.agents/rules/memory-bank/brief.md) for the full structure.

## Documentation

- [Product Overview](.agents/rules/memory-bank/product.md)
- [Architecture](.agents/rules/memory-bank/architecture.md)
- [Tech Stack](.agents/rules/memory-bank/tech.md)

## Planned Simplified Layout

This project is undergoing a refactor to a more simplified and SOLID-friendly package layout. New directories have been introduced to reflect this structure:

- `internal/platform`: For framework/adapter specific bootstrap logic (e.g., Fiber server).
- `internal/domains/<domain>/api`: HTTP adapters for domain-specific endpoints.
- `internal/domains/<domain>/biz`: Use cases and business logic for specific domains.
- `internal/domains/<domain>/store`: Repository implementations for specific domains.
- `internal/infrastructure/logging`: Logging adapter.
- `internal/infrastructure/trace`: Tracing adapter.
- `internal/shared/di`: Dependency Injection helpers and Fx module stubs.
- `internal/shared/pkg`: Common small utilities.
- `pkg`: Public utilities for cross-repository reuse.

For a detailed explanation of the new architecture, refer to the [Memory Bank: Architecture](.agents/rules/memory-bank/architecture.md) document.