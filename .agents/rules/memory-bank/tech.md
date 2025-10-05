# Technology Stack

## Core Framework & Language
- **Go**: Version 1.24.6+ (toolchain 1.25.0)
- **Fiber v2**: v2.52.9 - Express.js-inspired web framework for high-performance APIs

## Dependency Management
- **Go Modules**: Standard Go dependency management
- **Module Path**: `github.com/zercle/gofiber-skeleton`

## Key Dependencies

### Web Framework & HTTP
- `github.com/gofiber/fiber/v2` v2.52.9 - Core web framework
- `github.com/arsmn/fiber-swagger/v2` v2.31.1 - Swagger UI integration for Fiber

### Database
- **Driver**: `github.com/jackc/pgx/v5` v5.6.0 - PostgreSQL driver
- **Connection**: Uses `database/sql` with `pgx` stdlib driver
- **Query Generator**: sqlc v1.30.0 (generates type-safe Go from SQL)
- **Database**: PostgreSQL 18 (Alpine-based Docker image)

### Caching (Infrastructure)
- **Valkey**: v8 (Redis-compatible, Alpine-based)
- **Note**: Redis client not yet integrated in code, infrastructure-ready

### Authentication & Security
- `github.com/golang-jwt/jwt/v4` v4.5.2 - JWT token generation/validation
- `golang.org/x/crypto` v0.42.0 - bcrypt password hashing

### Configuration
- `github.com/joho/godotenv` v1.5.1 - .env file loading
- **Pattern**: Environment variables with fallback to .env file

### Logging
- `github.com/rs/zerolog` v1.34.0 - Structured JSON logging

### Utilities
- `github.com/google/uuid` v1.6.0 - UUID v7 generation

### API Documentation
- `github.com/swaggo/swag` v1.16.6 - Swagger documentation generator
- **Command**: `swag init -g cmd/server/main.go --output ./docs`

### Testing & Mocking
- `github.com/golang/mock` v1.6.0 - Mock generation for interfaces
- **Pattern**: `//go:generate mockgen -source=<file>.go -destination=mocks/<mock>.go -package=mocks`
- **Note**: `DATA-DOG/go-sqlmock` mentioned in brief but not yet in go.mod

## Development Tools

### Hot Reload
- **Air**: Mentioned in brief for development hot-reloading (not in dependencies yet)
- **Configuration**: Not present in codebase

### Database Migrations
- **Tool**: `golang-migrate/migrate` (mentioned in brief, not in go.mod)
- **Current**: Migrations exist in `db/migrations/`, tooling not integrated
- **Makefile**: Placeholder command `make migrate`

### SQL Code Generation
- **sqlc**: v1.30.0
- **Config**: `sqlc.yaml`
  - Engine: PostgreSQL
  - Schema: `db/migrations/`
  - Queries: `db/queries/`
  - Output: `internal/db/`
  - Features: JSON tags, empty slices, interface, exact table names off, no prepared queries, pointers for null types

### Code Quality
- **Formatter**: `go fmt` via `make fmt`
- **Linter**: `golangci-lint` via `make lint` (tool not in dependencies)
- **Testing**: Standard `go test` with race detector (`make test-race`)

### Build Tools
- **Make**: Automation for common tasks
- **Docker**: Multi-stage builds with Alpine base

## Infrastructure

### Containerization
- **Docker**: Multi-stage Dockerfile
  - **Builder**: golang:alpine
  - **Runtime**: alpine
  - **Optimizations**: Static linking, stripped binaries (`-ldflags="-s -w -extldflags '-static'"`)
  - **Security**: Non-root user (appuser/appgroup)

### Orchestration
- **Docker Compose**: v3 format
- **Services**:
  - `app`: Go Fiber application (port 8080)
  - `db`: PostgreSQL 18 Alpine (port 5432)
  - `redis`: Valkey 8 Alpine (port 6379)
- **Volumes**: Persistent storage for `db-data` and `redis-data`
- **Health Checks**: All services monitored

### Environment Configuration

