# E-commerce Management System Backend Boilerplate 🛍️

A production-ready, scalable backend boilerplate for e-commerce management systems built with Go, following Clean Architecture principles and SOLID design patterns.

## ✨ Features

- **🔐 Authentication & Authorization**: JWT-based authentication with role-based access control
- **📦 Product Management**: Full CRUD operations for products with stock management
- **🛒 Order Management**: Complete order lifecycle with status tracking
- **👥 User Management**: User registration, login, and role management
- **🏗️ Clean Architecture**: Well-structured, maintainable codebase
- **🧪 Testing Ready**: Mock interfaces and testing utilities included
- **🐳 Docker Support**: Containerized application with PostgreSQL
- **📊 Database**: PostgreSQL with SQLC for type-safe queries
- **🔄 Migrations**: Database schema management with golang-migrate

## 🏗️ Architecture

The project follows Clean Architecture principles with clear separation of concerns:

```
cmd/server/          # Application entry point
internal/
├── domain/          # Business entities and interfaces
├── handler/         # HTTP request handlers
├── usecase/         # Business logic layer
├── repository/      # Data access layer
└── infrastructure/  # Database and middleware setup
pkg/                 # Shared utilities
migrations/          # Database migrations
queries/             # SQLC query definitions
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.6 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

### Option 1: Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd gofiber-skeleton
   ```

2. **Start the application**
   ```bash
   docker compose up --build
   ```

3. **Access the API**
   - API: http://localhost:8080
   - Health check: http://localhost:8080/health

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Start the application**
   ```bash
   make run
   ```

## 📚 API Endpoints

### Public Endpoints
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User authentication
- `GET /api/v1/products` - List all products
- `GET /api/v1/products/{id}` - Get product details

### Protected Endpoints (Require JWT)
- `POST /api/v1/products` - Create product (Admin only)
- `PUT /api/v1/products/{id}` - Update product (Admin only)
- `DELETE /api/v1/products/{id}` - Delete product (Admin only)
- `POST /api/v1/orders/create` - Create order
- `GET /api/v1/orders` - Get user orders
- `GET /api/v1/orders/{id}` - Get order details
- `GET /api/v1/orders/admin/all` - Get all orders (Admin only)
- `PUT /api/v1/orders/{id}/status` - Update order status (Admin only)

## 🛠️ Development

### Available Commands

```bash
make help              # Show available commands
make build             # Build the application
make run               # Run locally
make test              # Run tests
make generate-mocks    # Generate mock files
make docker-build      # Build Docker image
make docker-run        # Run with Docker Compose
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make lint              # Run linter
make fmt               # Format code
```

### Code Generation

The project uses several code generation tools:

- **SQLC**: Generates Go code from SQL queries
- **Go Mock**: Generates mock interfaces for testing

```bash
# Generate SQLC code
make sqlc-generate

# Generate mocks
make generate-mocks
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `ecommerce` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `PORT` | Server port | `8080` |

### Database Schema

The application includes the following tables:
- `users` - User accounts and authentication
- `products` - Product catalog with inventory
- `orders` - Customer orders
- `order_items` - Individual items within orders

## 🐳 Docker

### Building the Image

```bash
make docker-build
```

### Running with Docker Compose

```bash
make docker-run
```

### Environment Variables in Docker

The Docker Compose setup includes:
- PostgreSQL 15 with persistent storage
- Health checks for database readiness
- Automatic migration execution
- Network isolation

## 📝 Database Migrations

Migrations are managed using `golang-migrate`:

```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## 🧪 Testing Strategy

The project includes:
- **Unit Tests**: Testing individual components with mocks
- **Integration Tests**: Testing component interactions
- **Mock Generation**: Automatic mock creation for interfaces
- **Test Coverage**: Comprehensive testing of business logic

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt password encryption
- **Role-Based Access Control**: Admin and customer roles
- **Input Validation**: Request payload validation
- **SQL Injection Prevention**: Parameterized queries via SQLC

## 📈 Performance Features

- **Fiber Framework**: High-performance HTTP framework
- **Connection Pooling**: Database connection management
- **Efficient Queries**: SQLC-generated optimized queries
- **Middleware Pipeline**: Optimized request processing

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the code examples

## 🔮 Roadmap

- [ ] GraphQL API support
- [ ] Payment integration
- [ ] Email notifications
- [ ] Advanced search and filtering
- [ ] Analytics and reporting
- [ ] Multi-tenant support
- [ ] API rate limiting
- [ ] Caching layer
- [ ] WebSocket support for real-time updates
