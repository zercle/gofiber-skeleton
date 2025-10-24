# Go Fiber Microservice Template

A production-ready Go microservice template using the Fiber framework with Clean Architecture and Domain-Driven Design principles.

## Features

- **High Performance**: Fiber framework for ultra-fast HTTP handling
- **Clean Architecture**: Layered architecture (Handler → Usecase → Repository → Domain)
- **Type-Safe Database**: sqlc for type-safe SQL queries
- **Dependency Injection**: samber/do v2 for DI container management
- **JWT Authentication**: Secure token-based authentication
- **Configuration Management**: Viper for environment-based configuration
- **Structured Logging**: slog for structured logging with context propagation
- **Error Handling**: Comprehensive error handling with custom error types
- **Docker Support**: Multi-stage Dockerfile for optimized container builds
- **Database Migrations**: golang-migrate for version-controlled schema changes
- **Code Quality**: golangci-lint for code quality checks
- **API Documentation**: Swagger/OpenAPI ready (annotations in handlers)

## Project Structure

```
template-go-fiber/
├── cmd/service/              # Application entry point
├── internal/                 # Private application code
│   ├── handlers/             # HTTP handlers with Swagger annotations
│   ├── usecases/             # Business logic layer
│   ├── repositories/         # Data access layer (sqlc)
│   ├── domains/              # Domain interfaces and contracts
│   ├── middleware/           # HTTP middleware (auth, cors, etc.)
│   ├── infrastructure/       # Generated sqlc code
│   ├── config/               # Configuration management
│   └── errors/               # Custom error types
├── pkg/                      # Shared utilities (response, validation)
├── sql/                      # SQL files
│   ├── schema/               # Database schema
│   ├── queries/              # sqlc queries
│   └── migrations/           # Database migrations (golang-migrate)
├── docs/                     # API documentation and ADRs
├── Dockerfile                # Multi-stage container build
├── compose.yml               # Docker Compose for local development
├── Makefile                  # Development task automation
├── sqlc.yaml                 # sqlc configuration
├── .golangci.yml             # golangci-lint configuration
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
└── README.md                 # This file
```

## Quick Start

### Prerequisites

- Go 1.25 or higher
- Docker & Docker Compose (for containerized development)
- Make
- MariaDB 11+ (or use Docker Compose)

### Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd template-go-fiber
   make install-tools
   make init
   ```

2. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

3. **Start with Docker Compose**:
   ```bash
   make docker-up
   ```

4. **Run migrations**:
   ```bash
   make migrate-up
   ```

5. **Start development server**:
   ```bash
   make dev
   ```

   The server will be available at `http://localhost:3000`

### Build and Run

```bash
# Build the application
make build

# Run the binary
make run
```

## Development Workflow

### Code Generation

```bash
# Generate sqlc code from queries
sqlc generate

# Generate mocks for interfaces
go generate ./...

# Or use make
make generate
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Run linter
make lint

# Format code
make fmt
```

## API Endpoints

### Public Routes

- `POST /api/users/register` - Register a new user
- `GET /health` - Health check

### Protected Routes (Require JWT Token)

- `GET /api/users` - List users with pagination
- `GET /api/users/:id` - Get user by ID
- `GET /api/users/email?email=<email>` - Get user by email
- `PUT /api/users/:id` - Update user information
- `DELETE /api/users/:id` - Delete user

### Authentication

Include JWT token in the `Authorization` header:
```
Authorization: Bearer <token>
```

## Configuration

Configuration is managed via environment variables (see `.env.example`):

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=3000
SERVER_ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=template_go_fiber
DB_DRIVER=mysql

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## Database

### Migrations

```bash
# Run up migrations
make migrate-up

# Rollback migrations
make migrate-down
```

Migrations are located in `sql/migrations/` and use golang-migrate format.

### Queries

SQL queries for sqlc are in `sql/queries/`. After adding/modifying queries, regenerate:

```bash
make generate
```

## Docker

### Build Image

```bash
make docker-build
```

### Run with Compose

```bash
# Start services
make docker-up

# Stop services
make docker-down
```

### Manual Docker Run

```bash
docker build -t go-fiber-service:latest .
docker run -p 3000:3000 \
  -e DB_HOST=<db-host> \
  -e DB_USER=root \
  -e DB_PASSWORD=<password> \
  -e JWT_SECRET=<secret> \
  go-fiber-service:latest
```

## Architecture

### Clean Architecture Layers

1. **Handler Layer** (`internal/handlers`)
   - HTTP request/response handling
   - Input validation and parsing
   - Swagger documentation annotations

2. **Usecase Layer** (`internal/usecases`)
   - Business logic implementation
   - Input validation and orchestration
   - Error transformation

3. **Repository Layer** (`internal/repositories`)
   - Data persistence abstraction
   - Database query execution (via sqlc)
   - Transaction management

4. **Domain Layer** (`internal/domains`)
   - Entity definitions
   - Interface contracts
   - Business rules

### Dependency Injection

Uses `samber/do v2` for managing dependencies:

```go
// In config/di.go
injector, err := config.InitializeDI(cfg)
userRepo := do.MustInvoke[*repositories.UserRepository](injector)
```

### Error Handling

Custom error types in `internal/errors`:

```go
// Returns appropriate HTTP status codes
errors.NewNotFoundError("user not found")
errors.NewValidationError("invalid input", err)
errors.NewDatabaseError("query failed", err)
```

## Performance

- **Fiber**: ~25k req/s (single core) - one of the fastest Go frameworks
- **sqlc**: Type-safe, compile-time checked SQL queries - no runtime reflection
- **Connection pooling**: Configurable database connection pool
- **Graceful shutdown**: Proper resource cleanup on shutdown

## Security

- **JWT Authentication**: Token-based API security
- **Input Validation**: Request body validation
- **SQL Injection Protection**: Parameterized queries (sqlc)
- **CORS Support**: Configurable CORS headers
- **Non-root Container**: Docker image runs as non-root user

## Testing

```bash
# Unit tests with race detection
go test -v -race ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Mocks are auto-generated for all domain interfaces using:
```go
//go:generate mockgen -source=<file>.go ...
```

## Deployment

### Prerequisites

- Docker and Docker Compose
- Kubernetes (optional)
- Environment variables configured

### Kubernetes

Example deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-fiber-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-fiber-service
  template:
    metadata:
      labels:
        app: go-fiber-service
    spec:
      containers:
      - name: app
        image: go-fiber-service:latest
        ports:
        - containerPort: 3000
        env:
        - name: DB_HOST
          value: "db-service"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 10
```

## Troubleshooting

### Database Connection Issues

```bash
# Check database is running
docker logs go-fiber-db

# Verify connection
mysql -h localhost -u root -p -D template_go_fiber
```

### Port Already in Use

```bash
# Change port in .env
SERVER_PORT=3001

# Or kill existing process
lsof -ti:3000 | xargs kill -9
```

## Contributing

1. Follow Clean Architecture principles
2. Write tests for new features
3. Ensure linting passes: `make lint`
4. Add Swagger annotations to new handlers
5. Update migrations for schema changes

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or contributions, please open an issue or pull request on the repository.

## Roadmap

- [ ] GraphQL API support
- [ ] Service-to-service communication patterns
- [ ] Distributed tracing integration
- [ ] Metrics and monitoring
- [ ] Circuit breaker pattern implementation
- [ ] Cache integration (Redis/Valkey)
- [ ] Comprehensive example domain
