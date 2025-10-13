# Go Fiber Skeleton

Production-ready Go backend template built with **Fiber v2** implementing **Domain-Driven Clean Architecture**. Eliminates **80-90% of initial project setup** work, enabling developers to focus on business logic from day one.

## ✨ Features

- 🏗️ **Clean Architecture** with strict domain isolation
- 🔐 **JWT Authentication** with Argon2id password hashing
- 🗄️ **PostgreSQL + sqlc** for type-safe database operations
- 🚀 **High Performance** built on Fiber v2 web framework
- 📚 **Auto-generated API docs** with Swagger/OpenAPI
- 🐳 **Docker** ready for containerization
- 🧪 **Comprehensive testing** with mock generation
- 📦 **Type-safe DI** with Samber's do framework
- 🔧 **Development tools** pre-configured (Air, golangci-lint)
- 📊 **Structured logging** with correlation IDs

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (or use Docker setup)

### 5-Minute Setup

```bash
# Clone the repository
git clone https://github.com/zercle/gofiber-skeleton.git
cd gofiber-skeleton

# Start development environment
make setup
make docker-up

# Run migrations
make migrate-up

# Start the server (with hot-reload)
make dev
```

That's it! Your production-ready Go backend is running at `http://localhost:8080`

- 🌐 **API**: http://localhost:8080
- 📚 **Swagger Docs**: http://localhost:8080/swagger/index.html
- 🗄️ **Adminer**: http://localhost:8081 (Database GUI)

## 📁 Project Structure

```
gofiber-skeleton/
├── cmd/server/                 # Application entry point
├── internal/
│   ├── domains/               # Business domains (DDD)
│   │   └── user/              # Reference domain (entity, repository, usecase, handler)
│   ├── infrastructure/        # External implementations
│   │   ├── database/         # Database setup and sqlc code
│   │   ├── middleware/       # HTTP middleware
│   │   └── config/           # Configuration management
│   └── shared/               # Shared utilities
├── db/
│   ├── migrations/           # Database migrations
│   └── queries/              # SQL queries for sqlc
├── docs/                     # Documentation
├── compose.yml              # Docker Compose for development
├── Dockerfile               # Production container
└── Makefile                 # Development commands
```

## 🛠️ Available Commands

```bash
# Setup
make setup          # Install dependencies and tools
make sqlc           # Generate Go code from SQL queries
make swagger        # Generate API documentation
make mocks          # Generate mocks for testing

# Development
make dev            # Run with hot-reload
make run            # Run without hot-reload
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linter

# Database
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations
make migrate-create NAME=migration_name

# Docker
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make docker-logs    # Show logs

# Production
make build          # Build application
make docker-build   # Build Docker image
```

## 🏗️ Architecture Overview

### Clean Architecture Layers

1. **Handler Layer** - HTTP endpoints and request/response handling
2. **Usecase Layer** - Business logic and domain rules
3. **Repository Layer** - Data access interfaces
4. **Entity Layer** - Core domain models

### Reference Implementation

The **user domain** demonstrates all architectural patterns:

- **Entity**: User model with validation
- **Repository**: PostgreSQL integration with sqlc
- **Usecase**: Registration/login business logic
- **Handler**: HTTP endpoints with Swagger docs

## 📚 API Documentation

### Authentication Endpoints

```bash
# Register new user
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe"
}

# Login
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

### User Management Endpoints

```bash
# Get profile (requires auth)
GET /api/v1/users/profile
Authorization: Bearer <token>

# Update profile (requires auth)
PUT /api/v1/users/profile
{
  "full_name": "John Updated"
}

# Delete account (requires auth)
DELETE /api/v1/users/profile

# Get user by ID (requires auth)
GET /api/v1/users/{id}

# List users (requires auth)
GET /api/v1/users?page=1&limit=10
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/domains/user/usecase -v
```

## 🔧 Configuration

Configuration is managed through:

1. **Environment variables** (highest priority)
2. **`.env` file** (development)
3. **`config.yaml`** (default values)

See `.env.example` for all available configuration options.

## 🚀 Adding New Domains

Follow the user domain pattern:

1. **Create domain structure**: `internal/domains/{domain}/`
2. **Define entities**: `entity/`
3. **Create repository interfaces**: `repository/`
4. **Implement business logic**: `usecase/`
5. **Add HTTP handlers**: `handler/`
6. **Register routes** in `cmd/server/main.go`

Example for a `product` domain:

```bash
mkdir -p internal/domains/product/{entity,repository,usecase,handler}
# Follow the user domain implementation pattern
```

## 🐳 Docker Deployment

### Development

```bash
docker-compose up -d
```

### Production

```bash
# Build image
make docker-build

# Run container
docker run -p 8080:8080 \
  -e GS_DATABASE_HOST=your-db-host \
  -e GS_JWT_SECRET=your-secret \
  gofiber-skeleton
```

## 📊 Production Considerations

### Security

- ✅ JWT authentication with secure defaults
- ✅ Password hashing with Argon2id
- ✅ Security headers (HSTS, CSP, X-Frame-Options)
- ✅ Rate limiting and CORS configuration
- ✅ Input validation and sanitization

### Performance

- ✅ Connection pooling for database
- ✅ Structured logging with correlation IDs
- ✅ Graceful shutdown handling
- ✅ Health check endpoints
- ✅ Request timeout configuration

### Observability

- ✅ Structured JSON logging
- ✅ Request ID tracking
- ✅ Health checks for dependencies
- ✅ Performance metrics ready

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test lint`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - High-performance web framework
- [sqlc](https://github.com/sqlc-dev/sqlc) - Type-safe SQL builder
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin

## 📞 Support

- 📖 [Documentation](docs/)
- 🐛 [Issues](https://github.com/zercle/gofiber-skeleton/issues)
- 💬 [Discussions](https://github.com/zercle/gofiber-skeleton/discussions)

---

**Built with ❤️ for Go developers**