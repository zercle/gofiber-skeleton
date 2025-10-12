# Go Fiber Skeleton

A production-ready Go Fiber template implementing Clean Architecture with comprehensive features for rapid backend development.

## 🚀 Features

- **Clean Architecture**: Domain-driven design with strict separation of concerns
- **Modern Go Stack**: Fiber v2, SQLC, JWT, Argon2id, Docker
- **Database Integration**: PostgreSQL with migrations and type-safe queries
- **Authentication**: JWT-based auth with secure password hashing
- **API Documentation**: Auto-generated Swagger/OpenAPI docs
- **Testing**: Comprehensive test coverage with mocks
- **Development Tools**: Hot reload, linting, CI/CD ready
- **Production Ready**: Containerized, health checks, monitoring

## 📁 Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── server/            # Main HTTP server
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── domains/          # Business domains
│   │   └── user/         # User/auth domain
│   │       ├── entity/   # Domain entities
│   │       ├── repository/ # Data access interfaces
│   │       ├── usecase/  # Business logic
│   │       ├── delivery/ # HTTP handlers
│   │       └── tests/    # Domain tests
│   └── middleware/       # HTTP middleware
├── pkg/                  # Shared library code
│   ├── database/        # Database utilities
│   └── response/        # Response formatting
├── db/                  # Database-related files
│   ├── migrations/      # SQL migration files
│   ├── queries/         # SQLC query files
│   └── seeds/           # Database seeds
├── docs/                # Generated documentation
├── configs/             # Configuration files
└── scripts/             # Utility scripts
```

## 🛠️ Quick Start

### Prerequisites

- Go 1.25.0+
- Docker & Docker Compose
- PostgreSQL 18+
- Make (optional, but recommended)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd gofiber-skeleton
   ```

2. **Setup development environment**
   ```bash
   make quick-start
   ```

   This will:
   - Install Go dependencies
   - Start PostgreSQL and Valkey containers
   - Run database migrations
   - Generate SQLC code
   - Generate Swagger docs
   - Generate test mocks

3. **Start development server**
   ```bash
   make dev
   ```

   The server will be available at `http://localhost:3000`

### Manual Setup

If you prefer manual setup:

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Start services**
   ```bash
   docker-compose up -d
   ```

3. **Run migrations**
   ```bash
   make migrate-up
   ```

4. **Generate code**
   ```bash
   make sqlc swag mocks
   ```

5. **Run server**
   ```bash
   go run cmd/server/main.go
   ```

## 📚 API Documentation

- **Swagger UI**: `http://localhost:3000/swagger/`
- **Health Check**: `http://localhost:3000/health`

## 🔧 Available Commands

### Development
- `make dev` - Start development server with hot reload
- `make run` - Run server without hot reload
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make lint` - Run linter
- `make fmt` - Format Go code

### Database
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create NAME=migration_name` - Create new migration
- `make sqlc` - Generate SQLC code

### Code Generation
- `make swag` - Generate Swagger documentation
- `make mocks` - Generate test mocks

### Build & Deploy
- `make build` - Build for production
- `make build-all` - Build for multiple platforms
- `make docker-build` - Build Docker image

## 🏗️ Adding New Domains

Follow the user domain pattern:

1. **Create domain structure**
   ```bash
   mkdir -p internal/domains/yourdomain/{entity,repository,usecase,delivery,tests,mocks}
   ```

2. **Implement domain layers**
   - `entity/` - Domain models and business rules
   - `repository/` - Data access interfaces
   - `usecase/` - Business logic
   - `delivery/` - HTTP handlers

3. **Add database queries**
   - Create SQL files in `db/queries/yourdomain.sql`
   - Run `make sqlc` to generate code

4. **Register dependencies**
   - Add providers in `cmd/server/main.go`
   - Add routes in `setupRoutes()`

5. **Add tests**
   - Create unit tests with mocks
   - Test business logic thoroughly

## 🧪 Testing

### Running Tests
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests
make test-integration
```

### Testing Strategy
- **Unit Tests**: Test business logic in isolation with mocks
- **Integration Tests**: Test database operations and external integrations
- **API Tests**: Test HTTP endpoints
- **Coverage Target**: 90%+

## 🔐 Authentication

The template includes a complete authentication system:

- **JWT Tokens**: Stateless authentication with configurable expiry
- **Password Hashing**: Argon2id for secure password storage
- **Middleware**: Easy route protection
- **User Management**: Registration, login, profile management

### Example Usage
```go
// Protected route
protected := app.Group("/users")
protected.Use(middleware.JWTMiddleware(config.JWT.Secret))
protected.Get("/profile", userHandler.GetProfile)
```

## 📊 Configuration

Configuration is managed through environment variables and `.env` files:

```bash
# Application
APP_NAME=gofiber-skeleton
APP_ENV=development
APP_PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gofiber_skeleton

# JWT
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY=24h
```

See `.env.example` for all available options.

## 🐳 Docker Development

### Development with Hot Reload
```bash
docker-compose --profile dev up
```

### Production Build
```bash
docker build -t gofiber-skeleton .
docker run -p 3000:3000 gofiber-skeleton
```

## 🚀 Deployment

### Environment Setup
1. Set production environment variables
2. Configure `compose.yml` for production services
3. Build and deploy containers

### Health Checks
- Application: `/health`
- Database: Built-in connection health check
- Ready for container orchestration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make ci` to ensure quality
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Fiber Documentation](https://docs.gofiber.io/)
- [SQLC Documentation](https://sqlc.dev/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go.html)