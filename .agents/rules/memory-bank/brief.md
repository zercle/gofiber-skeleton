# Go Fiber Backend Mono-Repo Template

## Architecture
**Domain-Driven Clean Architecture** with mono-repo structure. Each domain contains complete Clean Architecture implementation with SOLID principles and domain isolation.

## Key Tools & Libraries
- **HTTP Framework:** Go Fiber v2
- **DI:** fx (Uber's framework)
- **Config:** Viper (yaml, .env, runtime env vars)
- **Migrations:** golang-migrate/migrate
- **SQL Generation:** sqlc
- **JSON Format:** omniti-labs/jsend
- **Mocking:** uber-go/mock, go-sqlmock
- **Auth:** golang-jwt
- **UUID:** UUIDv7 (index-friendly)
- **Validation:** go-playground/validator
- **Testing:** Go testing + Testify
- **Hot Reload:** Air
- **Docs:** swaggo/swag

## Project Structure
```
.
├── cmd/
│   ├── server/main.go          # Entry point
│   └── migrate/main.go         # Migration runner
├── internal/
│   ├── domains/                # Domain modules
│   │   ├── auth/               # Complete domain
│   │   │   ├── entities/
│   │   │   ├── usecases/
│   │   │   ├── repositories/
│   │   │   ├── handlers/
│   │   │   ├── routes/
│   │   │   ├── models/
│   │   │   └── tests/
│   │   └── posts/              # Another domain
│   ├── infrastructure/         # Shared infra
│   │   ├── database/
│   │   ├── middleware/
│   │   └── config/
│   └── shared/                 # Shared components
│       ├── types/
│       └── container/
├── pkg/utils/                  # Public utilities
├── db/
│   ├── queries/
│   └── migrations/
├── docs/
└── compose.yml
```

## Quick Start
```bash
# Initialize
go mod init project-name

# Setup environment
cp .env.example .env

# Run migrations
go run cmd/migrate/main.go

# Start server
go run cmd/server/main.go
# Or with hot reload
air
```

## Adding New Domain
1. Create directory structure:
   ```bash
   mkdir -p internal/{domain}/{entities,usecases,repositories,handlers,routes,models,tests}
   ```
2. Define entity with UUIDv7
3. Create models (DTOs)
4. Implement repository interface & implementation
5. Implement usecase interface & implementation
6. Create handlers
7. Define routes
8. Register dependencies in DI container
9. Write tests
10. Generate API docs

## Development Commands
```bash
# Development
go run cmd/server/main.go
air                             # Hot reload

# Database
go run cmd/migrate/main.go

# Testing
go test ./...

# Linting
golangci-lint run
gofmt -s -w .

# Documentation
swag init -g cmd/server/main.go -o docs

# Build
go build -o bin/server cmd/server/main.go
```

## Environment Variables
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
# Optional fallback if DB_* is not provided by the environment
DB_URL=postgres://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable

# Auth
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Cache (canonical)
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_PASSWORD=
VALKEY_DB=0
# Optional fallback if VALKEY_* is not provided by the environment
VALKEY_URL=redis://localhost:6379

# CORS
CORS_ORIGINS=*
# Notes:
# - DATABASE_URL is deprecated in favor of DB_* or DB_URL fallback
# - REDIS_URL is deprecated in favor of VALKEY_* or VALKEY_URL fallback
```

Precedence and fallback rules:
- Database: Prefer DB_* vars. If DB_* are not set, accept DB_URL. If both are present, DB_* take precedence. DATABASE_URL remains deprecated.
- Valkey: Prefer VALKEY_* vars. If VALKEY_* are not set, accept VALKEY_URL. If both are present, VALKEY_* take precedence. REDIS_URL remains deprecated.

## Docker Setup
- Multi-stage Dockerfile for production builds
- Docker Compose with PostgreSQL and Valkey
- Volume persistence for data