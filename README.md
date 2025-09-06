# Go Fiber Backend Mono-Repo Template

A production-ready **Go Fiber Backend Template** implementing **Domain-Driven Clean Architecture** with a mono-repo structure. This template provides a scalable foundation for building REST APIs with modern Go development practices.

## ✨ Features

- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Domain Isolation**: Each domain is self-contained with complete layers
- **Production Ready**: Comprehensive error handling, logging, and monitoring
- **Modern Stack**: Latest Go tools and frameworks integrated
- **Developer Experience**: Hot reloading, automated testing, and easy deployment
- **API Documentation**: Auto-generated Swagger documentation
- **Database Migrations**: Version-controlled database schema management
- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Request validation with detailed error messages
- **CORS Support**: Configurable cross-origin resource sharing
- **Health Checks**: Built-in health monitoring endpoints
- **Docker Support**: Complete containerization with Docker Compose

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────┐
│           HTTP Handlers             │  ← Delivery Layer
├─────────────────────────────────────┤
│            Use Cases                │  ← Business Logic
├─────────────────────────────────────┤
│      Entities & Repositories        │  ← Domain Layer
├─────────────────────────────────────┤
│      Database & External Services   │  ← Infrastructure
└─────────────────────────────────────┘
```

### Project Structure

```
.
├── cmd/                          # Application entry points
│   ├── server/main.go           # Main server application
│   └── migrate/main.go          # Database migration runner
├── internal/                    # Private application code
│   ├── domains/                 # Business domains
│   │   ├── auth/               # Authentication domain
│   │   │   ├── entities/       # Core business entities
│   │   │   ├── usecases/       # Business logic
│   │   │   ├── repositories/   # Data access interfaces
│   │   │   ├── handlers/       # HTTP handlers
│   │   │   ├── routes/         # Route definitions
│   │   │   ├── models/         # DTOs and request/response models
│   │   │   └── tests/          # Domain-specific tests
│   │   └── posts/              # Posts domain example
│   ├── infrastructure/         # External concerns
│   │   ├── database/           # Database connections
│   │   ├── middleware/         # HTTP middleware
│   │   └── config/             # Configuration management
│   └── shared/                 # Shared components
│       ├── types/              # Common types and interfaces
│       └── container/          # Dependency injection
├── pkg/                        # Public utilities
│   └── utils/                  # Utility functions
├── migrations/                 # Database migrations
├── docs/                       # API documentation
├── compose.yml          # Development environment
└── Dockerfile                  # Production container
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Redis (optional)
- Docker & Docker Compose (optional)

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd gofiber-skeleton

# Copy environment file
cp .env.example .env

# Install dependencies
make deps

# Install development tools
make install-tools
```

### 2. Database Setup

```bash
# Start PostgreSQL with Docker
make db-up

# Or run migrations against existing PostgreSQL
make migrate-up
```

### 3. Run the Application

```bash
# Development mode with hot reloading
make dev

# Or build and run
make build
make run
```

### 4. Access the API

- **API Base URL**: http://localhost:3000/api/v1
- **Health Check**: http://localhost:3000/health
- **API Documentation**: http://localhost:3000/docs (development only)

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | User login |
| GET | `/auth/profile` | Get user profile |
| PUT | `/auth/profile` | Update user profile |
| POST | `/auth/change-password` | Change password |

### Posts Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/posts` | List posts (public) |
| POST | `/posts` | Create new post |
| GET | `/posts/{id}` | Get post by ID |
| PUT | `/posts/{id}` | Update post |
| DELETE | `/posts/{id}` | Delete post |
| GET | `/posts/me` | List user's posts |
| POST | `/posts/{id}/publish` | Publish post |
| POST | `/posts/{id}/unpublish` | Unpublish post |

### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/users` | List all users |
| GET | `/admin/users/{id}` | Get user by ID |
| POST | `/admin/users/{id}/activate` | Activate user |
| POST | `/admin/users/{id}/deactivate` | Deactivate user |
| DELETE | `/admin/users/{id}` | Delete user |

## 🛠️ Development

### Available Commands

```bash
make help              # Show all available commands

# Development
make dev               # Run with hot reloading
make build             # Build application
make run               # Run application
make clean             # Clean build artifacts

# Testing
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Lint code
make fmt               # Format code

# Database
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make migrate-create    # Create new migration

# Documentation
make docs              # Generate API docs

# Docker
make docker-build      # Build Docker image
make docker-run        # Run with Docker Compose
```

### Adding a New Domain

1. **Create directory structure:**
```bash
mkdir -p internal/domains/{domain}/{entities,usecases,repositories,handlers,routes,models,tests}
```

2. **Follow the pattern:**
   - Define entities with business rules
   - Create repository interfaces
   - Implement use cases
   - Build HTTP handlers and routes
   - Register dependencies in DI container
   - Write comprehensive tests

3. **Register in main application:**
```go
// Add to cmd/server/main.go
domainRoutes.SetupDomainRoutes(app, domainHandler, cfg)
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific domain
go test ./internal/domains/auth/...

# Generate mocks
make generate
```

## 🐳 Docker Deployment

### Development

```bash
# Start all services
compose up

# Start in background
compose up -d

# View logs
compose logs -f app
```

### Production

```bash
# Build production image
docker build -t gofiber-skeleton .

# Run production container
docker run -p 3000:3000 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="your-secret" \
  gofiber-skeleton
```

## ⚙️ Configuration

Configuration is managed through multiple sources (in order of precedence):

1. **Environment Variables**
2. **Config Files** (config.yaml, .env)
3. **Default Values**

### Environment Variables

```bash
# Application
PORT=3000
ENV=production
APP_NAME=gofiber-skeleton

# Database
DATABASE_URL=postgres://user:pass@host:5432/db

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Redis (optional)
REDIS_URL=redis://localhost:6379
```

## 📝 API Documentation

API documentation is automatically generated using Swagger annotations:

```bash
# Generate docs
make docs

# View docs (development mode)
# Visit: http://localhost:3000/docs
```

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Input Validation**: Comprehensive request validation
- **CORS Protection**: Configurable cross-origin policies
- **Error Handling**: Secure error responses without sensitive data
- **Rate Limiting**: (Ready for implementation)
- **Request Logging**: Comprehensive request/response logging

## 🚀 Deployment

### Environment Setup

1. **Production Environment Variables**
2. **Database Migration**
3. **Secret Management**
4. **Load Balancing**
5. **Monitoring & Logging**

### Recommended Deployment

- **Container Orchestration**: Kubernetes, Docker Swarm
- **Database**: Managed PostgreSQL (AWS RDS, Google Cloud SQL)
- **Caching**: Redis cluster
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack or similar

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Go Fiber](https://gofiber.io/) - Web framework
- [Uber FX](https://uber-go.github.io/fx/) - Dependency injection
- [Viper](https://github.com/spf13/viper) - Configuration management
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [sqlc](https://sqlc.dev/) - SQL code generation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit