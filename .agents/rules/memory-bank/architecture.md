# Architecture Overview

## System Architecture
- The project follows Domain-Driven Clean Architecture within a mono-repo.
- Each domain is a self-contained module under internal/domains with layers: entities, usecases, repositories, handlers, routes, models, tests.

## Codebase Structure
```
.
├── cmd/
│   ├── server/
│   │   └── main.go          # Application entry point
│   └── migrate/
│       └── main.go          # Database migration runner
├── internal/
│   ├── app/
│   │   └── app.go           # Fiber app configuration
│   ├── domains/             # Domain-specific modules
│   │   ├── auth/            # Authentication domain
│   │   │   ├── entities/
│   │   │   │   └── user.go
│   │   │   ├── usecases/
│   │   │   │   ├── auth_usecase.go
│   │   │   │   └── interfaces/
│   │   │   ├── repositories/
│   │   │   │   ├── user_repository.go
│   │   │   │   └── interfaces/
│   │   │   ├── handlers/
│   │   │   │   └── auth_handler.go
│   │   │   ├── routes/
│   │   │   │   └── auth_routes.go
│   │   │   ├── models/
│   │   │   │   └── auth_models.go
│   │   │   └── tests/
│   │   │       ├── auth_usecase_test.go
│   │   │       └── auth_handler_test.go
│   │   ├── posts/           # Posts domain
│   │   │   ├── entities/
│   │   │   │   └── post.go
│   │   │   ├── usecases/
│   │   │   ├── repositories/
│   │   │   ├── handlers/
│   │   │   ├── routes/
│   │   │   ├── models/
│   │   │   └── tests/
│   │   └── greeting/        # Greeting domain (example)
│   │       ├── entities/
│   │       ├── usecases/
│   │       ├── repositories/
│   │       ├── handlers/
│   │       ├── routes/
│   │       ├── models/
│   │       └── tests/
│   ├── infrastructure/      # Shared infrastructure
│   │   ├── database/
│   │   │   ├── connection.go
│   │   │   └── base_repository.go
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go
│   │   │   ├── cors_middleware.go
│   │   │   └── validation_middleware.go
│   │   └── config/
│   │       └── app_config.go
│   └── shared/             # Shared components
│       ├── types/
│       │   └── common.go
│       └── container/
│           └── di_container.go
├── pkg/
│   └── utils/              # Public utilities
│       ├── uuidv7.go       # UUIDv7 generator
│       ├── response.go
│       └── validation.go
├── migrations/             # Database migration files
├── docs/                   # Swagger documentation
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .air.toml               # Air configuration for hot reload
└── README.md
```

## Key Technical Decisions and Patterns
- HTTP Framework: Go Fiber v2  
- Dependency Injection: fx (Uber’s DI framework)  
- Database ORM: GORM with migrations via golang-migrate and SQL code generation via sqlc  
- Configuration Management: Viper  
- JSON Response Format: omniti-labs/jsend  
- Mock Generation: go.uber.org/mock/mockgen and go-sqlmock  
- Validation: go-playground/validator  
- Testing: Go testing + Testify  
- OpenAPI Documentation: swaggo/swag  
- Hot Reload: Air  
- Linting & Formatting: golangci-lint, gofmt

## Component Relationships
- fx Module wires together configuration, database, repositories, usecases, handlers, and app.  
- Database infrastructure under internal/infrastructure/database.  
- Shared utilities under internal/shared and pkg/utils.  
- Domain modules under internal/domains.  
- API documentation generated in docs/ via swaggo.

## Critical Code Flows
- Application startup: cmd/server/main.go → fx container → app.NewApp() → Fiber app instantiation and route registration.  
- Database migrations: cmd/migrate/main.go.  
- Running tests: go test ./…  
- Generating docs: swag init -g cmd/server/main.go -o docs.  