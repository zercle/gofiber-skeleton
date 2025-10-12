# Go Fiber Production-Ready Backend Template

A minimal, production-ready Go Fiber backend service template with essential features for building robust APIs. This template provides a solid foundation with clean architecture, structured logging, graceful shutdown, and comprehensive middleware.

## 🚀 Features

- **Web Framework**: [Go Fiber v2](https://github.com/gofiber/fiber) - High-performance web framework
- **Structured Logging**: Built-in `slog` with JSON/text format support
- **Configuration Management**: Viper with environment variable and .env file support
- **Graceful Shutdown**: Proper signal handling and timeout management
- **Middleware Stack**: Recovery, CORS, request logging, request ID
- **Health Checks**: Complete health, liveness, and readiness endpoints
- **Sample API**: CRUD operations for todo items with error handling
- **Production Ready**: Error handling, timeouts, and proper HTTP status codes

## 📁 Project Structure

```
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── handlers/         # HTTP request handlers
│   │   ├── api.go       # Sample API handlers
│   │   └── health.go    # Health check handlers
│   ├── middleware/       # HTTP middleware
│   │   └── middleware.go
│   └── services/         # Business logic services
├── pkg/                 # Public library code
│   ├── logger/          # Structured logging
│   │   └── logger.go
│   └── server/          # Server setup and management
│       └── server.go
├── .env.example         # Environment variables template
├── .env                 # Local environment variables
├── go.mod              # Go module definition
├── package.json        # Development scripts
└── README.md           # This file
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.25 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/zercle/gofiber-skeleton.git
   cd gofiber-skeleton
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

The server will start at `http://localhost:3000`

### Development Workflow

For development with hot-reload:

1. **Install Air (for hot reloading)**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **Run in development mode**
   ```bash
   npm run dev
   # or
   air
   ```

## 📡 API Endpoints

### Health Checks

- `GET /health` - Comprehensive health check with system metrics
- `GET /health/live` - Liveness probe (Kubernetes style)
- `GET /health/ready` - Readiness probe (Kubernetes style)

### Sample API

- `GET /api/v1/todos` - Get all todo items
- `GET /api/v1/todos/:id` - Get a specific todo item
- `POST /api/v1/todos` - Create a new todo item
- `PUT /api/v1/todos/:id` - Update a todo item
- `DELETE /api/v1/todos/:id` - Delete a todo item
- `GET /api/v1/stats` - Get API statistics

### Root Endpoint

- `GET /` - Welcome message with basic information

## 🔧 Configuration

The application uses Viper for configuration management with the following priority:

1. Environment variables (highest priority)
2. `.env` file
3. Default values (lowest priority)

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 3000 | Server port |
| `HOST` | localhost | Server host |
| `ENV` | development | Environment (development/production) |
| `LOG_LEVEL` | info | Log level (debug/info/warn/error) |
| `LOG_FORMAT` | json | Log format (json/text) |
| `CORS_ALLOW_ORIGINS` | http://localhost:3000 | CORS allowed origins |
| `CORS_ALLOW_METHODS` | GET,POST,PUT,DELETE,OPTIONS | CORS allowed methods |
| `CORS_ALLOW_HEADERS` | Origin,Content-Type,Accept,Authorization | CORS allowed headers |
| `SHUTDOWN_TIMEOUT` | 30s | Graceful shutdown timeout |

## 📊 Logging

The application uses Go's structured logging (`slog`) with two output formats:

### JSON Format (Production)
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "info",
  "msg": "HTTP Request",
  "method": "GET",
  "path": "/api/v1/todos",
  "status": 200,
  "duration_ms": 15,
  "ip": "127.0.0.1"
}
```

### Text Format (Development)
```
2024-01-15T10:30:00Z | 200 | 15ms | 127.0.0.1 | GET | /api/v1/todos |
```

## 🏥 Health Checks

The application provides comprehensive health check endpoints:

### Main Health Check (`/health`)
Returns detailed health information including:
- Application status
- Uptime
- Memory usage
- Go runtime information
- System checks (database, external services)

### Liveness Probe (`/health/live`)
Simple liveness check for container orchestration:
```json
{
  "alive": true,
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "2h30m45s"
}
```

### Readiness Probe (`/health/ready`)
Readiness check with system dependency verification:
```json
{
  "ready": true,
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "database": "ok"
  }
}
```

## 🚦 Graceful Shutdown

The application implements graceful shutdown with:

1. **Signal Handling**: Catches SIGINT and SIGTERM signals
2. **Timeout Management**: Configurable shutdown timeout
3. **Connection Draining**: Waits for in-flight requests to complete
4. **Resource Cleanup**: Proper cleanup of resources and connections

## 🔄 API Examples

### Create a Todo
```bash
curl -X POST http://localhost:3000/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go Fiber"}'
```

### Get All Todos
```bash
curl http://localhost:3000/api/v1/todos
```

### Update a Todo
```bash
curl -X PUT http://localhost:3000/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

### Delete a Todo
```bash
curl -X DELETE http://localhost:3000/api/v1/todos/1
```

## 🛡️ Middleware Features

The middleware stack includes:

1. **Recovery**: Catches panics and returns proper error responses
2. **CORS**: Configurable cross-origin resource sharing
3. **Request Logging**: Structured logging of HTTP requests
4. **Request ID**: Unique request identification for tracing
5. **Security Headers**: Basic security headers

## 📦 Build and Deploy

### Build Binary
```bash
go build -o bin/app cmd/main.go
```

### Docker Deployment
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
```

### Production Considerations

1. **Environment**: Set `ENV=production` in production
2. **Logging**: Use JSON format for structured logging
3. **CORS**: Configure appropriate origins for your domain
4. **Timeouts**: Adjust timeouts based on your requirements
5. **Health Checks**: Configure health check intervals and timeouts
6. **Monitoring**: Use the health endpoints for monitoring

## 🧪 Testing

Run tests:
```bash
go test ./...
```

## 📝 Scripts

The `package.json` includes useful scripts:
- `npm run dev` - Development with hot reload
- `npm run build` - Build application binary
- `npm run start` - Start application
- `npm run test` - Run tests
- `npm run tidy` - Clean up dependencies

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Go Fiber](https://github.com/gofiber/fiber) - The web framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [slog](https://pkg.go.dev/log/slog) - Structured logging

## 📞 Support

If you have any questions or issues, please open an issue on GitHub.