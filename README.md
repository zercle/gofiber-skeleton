# Go Fiber Production-Ready Template

A production-ready template for building scalable Go backend services using Fiber v2 framework with Clean Architecture and Domain-Driven Design principles.

## 🚀 Features

- ✅ **Clean Architecture** - Strict domain isolation with SOLID principles
- ✅ **Fiber v2** - High-performance web framework
- ✅ **PostgreSQL** - Production-ready database with migrations
- ✅ **Type-safe SQL** - sqlc for compile-time query validation
- ✅ **JWT Authentication** - Secure stateless authentication
- ✅ **Dependency Injection** - Samber's do with generics
- ✅ **Hot Reload** - Air for fast development feedback
- ✅ **Comprehensive Testing** - Mocks, fixtures, 90%+ coverage
- ✅ **API Documentation** - Auto-generated Swagger/OpenAPI
- ✅ **Docker Support** - Multi-stage builds and Docker Compose
- ✅ **Production Ready** - Health checks, logging, security headers
- ✅ **Code Quality** - golangci-lint with comprehensive rules

## 📋 Prerequisites

- Go 1.25 or higher
- PostgreSQL 12+
- Docker & Docker Compose (optional)
- Make (optional, but recommended)

## 🏁 Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/zercle/gofiber-skeleton.git my-project
cd my-project

# Install development tools and setup environment
make setup
```

### 2. Configure Environment

Edit `.env` file with your configuration:

```bash
# Copy from example
cp .env.example .env

# Edit as needed
vim .env
```

### 3. Start Development Environment

#### Option A: With Docker Compose (Recommended)

```bash
# Start all services (PostgreSQL, Redis, App)
make docker-up

# View logs
docker-compose logs -f app
```

#### Option B: Local Development

```bash
# Start PostgreSQL and Redis only
docker-compose up -d postgres redis

# Run migrations
make migrate-up

# Start with hot reload
make dev
```

The API will be available at `http://localhost:3000`

## 📚 Project Structure

```
.
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connections and migrations
│   ├── middleware/       # HTTP middleware
│   ├── response/         # Response utilities (JSend)
│   ├── validator/        # Input validation
│   └── domains/          # Business domains
│       └── user/         # Reference domain (User/Auth)
│           ├── entity/   # Domain entities
│           ├── repository/ # Repository interfaces
│           ├── usecase/  # Business logic
│           └── handler/  # HTTP handlers
├── db/
│   ├── migrations/       # SQL migration files
│   └── queries/          # SQLC query files
├── docs/                 # Documentation and guides
├── configs/              # Configuration files
└── scripts/              # Utility scripts
```

## 🛠️ Development

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make dev               # Run with hot reload (Air)
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make lint              # Run linter
make format            # Format code
make docs              # Generate Swagger documentation
make migrate-up        # Run database migrations
make migrate-down      # Rollback last migration
make migrate-create    # Create new migration (NAME=migration_name)
make sqlc              # Generate Go code from SQL
make generate-mocks    # Generate mock implementations
make docker-build      # Build Docker image
make docker-up         # Start Docker Compose services
make docker-down       # Stop Docker Compose services
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/config
```

### Database Migrations

```bash
# Create new migration
make migrate-create NAME=add_users_table

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down
```

### Generating Code

```bash
# Generate type-safe Go code from SQL queries
make sqlc

# Generate Swagger documentation
make docs

# Generate mocks for testing
make generate-mocks
```

## 🔧 Adding a New Domain

See [docs/ADDING_NEW_DOMAIN.md](docs/ADDING_NEW_DOMAIN.md) for detailed instructions on adding new business domains.

Quick overview:
1. Create domain structure in `internal/domains/your-domain/`
2. Create migration files in `db/migrations/`
3. Create SQL queries in `db/queries/`
4. Run `make sqlc` to generate repository code
5. Implement usecase, handler, and tests
6. Register routes in router

## 🚢 Deployment

### Docker Deployment

```bash
# Build image
make docker-build

# Run container
docker run -p 3000:3000 \
  -e DB_HOST=your-db-host \
  -e JWT_SECRET=your-secret \
  gofiber-skeleton:latest
```

### Environment Variables

See `.env.example` for all available configuration options.

Critical production settings:
- `APP_ENV=production`
- `JWT_SECRET` - Strong secret key
- `DB_*` - Production database credentials
- `REDIS_*` - Production Redis credentials

## 📖 API Documentation

When the server is running, access interactive API documentation at:
- Swagger UI: `http://localhost:3000/swagger/`

Generate documentation:
```bash
make docs
```

## 🔐 Health Checks

The template includes built-in health check endpoints:

- `/health` - Overall health status with database stats
- `/health/ready` - Kubernetes readiness probe
- `/health/live` - Kubernetes liveness probe

## 🧪 Testing

The template follows these testing principles:
- **Unit Tests** - Business logic with 90%+ coverage
- **Mock-based** - Isolated tests with go.uber.org/mock
- **Integration Tests** - Database-backed tests
- **Test Fixtures** - Reusable test data

Example:
```go
func TestUserUsecase_Register(t *testing.T) {
    // See internal/domains/user/usecase/auth_test.go
}
```

## 📝 Code Quality

### Linting

```bash
# Run linter
make lint

# Auto-fix issues
golangci-lint run --fix ./...
```

### Formatting

```bash
# Format code
make format
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [sqlc](https://sqlc.dev/) - Type-safe SQL
- [Viper](https://github.com/spf13/viper) - Configuration
- [Samber's do](https://github.com/samber/do) - Dependency injection

## 📧 Support

- GitHub Issues: https://github.com/zercle/gofiber-skeleton/issues
- Documentation: [docs/](docs/)

---

**Made with ❤️ for the Go community**
