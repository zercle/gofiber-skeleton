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
├── migrations/
├── docs/
└── docker-compose.yml
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
DATABASE_URL=postgres://user:pass@localhost:5432/db?sslmode=disable
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
REDIS_URL=redis://localhost:6379
CORS_ORIGINS=*
```

## Docker Setup
- Multi-stage Dockerfile for production builds
- Docker Compose with PostgreSQL and Redis
- Volume persistence for data