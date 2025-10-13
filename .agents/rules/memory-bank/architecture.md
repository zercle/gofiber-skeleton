# **System Architecture: Go Fiber Skeleton**

## **1. Architecture Overview**

Production-ready Go backend template implementing **Domain-Driven Clean Architecture** with **strict domain isolation** and **SOLID principles**. **Eliminates 80-90% of initial project setup work** by providing complete infrastructure, tooling, and reference implementation.

### **Core Principles**

* **Clean Architecture Layers:** Domain → Application → Infrastructure
* **Domain-Driven Design:** Business logic encapsulated in isolated domains
* **Dependency Injection:** Type-safe DI with Samber's do framework
* **Mono-repo Structure:** Organized codebase with clear boundaries
* **Test-Driven Development:** Comprehensive testing with mocks

## **2. Directory Structure**

```
gofiber-skeleton/
├── cmd/server/                 # Application entry point
├── internal/
│   ├── domains/               # Business domains (DDD)
│   │   └── user/              # Reference domain
│   │       ├── entity/        # Domain models
│   │       ├── repository/    # Repository interfaces
│   │       ├── usecase/       # Business logic
│   │       └── handler/       # HTTP handlers
│   ├── infrastructure/        # External implementations
│   │   ├── database/         # Database setup
│   │   ├── middleware/       # HTTP middleware
│   │   └── config/           # Configuration
│   └── shared/               # Shared utilities
├── db/
│   ├── migrations/           # Schema migrations
│   └── queries/              # SQL queries for sqlc
├── docs/                     # Documentation
├── scripts/                  # Utility scripts
├── compose.yml              # Development environment
├── Dockerfile               # Production container
└── Makefile                 # Development commands
```

## **3. Request Flow**

```mermaid
graph TD
    A[HTTP Request] --> B[Middleware Chain]
    B --> C[Handler Layer]
    C --> D[Usecase Layer]
    D --> E[Repository Interface]
    E --> F[Repository Implementation]
    F --> G[Database/External Service]
    G --> H[HTTP Response]
    
    I[DI Container] --> C
    I --> D
    I --> F
```

## **4. Layer Dependencies**

* **Handler:** Depends on Usecase interfaces
* **Usecase:** Depends on Repository interfaces + Domain entities
* **Repository:** Implements interfaces, depends on infrastructure
* **Entity:** Pure domain models, no external dependencies
* **Infrastructure:** Provides implementations for inner layers

## **5. Domain Pattern**

Each domain follows this consistent structure:

1. **Entity:** Core business models with validation
2. **Repository:** Interface definitions for data access
3. **Usecase:** Business logic with dependency injection
4. **Handler:** HTTP endpoints with Swagger documentation

## **6. Technology Integration**

### **Database Layer**
* **PostgreSQL:** Primary database with pgx driver
* **Migrations:** golang-migrate for schema versioning
* **Queries:** sqlc for type-safe SQL generation
* **Transactions:** Consistent transaction handling

### **Web Layer**
* **Fiber v2:** High-performance web framework
* **Middleware:** CORS, security, rate limiting, logging
* **Authentication:** JWT-based stateless auth
* **Documentation:** Auto-generated Swagger/OpenAPI

### **Development Tools**
* **Dependency Injection:** Samber's do framework
* **Configuration:** Viper with environment hierarchy
* **Testing:** uber-go/mock for interface mocking
* **Hot Reloading:** Air for development efficiency

## **7. Reference Implementation**

The **user domain** demonstrates all architectural patterns:

* **Entity:** User model with password hashing
* **Repository:** PostgreSQL integration with sqlc
* **Usecase:** Registration/login business logic
* **Handler:** HTTP endpoints with comprehensive tests
* **Testing:** 90%+ coverage with proper mocking

## **8. Configuration Management**

### **Priority Order**
1. Environment Variables (production)
2. .env File (development)
3. Default Values (fallback)

### **Key Areas**
* Server settings (host, port, timeouts)
* Database configuration (connection, pooling)
* JWT settings (secrets, expiration)
* Logging configuration (levels, output)

## **9. Testing Strategy**

* **Unit Tests:** Business logic with mocked dependencies
* **Integration Tests:** Database and external service testing
* **Handler Tests:** HTTP endpoint testing
* **Mock Generation:** Automatic with //go:generate annotations

## **10. Production Considerations**

### **Containerization & Deployment**
* **Multi-stage Docker builds:** Optimized production images
* **Health checks:** Liveness and readiness probes
* **Graceful shutdown:** Proper resource cleanup on termination
* **Resource limits:** CPU and memory constraints
* **Security scanning:** Container vulnerability assessment

### **Observability & Monitoring**
* **Structured logging:** JSON-formatted logs with correlation IDs
* **Metrics collection:** Application performance metrics
* **Error tracking:** Centralized error reporting
* **Request tracing:** Distributed tracing for debugging
* **Health monitoring:** External dependency health checks

### **Security Best Practices**
* **Input validation:** Comprehensive validation and sanitization
* **Rate limiting:** API endpoint protection
* **CORS configuration:** Proper cross-origin resource sharing
* **Security headers:** HSTS, CSP, X-Frame-Options
* **Secrets management:** Environment-based configuration
* **Dependency updates:** Regular security patching

### **Performance Optimization**
* **Connection pooling:** Database connection optimization
* **Caching strategy:** Multi-level caching implementation
* **Response compression:** Gzip for API responses
* **Database indexing:** Proper query optimization
* **Memory management:** Efficient Go patterns

### **Scalability Considerations**
* **Horizontal scaling:** Stateless application design
* **Database scaling:** Read replicas and connection management
* **Load balancing:** Multiple instance support
* **Caching layer:** Distributed caching with Valkey
* **Background jobs:** Asynchronous task processing