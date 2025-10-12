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
│   │   │   ├── entity/   # Domain entities and sqlc generated code
│   │   │   ├── repository/ # Repository interfaces
│   │   │   ├── usecase/  # Business logic use cases
│   │   │   ├── delivery/ # HTTP handlers/transport
│   │   │   ├── tests/    # Domain-specific tests
│   │   │   └── mocks/    # Generated mocks
│   │   ├── post/         # Post domain
│   │   │   ├── entity/   # Domain entities and sqlc generated code
│   │   │   ├── repository/ # Repository interfaces
│   │   │   ├── usecase/  # Business logic use cases
│   │   │   ├── delivery/ # HTTP handlers/transport
│   │   │   ├── tests/    # Domain-specific tests
│   │   │   └── mocks/    # Generated mocks
│   │   └── [domain]/     # Additional domains
│   ├── middleware/       # HTTP middleware
│   └── config/          # Configuration management
├── pkg/                  # Shared library code
│   ├── cache/           # Cache utilities
│   └── response/        # Response formatting
├── db/                  # Database-related files
│   ├── migrations/      # SQL migration files
│   ├── queries/         # SQLC query files (per domain)
│   │   ├── user.sql     # User domain queries
│   │   └── post.sql     # Post domain queries
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
├── entity/           # Domain entities and sqlc generated code
├── repository/       # Repository interfaces
├── usecase/          # Business logic use cases
├── delivery/         # HTTP handlers/transport
├── tests/            # Domain-specific tests
└── mocks/            # Generated mocks
```

### **Domain Components**

1. **Entity Layer**: Pure business objects with no external dependencies and sqlc generated code
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
do.Provide(container, NewPostRepository)
do.Provide(container, NewUserUsecase)
do.Provide(container, NewPostUsecase)
do.Provide(container, NewUserHandler)
do.Provide(container, NewPostHandler)
```

## **Database Architecture**

### **Database Components**

1. **Migrations**: Version-controlled schema changes using `golang-migrate`
2. **Queries**: Type-safe SQL queries using `sqlc`
3. **Generated Code**: sqlc-generated Go code in `internal/domains/*/entity/` packages
4. **Repositories**: Data access layer with transaction management and aggregation
5. **Connection Pooling**: Optimized database connection management with pgx/v5

### **Database Flow**

```mermaid
graph LR
    A[SQL Files] --> B[sqlc Generator]
    B --> C[Domain Entity Code]
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
        A[db/queries/user.sql] --> B[sqlc Generator]
        A2[db/queries/post.sql] --> B
        B --> C[internal/domains/user/entity/db.go]
        B --> D[internal/domains/user/entity/models.go]
        B --> E[internal/domains/user/entity/queries.sql.go]
        B --> F[internal/domains/post/entity/db.go]
        B --> G[internal/domains/post/entity/models.go]
        B --> H[internal/domains/post/entity/queries.sql.go]
    end
    
    subgraph "Repository Layer"
        I[User Repository Interface] --> J[User Repository Implementation]
        K[Post Repository Interface] --> L[Post Repository Implementation]
        J --> C
        J --> D
        J --> E
        L --> F
        L --> G
        L --> H
        J --> M[Transaction Management]
        L --> M
        J --> N[Data Aggregation]
        L --> N
    end
    
    subgraph "Business Logic"
        O[User Use Cases] --> I
        P[Post Use Cases] --> K
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
    queries *entity.Queries // sqlc generated in entity package
}

type PostRepository interface {
    // Single operations
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
    Create(ctx context.Context, post *entity.Post) error
    
    // User-related operations
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Post, error)
    
    // Aggregation operations
    GetPostsWithAuthor(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error)
    GetUserPostStats(ctx context.Context, userID uuid.UUID) (*entity.PostStats, error)
}

type postRepository struct {
    db      *database.DB
    queries *entity.Queries // sqlc generated in entity package
}
```

## **Domain Relationships**

### **User-Post Relationship**

```mermaid
graph TB
    subgraph "User Domain"
        A[User Entity] --> B[User Repository]
        B --> C[User Use Cases]
        C --> D[User Handlers]
    end
    
    subgraph "Post Domain"
        E[Post Entity] --> F[Post Repository]
        F --> G[Post Use Cases]
        G --> H[Post Handlers]
    end
    
    subgraph "Database"
        I[users table] --> J[posts table]
        J --> I
    end
    
    B --> I
    F --> J
    F --> I
```

### **Cross-Domain Interactions**

1. **User-Post Relationship**: Posts belong to users
2. **Shared Transactions**: Operations spanning multiple domains
3. **Data Aggregation**: Posts with user information
4. **Authorization**: User-based post access control

## **Security Architecture**

### **Authentication & Authorization**

1. **JWT Authentication**: Stateless token-based authentication
2. **Password Hashing**: Argon2id for secure password storage
3. **Middleware Protection**: Route-level authentication checks
4. **Input Validation**: Request validation and sanitization
5. **Resource Authorization**: User can only access their own posts

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

### **API Endpoints**

#### **User Endpoints**
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

#### **Post Endpoints**
- `GET /api/v1/posts` - List posts (with pagination)
- `GET /api/v1/posts/:id` - Get post by ID
- `POST /api/v1/posts` - Create new post
- `PUT /api/v1/posts/:id` - Update post
- `DELETE /api/v1/posts/:id` - Delete post
- `GET /api/v1/users/:id/posts` - Get user's posts

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