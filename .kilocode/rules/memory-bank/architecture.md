# Architecture

## System Architecture Overview
The E-commerce Backend Boilerplate follows a Clean Architecture approach, decoupling high-level business rules from framework and infrastructure concerns. The system comprises:
- REST API entrypoint (Fiber)  
- HTTP handlers orchestrating requests (internal/handler)  
- Use case services encapsulating business logic (internal/usecase)  
- Domain entities and interfaces defining business contracts (internal/domain)  
- Repository implementations managing data persistence (internal/repository/db)  
- Supporting shared utilities (pkg)  
- Database migrations (migrations) and SQL query generation (SQLC)

## Layered Architecture
- Presentation / API Layer  
  • Location: cmd/server/main.go  
  • Responsibilities: Bootstraps Fiber app, registers routes and middleware.
- Handler Layer  
  • Location: internal/handler  
  • Responsibilities: Validate request payloads, call use cases, format HTTP responses.
- Use Case / Service Layer  
  • Location: internal/usecase  
  • Responsibilities: Encapsulate business rules and orchestrate repository calls.
- Domain Layer
  • Location: internal/domain
  • Responsibilities: Contain only domain models (entities) and interfaces for repository and usecase boundaries; do not include business logic or implementations.
- Repository / Infrastructure Layer  
  • Location: internal/repository/db and internal/repository  
  • Responsibilities: Implement domain repository interfaces using SQLC-generated queries, manage database connections.
- Shared Utilities  
  • Location: pkg  
  • Responsibilities: Common helpers, error definitions, configuration loaders.

## Component Relationships
mermaid
graph LR
  A[cmd/server] --> B[internal/handler]
  B --> C[internal/usecase]
  C --> D[internal/domain/interfaces]
  C --> E[internal/repository/db]
  E --> F[(Postgres via SQLC)]
  A --> G[pkg utilities]

## Data Flow
1. Client issues HTTP request to Fiber server.  
2. Fiber router dispatches to appropriate handler.  
3. Handler validates input and converts to domain model.  
4. Handler invokes a use case service method.  
5. Use case orchestrates domain logic and calls repository interface.  
6. Repository implementation executes SQLC-generated queries against Postgres.  
7. Results propagate back through use case and handler to the HTTP response.

## Design Patterns & Key Decisions
- Clean Architecture for testable, modular boundaries.  
- Dependency inversion: handlers and use cases depend on interfaces, not concrete implementations.  
- SQLC for compile-time safe SQL queries and type generation.
- SQLC-generated code acts as entity providers; repositories should orchestrate and aggregate data before returning to use cases.
- Golang-migrate for reliable, reversible database migrations.  
- JWT middleware for stateless authentication.  
- Go-mock and go-sqlmock for interface and database mocks in unit tests.
- Place generated mock files in a `mock` subpackage within each owner package (e.g., `internal/repository/mock`, `internal/usecase/mock`).
- Docker Compose for isolated service orchestration (app + Postgres).

## Critical Implementation Paths
- Initializing and wiring up dependencies in cmd/server (constructor injection).  
- SQLC code generation: queries in queries/*.sql → generate Go types and methods in internal/repository/db.  
- Migrations: ensure schema aligns with SQLC queries and domain models.  
- Middleware pipeline: logging, recovery, JWT validation before handler execution.  
- Unit and integration tests validating handler-use case interactions with mocks.