# Template Setup Guide

This guide will help you set up a new project from this Go Fiber template.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.25.0 or higher
- **Docker**: For running PostgreSQL and Valkey/Redis
- **Docker Compose**: For orchestrating containers
- **Make**: For running convenience commands
- **Git**: For version control

## Initial Setup Steps

### 1. Clone or Use Template

**Option A: Use GitHub Template**
1. Click "Use this template" button on GitHub
2. Create your new repository
3. Clone your repository:
   ```bash
   git clone https://github.com/your-username/your-project.git
   cd your-project
   ```

**Option B: Clone Directly**
```bash
git clone https://github.com/zercle/gofiber-skeleton.git your-project
cd your-project
rm -rf .git
git init
```

### 2. Update Module Name

Replace the module name throughout the project:

```bash
# Find all files with the old module name
grep -r "github.com/zercle/gofiber-skeleton" .

# Use your editor's find-and-replace or run:
# macOS/BSD:
find . -type f -name "*.go" -exec sed -i '' 's|github.com/zercle/gofiber-skeleton|github.com/your-username/your-project|g' {} +

# Linux:
find . -type f -name "*.go" -exec sed -i 's|github.com/zercle/gofiber-skeleton|github.com/your-username/your-project|g' {} +
```

Update `go.mod`:
```go
module github.com/your-username/your-project

go 1.25
```

### 3. Configure Environment

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your project-specific values:
```bash
# Update these critical values
DB_NAME=your_project_db
JWT_SECRET=your-super-secret-jwt-key-CHANGE-THIS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://your-domain.com

# Optional: Update other values as needed
SERVER_PORT=3000
LOG_LEVEL=info
```

### 4. Install Development Tools

Install required development tools:
```bash
make install-tools
```

This installs:
- **air**: Hot reload development server
- **sqlc**: SQL code generator
- **swag**: API documentation generator
- **golangci-lint**: Code quality tool
- **mockgen**: Mock generator for testing
- **migrate**: Database migration tool

### 5. Initialize Dependencies

Download Go module dependencies:
```bash
go mod download
go mod tidy
```

### 6. Start Development Environment

Start PostgreSQL and Valkey containers:
```bash
make docker-up
```

Verify containers are running:
```bash
docker ps
```

You should see:
- `gofiber_postgres` on port 5432
- `gofiber_valkey` on port 6379

### 7. Run Database Migrations

Apply the initial database schema:
```bash
make migrate-up
```

Verify migration status:
```bash
make migrate-version
```

### 8. Generate Code

Generate SQL code and API documentation:
```bash
# Generate type-safe SQL code
make sqlc

# Generate API documentation
make swag

# Generate test mocks (optional, will be auto-generated when needed)
make mocks
```

### 9. Start Development Server

Start the server with hot reload:
```bash
make dev
```

The server will start on http://localhost:3000

### 10. Verify Installation

Test the installation:

**Health Check:**
```bash
curl http://localhost:3000/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "go-fiber-skeleton",
  "database": "healthy",
  "cache": "healthy",
  "time": "2025-10-11T..."
}
```

**API Documentation:**

Visit http://localhost:3000/swagger/ in your browser

**Register a Test User:**
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123",
    "full_name": "Test User"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123"
  }'
```

## Project Customization

### Rename Application

Update application name in these files:

1. **README.md**: Replace project title and descriptions
2. **cmd/server/main.go**: Update `@title` and `AppName` in config
3. **Docker files**: Update image names and labels
4. **GitHub workflows**: Update repository references

### Remove Example Domain (Optional)

If you want to start fresh without the user domain:

```bash
# Remove user domain
rm -rf internal/domains/user

# Remove user migrations
rm db/migrations/000001_create_users_table.*

# Remove user queries
rm db/queries/users.sql

# Update main.go to remove user routes
```

Then add your own domains following the same structure.

### Update Database Schema

Add your own migrations:
```bash
make migrate-create NAME=create_your_table
```

This creates:
- `db/migrations/XXXXXX_create_your_table.up.sql`
- `db/migrations/XXXXXX_create_your_table.down.sql`

Write your SQL, then apply:
```bash
make migrate-up
```

### Add SQL Queries

Create query files in `db/queries/`:
```sql
-- name: CreateYourEntity :one
INSERT INTO your_table (column1, column2)
VALUES ($1, $2)
RETURNING *;

-- name: GetYourEntityByID :one
SELECT * FROM your_table
WHERE id = $1 LIMIT 1;
```

Generate Go code:
```bash
make sqlc
```

## Development Workflow

### Daily Development

1. Start containers: `make docker-up`
2. Start dev server: `make dev`
3. Make code changes (auto-reload)
4. Run tests: `make test`
5. Check code quality: `make lint`

### Before Committing

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Build to verify
make build
```

### Adding New Features

1. Create migration if needed
2. Add SQL queries
3. Generate SQL code: `make sqlc`
4. Implement domain layers (entity → repository → usecase → delivery)
5. Write tests
6. Update API docs: `make swag`
7. Test endpoints

## Docker Deployment

### Build Production Image

```bash
make docker-build
```

### Run Production Container

```bash
docker run -d \
  --name your-project \
  -p 3000:3000 \
  --env-file .env.production \
  your-project:latest
```

### Docker Compose Production

Create `docker-compose.prod.yml`:
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - SERVER_ENV=production
      - DB_HOST=postgres
      - VALKEY_HOST=valkey
    depends_on:
      - postgres
      - valkey

  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_DB: your_project_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data

  valkey:
    image: valkey/valkey:latest
    volumes:
      - valkey_data:/data

volumes:
  postgres_data:
  valkey_data:
```

Run:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 3000
lsof -i :3000

# Kill process
kill -9 <PID>
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs gofiber_postgres

# Restart container
make docker-down
make docker-up
```

### Migration Failed

```bash
# Check migration status
make migrate-version

# Force to specific version
make migrate-force VERSION=1

# Rollback and retry
make migrate-down
make migrate-up
```

### Module Issues

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

### Hot Reload Not Working

```bash
# Check Air configuration
cat .air.toml

# Restart Air
# Press Ctrl+C and run: make dev
```

## Next Steps

1. **Read Architecture Docs**: `docs/ADDING_NEW_DOMAIN.md`
2. **Review User Domain**: `internal/domains/user/`
3. **Add Your First Domain**: Follow the user domain pattern
4. **Configure CI/CD**: Update `.github/workflows/`
5. **Set Up Production**: Configure deployment environment
6. **Add Monitoring**: Integrate logging and metrics

## Support

- **Issues**: https://github.com/zercle/gofiber-skeleton/issues
- **Discussions**: https://github.com/zercle/gofiber-skeleton/discussions
- **Documentation**: See `docs/` directory

## Checklist

Before deploying to production:

- [ ] Changed JWT secret from default
- [ ] Updated CORS allowed origins
- [ ] Configured production database
- [ ] Set up SSL/TLS certificates
- [ ] Configured logging and monitoring
- [ ] Set up backup strategy
- [ ] Configured rate limiting
- [ ] Reviewed security settings
- [ ] Added health check monitoring
- [ ] Set up CI/CD pipeline
- [ ] Documented API endpoints
- [ ] Added integration tests
- [ ] Configured error tracking
- [ ] Set up performance monitoring
