# **Architecture Documentation: Go Fiber Skeleton**

## **1. System Overview**

This is a **Go Fiber production-ready template repository** implementing Domain-Driven Clean Architecture within a mono-repo structure. The project is currently in a **refactored state** with significant structural simplification - most domain implementations have been removed, leaving only the core foundation and agent configuration files.

### **Current State: Post-Refactoring**
- **Recent Activity**: Major refactoring to "simplify code" with removal of complete project structure
- **Remaining Files**: Core Go module, agent configuration (CLAUDE.md, agents.md, GEMINI.md)
- **Missing Components**: Domain implementations, internal package structure, database configuration
- **Architecture Pattern**: Clean Architecture with Domain-Driven Design (planned for re-implementation)

## **2. Project Structure (Template Design)**

### **Intended Directory Layout**
```
gofiber-skeleton/
├── cmd/                          # Application entry points
│   └── server/
│       └── main.go              # Main server bootstrap
├── internal/                     # Private application code
│   ├── domains/                 # Business domains
│   │   ├── user/               # User domain (reference implementation)
│   │   │   ├── entity/         # Domain entities
│   │   │   ├── repository/     # Repository interfaces
│   │   │   ├── usecase/        # Business logic
│   │   │   └── handler/        # HTTP handlers
│   │   └── [other domains]...  # Additional business domains
│   ├── infrastructure/          # Infrastructure layer
│   │   ├── database/           # Database connections
│   │   ├── middleware/         # HTTP middleware
│   │   └── config/             # Configuration management
│   └── shared/                  # Shared utilities
├── db/                          # Database management
│   ├── migrations/             # SQL migration files
│   └── queries/               # sqlc generated code
├── docs/                       # Documentation
├── pkg/                        # Public library code
├── scripts/                    # Development scripts
├── compose.yml                 # Docker Compose configuration
├── Dockerfile                  # Container configuration
├── Makefile                    # Development tasks
├── go.mod                      # Go module definition
└── README.md                   # Project documentation
```

### **Current File Structure (Post-Refactoring)**
```
gofiber-skeleton/
├── .agents/
│   └── rules/
│       └── memory-bank/        # Agent memory system (being created)
├── CLAUDE.md                   # Claude agent configuration
├── agents.md                   # Alternative agent configuration
├── GEMINI.md                   # Gemini agent configuration
├── go.mod                      # Go module (v1.25.0)
└── brief.md                    # Project requirements (in memory-bank)
```

## **3. Architecture Patterns**

### **Clean Architecture Implementation**
- **Domain Layer**: Business entities and rules
- **Repository Layer**: Data access abstractions
- **Usecase Layer**: Application business logic
- **Handler Layer**: HTTP request/response handling
- **Infrastructure Layer**: External dependencies

### **Domain-Driven Design (DDD)**
- **Domain Isolation**: Each domain is self-contained
- **Ubiquitous Language**: Domain-specific terminology
- **Bounded Contexts**: Clear domain boundaries
- **Aggregate Roots**: Domain entity relationships

## **4. Core Components (Planned Implementation)**

### **4.1. Web Framework & Routing**
- **Framework**: Go Fiber v2 (Express.js-inspired)
- **Router**: Centralized routing configuration
- **Middleware Stack**: CORS, security, rate limiting, logging
- **Handler Organization**: Domain-specific HTTP handlers

### **4.2. Data Layer**
- **ORM**: sqlc for type-safe SQL generation
- **Database**: PostgreSQL (primary)
- **Cache**: Valkey (Redis-compatible)
- **Migrations**: golang-migrate for schema evolution

### **4.3. Dependency Injection**
- **Framework**: Samber's do (generics-based)
- **Container**: Application-wide DI container
- **Interface-based**: All dependencies through interfaces

### **4.4. Authentication & Security**
- **JWT**: golang-jwt for token management
- **Password Hashing**: Argon2id algorithm
- **Middleware**: Authentication middleware
- **Validation**: Input validation and sanitization

## **5. Technology Stack**

### **Core Technologies**
- **Language**: Go 1.25.0
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL + sqlc
- **Cache**: Valkey (Redis-compatible)
- **Authentication**: JWT with Argon2id
- **Configuration**: Viper
- **Testing**: uber-go/mock + go-sqlmock

