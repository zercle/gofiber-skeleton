# Go Fiber Production-Ready Template

A production-ready backend template using Go and Fiber v2 framework, implementing Clean Architecture with Domain-Driven Design principles.

## Features

- **Clean Architecture**: Strict separation of concerns with domain isolation
- **Domain-Driven Design**: Business logic organized by domains
- **Type-Safe Database**: Using `sqlc` for compile-time SQL validation
- **JWT Authentication**: Secure token-based authentication
- **API Documentation**: Auto-generated Swagger/OpenAPI docs
- **Comprehensive Testing**: Unit tests with mocks and 90%+ coverage
- **Docker Support**: Containerized development and deployment
- **Hot Reloading**: Fast development with Air
- **Database Migrations**: Version-controlled schema changes
- **Caching**: Valkey/Redis integration for performance
- **Middleware Stack**: CORS, rate limiting, logging, security headers
- **CI/CD Ready**: GitHub Actions workflows included

## Tech Stack

- **Framework**: [Fiber v2](https://gofiber.io/) - Express-inspired web framework
- **Database**: PostgreSQL 18 with [pgx](https://github.com/jackc/pgx) driver
- **Cache**: Valkey (Redis-compatible)
- **ORM**: [sqlc](https://sqlc.dev/) - Type-safe SQL code generation
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **DI**: [Samber's do](https://github.com/samber/do) - Type-safe dependency injection
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt)
- **Documentation**: [swag](https://github.com/swaggo/swag) - Swagger/OpenAPI
- **Testing**: [gomock](https://github.com/uber-go/mock) - Mock generation

## Quick Start

### Prerequisites

- Go 1.25.0 or higher
- Docker & Docker Compose
- Make (optional, but recommended)

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/zercle/gofiber-skeleton.git
   cd gofiber-skeleton
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Install development tools**
   ```bash
   make install-tools
   ```

4. **Start development environment**
   ```bash
   make docker-up    # Start PostgreSQL and Valkey
   make migrate-up   # Run database migrations
   make dev          # Start development server with hot reload
   ```

5. **Access the application**
   - API: http://localhost:3000
   - Health Check: http://localhost:3000/health
   - Swagger Docs: http://localhost:3000/swagger/

## Project Structure

```
.
├── cmd/
│   ├── server/           # Main HTTP server
│   └── migrate/          # Migration tool
├── internal/
│   ├── config/          # Configuration management
│   ├── middleware/      # HTTP middleware
│   └── domains/         # Business domains
│       └── user/        # User domain (reference implementation)
│           ├── entity/      # Domain entities
│           ├── repository/  # Data access layer
│           ├── usecase/     # Business logic
│           ├── delivery/    # HTTP handlers
│           ├── tests/       # Domain tests
│           └── mocks/       # Generated mocks
├── pkg/
│   ├── auth/            # Authentication utilities
│   ├── cache/           # Cache utilities
│   ├── database/        # Database utilities
│   ├── response/        # HTTP response formatting
│   └── validator/       # Input validation
├── db/
│   ├── migrations/      # SQL migration files
│   ├── queries/         # SQL query files for sqlc
│   └── seeds/           # Database seed data
├── docs/                # Generated API documentation
├── configs/             # Configuration files
└── scripts/             # Utility scripts
```

## Available Commands

### Development
```bash
make dev              # Start development server with hot reload
make run              # Run application without hot reload
make build            # Build production binary
make clean            # Clean build artifacts
```

### Testing
```bash
make test             # Run all tests
make test-coverage    # Run tests with coverage report
make test-unit        # Run unit tests only
```

### Code Quality
```bash
make lint             # Run linter
make fmt              # Format code
make vet              # Run go vet
```

### Database
```bash
make migrate-up       # Apply all migrations
make migrate-down     # Rollback last migration
make migrate-create NAME=migration_name  # Create new migration
make migrate-version  # Show current migration version
```

### Code Generation
```bash
make sqlc             # Generate SQL code
make swag             # Generate API documentation
make mocks            # Generate test mocks
```

### Docker
```bash
make docker-up        # Start Docker containers
make docker-down      # Stop Docker containers
make docker-logs      # Show Docker logs
make docker-build     # Build Docker image
```

## Environment Variables

Key environment variables (see `.env.example` for full list):

```bash
# Server
SERVER_PORT=3000
SERVER_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=gofiber_skeleton
DB_USER=postgres
DB_PASSWORD=postgres

# Cache
VALKEY_HOST=localhost
VALKEY_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
```

## API Documentation

Once the server is running, access the interactive API documentation at:

http://localhost:3000/swagger/

The documentation is auto-generated from code comments using Swagger annotations.

## Testing

The template includes comprehensive test coverage:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/domains/user/usecase/...
```

Tests use mocks generated via `go:generate` directives for dependency isolation.

## Adding a New Domain

To add a new business domain, follow the user domain structure:

1. Create domain directory structure:
   ```bash
   mkdir -p internal/domains/[domain]/{entity,repository,usecase,delivery,tests,mocks}
   ```

2. Implement layers following the user domain pattern:
   - `entity/`: Domain models and DTOs
   - `repository/`: Data access interfaces and implementations
   - `usecase/`: Business logic
   - `delivery/`: HTTP handlers
   - `tests/`: Unit and integration tests

3. Add database migrations in `db/migrations/`
4. Add SQL queries in `db/queries/`
5. Generate SQL code: `make sqlc`
6. Register routes in main server

See `docs/ADDING_NEW_DOMAIN.md` for detailed instructions.

## Deployment

### Docker

Build and run with Docker:

```bash
# Build image
docker build -t gofiber-skeleton:latest .

# Run container
docker run -p 3000:3000 --env-file .env gofiber-skeleton:latest
```

### Binary

Build and deploy binary:

```bash
# Build for Linux
make build-linux

# Copy binary and configs to server
# Run migrations
# Start server
./bin/gofiber-skeleton-linux
```

## CI/CD

GitHub Actions workflows are included:

- **ci.yml**: Runs tests and linting on pull requests
- **cd.yml**: Builds and deploys on main branch
- **go-ci.yml**: Go-specific CI pipeline

## Security

- Password hashing with bcrypt
- JWT token authentication
- Input validation
- SQL injection prevention via sqlc
- CORS configuration
- Rate limiting
- Security headers (helmet middleware)

## Performance

- Connection pooling for database and cache
- Efficient query patterns
- Caching layer with Valkey/Redis
- Optimized Docker images
- Compression middleware

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run linting: `make lint`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or contributions:

- GitHub Issues: https://github.com/zercle/gofiber-skeleton/issues
- Documentation: See `docs/` directory
- Example: User domain implementation in `internal/domains/user/`
