# **Architecture Documentation**

## **System Overview**

This Go Fiber skeleton implements **Domain-Driven Clean Architecture** with strict separation of concerns and SOLID principles. The architecture follows a mono-repo structure with isolated business domains, shared infrastructure, and comprehensive tooling for production-ready applications.

## **Architectural Patterns**

### **Clean Architecture Layers**

```mermaid
graph TB
    subgraph "Presentation Layer"
        A[HTTP Handlers] --> B[Middleware]
        B --> C[Router]
    end
    
    subgraph "Application Layer"
        D[Use Cases] --> E[Domain Services]
    end
    
    subgraph "Domain Layer"
        F[Entities] --> G[Domain Interfaces]
        G --> H[Business Rules]
    end
    
    subgraph "Infrastructure Layer"
        I[Repositories] --> J[Database]
        I --> K[External APIs]
        L[Config] --> M[Environment]
    end
    
    A --> D
    D --> F
    D --> I
    I --> F
```

### **Directory Structure**

```
.
├── cmd/                    # Application entry points
│   ├── server/            # Main HTTP server
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── domains/          # Business domains
│   │   ├── user/         # User/auth domain
│   │   └── [domain]/     # Additional domains
│   ├── middleware/       # HTTP middleware
│   └── config/          # Configuration management
├── pkg/                  # Shared library code
│   ├── database/        # Database utilities and sqlc generated code
│   ├── cache/           # Cache utilities
│   └── response/        # Response formatting
├── db/                  # Database-related files
│   ├── migrations/      # SQL migration files
│   ├── queries/         # SQLC query files
│   └── seeds/           # Database seeds
├── docs/                # Generated documentation
├── configs/             # Configuration files
└── scripts/             # Utility scripts
```

## **Domain Architecture**

### **Domain Structure Pattern**

Each domain follows the same architectural pattern:

```
domains/[domain]/
├── entity/           # Domain entities and models
├── repository/       # Repository interfaces
├── usecase/          # Business logic use cases
├── delivery/         # HTTP handlers/transport
├── tests/            # Domain-specific tests
└── mocks/            # Generated mocks
```

### **Domain Components**

1. **Entity Layer**: Pure business objects with no external dependencies
2. **Repository Layer**: Data access interfaces with transaction management and aggregation
3. **Usecase Layer**: Application business logic and workflows
4. **Delivery Layer**: HTTP handlers and request/response processing

## **Data Flow Architecture**

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Usecase
    participant Repository
    participant Database
    
    Client->>Handler: HTTP Request
    Handler->>Handler: Validate Request
    Handler->>Usecase: Call Business Logic
    Usecase->>Repository: Data Operations
    Repository->>Repository: Transaction Management
    Repository->>Repository: Data Aggregation
    Repository->>Database: sqlc Generated Queries
    Database-->>Repository: Data Results
    Repository-->>Usecase: Domain Entities
    Usecase->>Usecase: Apply Business Rules
    Usecase-->>Handler: Use Case Results
    Handler->>Handler: Format Response
    Handler-->>Client: HTTP Response
```

## **Dependency Injection**

The project uses **Samber's do** for type-safe dependency injection:

```go
// Container setup
container := do.New()

// Register dependencies
do.Provide(container, NewConfig)
do.Provide(container, NewDatabase)
do.Provide(container, NewUserRepository)
do.Provide(container, NewUserUsecase)
do.Provide(container, NewUserHandler)
```

## **Database Architecture**

### **Database Components**

1. **Migrations**: Version-controlled schema changes using `golang-migrate`
2. **Queries**: Type-safe SQL queries using `sqlc`
3. **Generated Code**: sqlc-generated Go code in `pkg/database` package
4. **Repositories**: Data access layer with transaction management and aggregation
5. **Connection Pooling**: Optimized database connection management with pgx/v5

### **Database Flow**

```mermaid
graph LR
    A[SQL Files] --> B[sqlc Generator]
    B --> C[Generated Go Code]
    C --> D[Repository Implementation]
    D --> E[Transaction Management]
    E --> F[Data Aggregation]
    F --> G[Use Case Layer]
    G --> H[Handler Layer]
