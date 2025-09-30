# Current Project Context: Production-Ready Go Fiber Backend Template

This document captures the current state of the Go Fiber backend template following production hardening.

## 1. Project Status: Production-Ready ✅

The template has been fully migrated from GORM to sqlc and hardened for production deployment with enterprise-grade features.

## 2. Completed Implementation

### Core Architecture
- ✅ **Feature-Based Structure**: User and post domains fully implemented under `internal/`
- ✅ **Type-Safe SQL**: Complete migration from GORM to sqlc with generated type-safe queries
- ✅ **Clean Architecture**: Strict separation of handlers, usecases, repositories, and entities

### Authentication & Authorization
- ✅ **User Registration & Login**: Bcrypt password hashing + JWT tokens
- ✅ **Protected Routes**: JWT middleware for authenticated endpoints
- ✅ **Rate Limiting**: Aggressive rate limiting on auth endpoints (5 req/min) to prevent brute force attacks

### API Features
- ✅ **JSend Responses**: All endpoints return standardized JSend-compliant JSON responses
- ✅ **Input Validation**: go-playground/validator with clear error messages
- ✅ **Swagger Documentation**: Comprehensive OpenAPI docs at `/swagger/*` endpoints
- ✅ **CRUD Operations**: Full post management with ownership validation

### Production Middleware
- ✅ **Request ID Tracking**: UUID-based request tracing for distributed systems
- ✅ **Structured Logging**: Zerolog with request context (method, path, status, duration, user_agent)
- ✅ **Rate Limiting**: 100 req/min for API, 5 req/min for auth endpoints
- ✅ **Panic Recovery**: Graceful error handling and logging
- ✅ **Graceful Shutdown**: SIGTERM/SIGINT handling with 30s timeout

### Infrastructure
- ✅ **PostgreSQL**: Type-safe queries via sqlc
- ✅ **Valkey/Redis**: Cache layer configured in docker-compose
- ✅ **Health Checks**: `/health` (liveness) and `/ready` (readiness with DB ping)
- ✅ **Docker Support**: Multi-stage Dockerfile with non-root user
- ✅ **CI/CD**: GitHub Actions workflow for testing, linting, and building

## 3. Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Framework | Fiber v2 | High-performance HTTP server |
| Database | PostgreSQL 18 | Relational data store |
| Cache | Valkey 8 | Redis-compatible cache |
| SQL Toolkit | sqlc | Type-safe SQL query generation |
| Auth | JWT + Bcrypt | Stateless authentication |
| Validation | go-playground/validator/v10 | Request validation |
| Logging | zerolog | Structured logging |
| API Docs | swaggo/swag | OpenAPI/Swagger generation |

## 4. Current Development Focus

**Status**: Ready for deployment and extension with new features.

### Immediate Priorities
- Deploy to staging environment for integration testing
- Add integration tests for auth and post flows
- Implement caching layer using Valkey/Redis
- Add observability (metrics, tracing) for production monitoring

### Future Enhancements
- Thread management feature (currently posts reference threads but CRUD not implemented)
- Comment system for posts
- User profile management
- Role-based access control (RBAC)
- Email verification for user registration
- Password reset functionality
- File upload support for posts
- Full-text search for posts

## 5. Key Files & Locations

```
├── cmd/server/main.go              # Application entrypoint with Swagger annotations
├── internal/
│   ├── config/                     # Configuration management
│   ├── db/                         # sqlc generated code
│   ├── middleware/                 # Production middleware (request ID, logging, rate limiting)
│   ├── response/                   # JSend response utilities
│   ├── user/                       # User feature (auth, registration)
│   └── post/                       # Post feature (CRUD operations)
├── pkg/validator/                  # Reusable validation utilities
├── db/
│   ├── migrations/                 # SQL schema migrations
│   └── queries/                    # sqlc query definitions
├── docs/                           # Generated Swagger documentation
├── docker-compose.yml              # Local development environment (app, db, redis)
└── Dockerfile                      # Production-ready multi-stage build
```

## 6. Quick Start Commands

```bash
# Generate code
make sqlc                    # Generate sqlc code
make generate-docs           # Generate Swagger docs
go generate ./...            # Generate mocks

# Development
docker-compose up -d         # Start services
go run cmd/server/main.go    # Run server
make ci                      # Run full CI pipeline locally

# Testing
go test -v -race ./...       # Run all tests with race detector

# Deployment
docker-compose up --build    # Build and run containerized app
```

## 7. API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token

### Posts (Protected)
- `POST /api/v1/posts` - Create post (requires auth)
- `GET /api/v1/posts/:id` - Get post by ID
- `GET /api/v1/users/:user_id/posts` - List user's posts
- `PUT /api/v1/posts/:id` - Update post (requires auth + ownership)
- `DELETE /api/v1/posts/:id` - Delete post (requires auth + ownership)

### Health
- `GET /health` - Liveness check
- `GET /ready` - Readiness check (includes DB ping)