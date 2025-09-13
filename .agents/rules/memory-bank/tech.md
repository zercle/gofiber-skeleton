# Tech Stack and Environment Setup

## Primary Technologies
- Go Fiber v2 (github.com/gofiber/fiber/v2)
- Fx (Uber’s dependency injection framework)
- Viper for configuration (YAML, .env, env vars)
- PostgreSQL with pgx and sqlc for type-safe SQL
- golang-migrate/migrate for schema migrations
- go-playground/validator for request and model validation
- github.com/omniti-labs/jsend for standardized JSON responses
- github.com/golang-jwt/jwt/v5 for JWT authentication
- UUIDv7 for domain entity identifiers
- uber-go/mock and go-sqlmock for mocking in tests
- Go testing and Testify for unit and integration tests
- Air for hot reload during development
- swaggo/swag for API documentation generation

## Tooling & Commands
- Docker and Docker Compose for containerization and local services
- golangci-lint and gofmt for linting and formatting
- sqlc for generating Go code from SQL queries
- swag for generating Swagger docs
- git for version control

## Environment Variables
See `.env.example` for the full list. Key variables include:
- PORT, ENV
- DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SCHEMA, DB_SSLMODE, DB_URL
- JWT_SECRET, JWT_EXPIRES_IN
- VALKEY_HOST, VALKEY_PORT, VALKEY_PASSWORD, VALKEY_DB, VALKEY_URL
- CORS_ORIGINS

## Development Commands
```bash
go mod tidy
cp .env.example .env
docker-compose up -d
go run cmd/migrate/main.go
air