```

### **sqlc Data Access Layer Architecture**

The project uses **sqlc** as the primary data access layer, providing type-safe SQL operations with compile-time validation:

```mermaid
graph TB
    subgraph "sqlc Architecture"
        A[db/queries/*.sql] --> B[sqlc Generator]
        B --> C[pkg/database/db.go]
        B --> D[pkg/database/models.go]
        B --> E[pkg/database/queries.sql.go]
    end
    
    subgraph "Repository Layer"
        F[Repository Interface] --> G[Repository Implementation]
        G --> C
        G --> D
        G --> E
        G --> H[Transaction Management]
        G --> I[Data Aggregation]
    end
    
    subgraph "Business Logic"
        J[Use Cases] --> F
    end
```

### **Repository Layer Responsibilities**

The repository layer is responsible for:

1. **Data Access**: Using sqlc-generated code for type-safe database operations
2. **Transaction Management**: Controlling transaction state and boundaries
3. **Data Aggregation**: Performing data aggregation and complex queries
4. **Error Handling**: Translating database errors to domain errors
5. **Mapping**: Converting between database models and domain entities

### **Transaction State Management**

Transaction management is handled at the repository layer:

```mermaid
graph TB
    A[Use Case] --> B[Repository Method]
    B --> C{Transaction Needed?}
    C -->|Yes| D[Begin Transaction]
    C -->|No| E[Direct Query]
    D --> F[Execute Operations]
    F --> G{Success?}
    G -->|Yes| H[Commit Transaction]
    G -->|No| I[Rollback Transaction]
    E --> J[Return Result]
    H --> J
    I --> K[Return Error]
```

### **Data Aggregation Patterns**

Data aggregation is performed in the repository layer:

1. **Simple Aggregation**: COUNT, SUM, AVG operations
2. **Complex Queries**: JOIN operations with multiple tables
3. **Pagination**: LIMIT/OFFSET with total count queries
4. **Filtering**: Dynamic WHERE clause construction
5. **Sorting**: Multi-column sorting with direction control

### **Repository Implementation Pattern**

```go
type UserRepository interface {
    // Single operations
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    Create(ctx context.Context, user *entity.User) error
    
    // Transactional operations
    CreateWithProfile(ctx context.Context, user *entity.User, profile *entity.Profile) error
    
    // Aggregation operations
    GetUsersWithStats(ctx context.Context, filter UserFilter) ([]*entity.UserWithStats, error)
    CountByStatus(ctx context.Context) (map[string]int, error)
}

type userRepository struct {
    db      *database.DB
    queries *database.Queries // sqlc generated
}
```

## **Security Architecture**

### **Authentication & Authorization**

1. **JWT Authentication**: Stateless token-based authentication
2. **Password Hashing**: Argon2id for secure password storage
3. **Middleware Protection**: Route-level authentication checks
4. **Input Validation**: Request validation and sanitization

### **Security Middleware Stack**

```mermaid
graph TB
    A[CORS] --> B[Rate Limiting]
    B --> C[Request ID]
    C --> D[Logging]
    D --> E[Authentication]
    E --> F[Authorization]
    F --> G[Business Logic]
```

## **Testing Architecture**

### **Testing Strategy**

1. **Unit Tests**: Isolated business logic testing with mocks
2. **Integration Tests**: Database and external service testing
3. **End-to-End Tests**: Full request/response cycle testing
4. **Mock Generation**: Automated mock generation with `uber-go/mock` and `DATA-DOG/go-sqlmock`

### **Test Structure**

```
tests/
├── unit/              # Unit tests with mocks
├── integration/       # Integration tests
├── e2e/              # End-to-end tests
└── fixtures/         # Test data fixtures
```

## **Configuration Architecture**

### **Configuration Management**

1. **Environment-Based**: Environment variables for production
2. **File-Based**: `.env` files for local development
3. **Type Safety**: Structured configuration with validation
4. **Precedence Rules**: Clear configuration override order

```mermaid
graph TB
    A[Environment Variables] --> D[Final Config]
    B[.env File] --> D
    C[Default Values] --> D
```

## **API Architecture**

### **RESTful API Design**

1. **Resource-Based URLs**: Clear resource naming conventions
2. **HTTP Methods**: Proper use of HTTP verbs
3. **Status Codes**: Consistent HTTP status code usage
4. **Response Format**: Structured JSON responses with JSend format

### **API Documentation**

1. **Swagger Generation**: Automatic API documentation from code comments
2. **Interactive Docs**: Browser-based API exploration
3. **Schema Definitions**: Clear request/response schemas

## **Deployment Architecture**

### **Containerization**

1. **Multi-stage Docker**: Optimized container builds
2. **Docker Compose**: Local development environment
3. **Health Checks**: Application health monitoring
4. **Configuration Injection**: Environment-based configuration

### **CI/CD Pipeline**

```mermaid
graph TB
    A[Code Push] --> B[CI Pipeline]
    B --> C[Tests]
    C --> D[Linting]
    D --> E[Security Scan]
    E --> F[Build Image]
    F --> G[Deploy Staging]
    G --> H[Deploy Production]
```

## **Performance Architecture**

### **Performance Optimizations**

1. **Connection Pooling**: Database connection management
2. **Caching Layer**: Redis/Valkey integration for caching
3. **Efficient Queries**: Optimized SQL with proper indexing
4. **Memory Management**: Proper resource cleanup and garbage collection

### **Monitoring & Observability**

1. **Structured Logging**: Consistent log formats
2. **Request Tracing**: Request ID propagation
3. **Error Tracking**: Comprehensive error handling
4. **Metrics Collection**: Performance monitoring capabilities

## **Scalability Architecture**

### **Scaling Considerations**

1. **Horizontal Scaling**: Stateless application design
2. **Database Scaling**: Connection pooling and query optimization
3. **Cache Strategy**: Distributed caching for performance
4. **Load Balancing**: Ready for load balancer deployment

### **Domain Scalability**

1. **Domain Isolation**: Independent domain development
2. **Microservice Ready**: Easy extraction to microservices
3. **Shared Infrastructure**: Common utilities and patterns
4. **Standardized Patterns**: Consistent development approach