### **Development Tools**
- **Hot Reload**: Air
- **Documentation**: swaggo/swag (Swagger)
- **Linting**: golangci-lint
- **Testing**: Built-in Go testing + mocks
- **Containerization**: Docker + Docker Compose

## **6. Data Flow Architecture**

### **Request Processing Flow**
1. **HTTP Request** → Middleware Chain
2. **Middleware** → Authentication/Validation
3. **Handler** → Usecase Layer
4. **Usecase** → Business Logic Processing
5. **Repository** → Database Operations
6. **Response** → HTTP Response

### **Dependency Flow**
```
Handler → Usecase → Repository → Database
   ↑         ↑          ↑
Middleware  Config   Cache/Middleware
```

## **7. Integration Points**

### **External Systems**
- **Database**: PostgreSQL connection pool
- **Cache**: Valkey for session/storage
- **Authentication**: JWT token validation
- **Monitoring**: Health check endpoints

### **API Design**
- **RESTful**: Standard HTTP methods and status codes
- **JSON**: Request/response format
- **Documentation**: Auto-generated Swagger docs
- **Versioning**: API versioning strategy

## **8. Development Patterns**

### **Code Organization**
- **Domain-First**: Business logic drives architecture
- **Interface-Based**: All dependencies through interfaces
- **Single Responsibility**: Clear separation of concerns
- **Dependency Inversion**: High-level modules don't depend on low-level

### **Testing Strategy**
- **Unit Tests**: Isolated business logic testing
- **Integration Tests**: Database and external service testing
- **Mock-Based**: Interface mocking for isolation
- **Coverage Target**: 90%+ test coverage

## **9. Configuration Management**

### **Environment-Aware Configuration**
- **Priority**: Environment variables > .env files > defaults
- **Validation**: Configuration validation on startup
- **Typesafe**: Structured configuration with validation
- **Secrets**: Secure handling of sensitive data

## **10. Deployment Architecture**

### **Container Strategy**
- **Multi-stage**: Optimized Docker builds
- **Health Checks**: Container health monitoring
- **Environment**: Environment-specific configurations
- **Scalability**: Horizontal scaling support

### **Development Environment**
- **Docker Compose**: Local development setup
- **Hot Reload**: Automatic server restart
- **Database**: Local PostgreSQL + Valkey
- **Development Tools**: Integrated tooling

## **11. Current Implementation Status**

### **Completed (Template Foundation)**
- ✅ Go module setup (v1.25.0)
- ✅ Agent configuration systems
- ✅ Memory Bank framework (being created)
- ✅ Project documentation framework

### **Removed (Recent Refactoring)**
- ❌ Domain implementations
- ❌ Internal package structure
- ❌ Database configuration
- ❌ API implementations
- ❌ Middleware stack
- ❌ Testing infrastructure

### **Next Implementation Priorities**
1. **Re-establish Core Architecture**: Recreate internal package structure
2. **Database Integration**: Set up PostgreSQL + sqlc
3. **Web Framework**: Implement Fiber v2 with middleware
4. **Reference Domain**: Implement user/auth domain
5. **Testing Framework**: Set up testing infrastructure with mocks
6. **Development Environment**: Docker Compose + Makefile

## **12. Architectural Decisions**

### **Key Design Choices**
- **Domain Isolation**: Strict separation between business domains
- **Type Safety**: sqlc for compile-time SQL validation
- **Interface-First**: All dependencies through interfaces
- **Configuration-Driven**: Environment-aware configuration
- **Container-First**: Docker-based development and deployment

### **Trade-offs**
- **Complexity vs. Maintainability**: Choose maintainability
- **Performance vs. Simplicity**: Balanced approach
- **Flexibility vs. Convention**: Convention over configuration
- **Tooling vs. Dependencies**: Essential dependencies only

## **13. Future Architecture Evolution**

### **Scalability Considerations**
- **Microservices Ready**: Domain isolation enables service split
- **Database Scaling**: Connection pooling + query optimization
- **Caching Strategy**: Multi-level caching architecture
- **Load Balancing**: Stateless service design

### **Technology Roadmap**
- **Go Version Updates**: Follow stable Go releases
- **Framework Updates**: Track Fiber v2 evolution
- **Security Updates**: Regular dependency updates
- **Performance Monitoring**: APM integration planning

---

**Note**: This architecture documentation reflects the intended template design. The current codebase has been significantly refactored and requires re-implementation of the core architectural components according to these specifications.