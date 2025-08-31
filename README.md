# Go Fiber Backend Boilerplate

A production-ready Go web API boilerplate built with Fiber framework, PostgreSQL, and clean architecture principles.

## Features

- **Go Fiber**: Fast HTTP web framework
- **PostgreSQL**: Robust relational database
- **SQLC**: Type-safe database queries
- **Clean Architecture**: Separation of concerns with domain, usecase, and infrastructure layers
- **Docker Support**: Development environment with Docker Compose
- **Validation**: Request validation with custom error handling
- **Middleware**: CORS, logging, error handling, and recovery
- **Health Check**: Database health monitoring endpoint
- **Configuration**: Environment-based configuration management

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL (if running locally)
- golang-migrate (for database migrations)
- sqlc (for generating type-safe queries)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd gofiber-skeleton
   ```

2. **Copy environment variables**
   ```bash
   cp .env.example .env
   ```

3. **Start with Docker Compose**
   ```bash
   make docker-up
   ```

4. **Run migrations**
   ```bash
   make migrate-up
   ```

### Available Commands

```bash
# Development
make dev          # Run with hot reload (requires air)
make build        # Build the application
make run          # Run the built application

# Database
make migrate-up           # Run database migrations
make migrate-down         # Rollback database migrations
make migrate-create name=migration_name  # Create new migration
make sqlc-generate        # Generate SQLC code

# Testing
make test         # Run tests
make test-coverage # Run tests with coverage

# Docker
make docker-build # Build Docker image
make docker-up    # Start services
make docker-down  # Stop services
make docker-logs  # Show logs

# Utilities
make clean        # Clean build artifacts
make lint         # Run linter
```

## Project Structure

```
├── cmd/
│   └── server/          # Application entrypoints
├── internal/
│   ├── config/          # Configuration management
│   ├── infrastructure/  # External concerns (database, validation)
│   ├── middleware/      # HTTP middleware
│   ├── pkg/            # Shared utilities
│   └── user/           # Domain modules
│       ├── handler/    # HTTP handlers
│       ├── repository/ # Data access layer
│       ├── usecase/    # Business logic
│       └── user.go     # Domain models and interfaces
├── db/
│   ├── migrations/     # Database migrations
│   └── queries/        # SQL queries for SQLC
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## API Endpoints

### Health Check
- `GET /health` - Check server and database health

### Users
- `POST /api/v1/users/register` - Register a new user

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_PORT | Server port | 8080 |
| APP_ENV | Environment (development/production) | development |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | gofiber_skeleton |
| DB_SSL_MODE | Database SSL mode | disable |
| JWT_SECRET | JWT secret key | your_jwt_secret_here |
| JWT_EXPIRES_IN | JWT expiration time | 24h |
| LOG_LEVEL | Log level | info |

## Database

The project uses PostgreSQL with SQLC for type-safe database queries. Migrations are handled using golang-migrate.

### Adding New Migrations

```bash
make migrate-create name=add_new_table
```

### Generating SQLC Code

After modifying queries in `db/queries/`, regenerate the code:

```bash
make sqlc-generate
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.