**Required Variables** (.env or environment):
```
PORT=8080
DATABASE_DSN="host=localhost user=user password=password dbname=fiber_forum port=5432 sslmode=disable TimeZone=Asia/Shanghai"
JWT_SECRET="supersecretjwtkey"
```

**Docker Compose Defaults**:
```
POSTGRES_DB=fiber_forum
POSTGRES_USER=user
POSTGRES_PASSWORD=password
REDIS_URL=redis:6379
```

## Build & Deployment

### Local Development
```bash
make fmt          # Format code
make sqlc         # Generate DB code
make build        # Build binary to bin/server
make run          # Build and run
make test         # Run tests
make test-race    # Run tests with race detector
make lint         # Run linter
```

### CI Pipeline (Make target)
```bash
make ci           # fmt → sqlc → lint → test-race → build → generate-docs
```

### Production Build
- **Binary Location**: `bin/server`
- **Static Binary**: Yes (fully static linking)
- **Size Optimization**: Stripped symbols (`-s -w`)
- **Entry Point**: `cmd/server/main.go`

## Database Schema Management

### Migration Files
- **Location**: `db/migrations/`
- **Naming**: `{number}_{description}.up.sql`
- **Existing**:
  - `001_create_users.up.sql`
  - `002_create_roles.up.sql`
  - `003_create_threads_posts_comments.up.sql`
  - `004_create_sessions.up.sql`

### Query Files
- **Location**: `db/queries/`
- **Format**: sqlc annotations (`-- name: FunctionName :one/:many/:exec`)
- **Existing**: `users.sql`, `posts.sql`

## Network Architecture

### Development (Docker Compose)
- **App**: http://localhost:8080
- **Swagger**: http://localhost:8080/swagger/
- **Health**: http://localhost:8080/health
- **Readiness**: http://localhost:8080/ready
- **PostgreSQL**: localhost:5432
- **Redis/Valkey**: localhost:6379

### API Versioning
- **Base Path**: `/api/v1`
- **Auth Routes**: `/api/v1/auth/*`
- **Post Routes**: `/api/v1/posts/*`

## Security Measures

### Authentication
- **Method**: JWT Bearer tokens in Authorization header
- **Token Expiry**: 72 hours
- **Password Hashing**: bcrypt with DefaultCost (10)

### Rate Limiting
- **API**: General rate limiting via `middleware.APIRateLimit()`
- **Auth**: Stricter limits via `middleware.AuthRateLimit()`

### Application Security
- **Panic Recovery**: `fiber/v2/middleware/recover`
- **Request ID**: Unique ID per request for tracing
- **Graceful Shutdown**: 30-second timeout for in-flight requests

## Monitoring & Observability

### Logging
- **Format**: Structured JSON (zerolog)
- **Fields**: request_id, method, path, ip, status, duration, body_size, user_agent
- **Levels**: Info (requests), Error (failures), Fatal (startup failures)

### Health Checks
- **Liveness**: `/health` - Always returns 200 if server running
- **Readiness**: `/ready` - Checks database connectivity via Ping()

## Technical Constraints

### Language Version
- **Minimum**: Go 1.24.6
- **Toolchain**: Go 1.25.0

### Database
- **Engine**: PostgreSQL (pgx driver only, no MySQL/SQLite support)
- **Connection**: `database/sql` interface, not direct pgx connection pool

### Dependency Injection
- **Current**: Manual DI in `router.go`
- **Planned**: Uber fx (mentioned in brief, not implemented)

## Missing/Planned Integrations

Based on brief requirements not yet in codebase:
1. **Air** (hot-reloading) - Mentioned but not configured
2. **golang-migrate** - Migration runner not integrated
3. **Uber fx** - DI framework mentioned but not used
4. **go-sqlmock** - For repository testing, not in go.mod
5. **Redis Client** - Infrastructure ready, client not integrated
6. **Viper** - Brief mentions it, project uses godotenv instead

## Version Control & CI/CD

### Git
- **Platform**: GitHub (module path suggests)
- **Organization**: zercle
- **Repository**: gofiber-skeleton

### Docker Registry
- **Images Used**: Official golang:alpine, alpine, postgres:18-alpine, valkey/valkey:8-alpine
- **Custom Image**: Built from Dockerfile (not published)