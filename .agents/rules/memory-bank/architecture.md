# Architecture: Go Fiber Template Repository

## Core Architectural Principles

This project is built upon a **Domain-Driven Clean Architecture** within a **mono-repo structure**. This design choice prioritizes:

1.  **Domain Isolation:** Each business domain (e.g., `auth`, `posts`) is a self-contained module, minimizing dependencies between domains. This enhances modularity and reduces the blast radius of changes.
2.  **SOLID Principles:** Strict adherence to Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion principles.
3.  **Testability:** Components are designed for independent testing, with interfaces and dependency injection facilitating easy mocking.
4.  **Scalability & Maintainability:** The clear separation of concerns and modular design support long-term project health and team collaboration.
5.  **Abstraction-Driven Development:** Public components interact with abstractions (interfaces) rather than concrete implementations, promoting flexibility and testability.
6.  **Composition over Inheritance:** Favoring smaller, purpose-specific abstractions.

## Feature-Based Architecture Overview

The architecture is organized around individual features, each encapsulated within its own directory structure. This design prioritizes feature cohesion and simplifies navigation by grouping all code related to a specific feature (handlers, usecases, repositories, tests, configs) together.

### Directory Layout

```text
project/
├── cmd/
│   └── server/
│       └── main.go      # Main application logic
├── internal/
│   ├── config/          # Configuration loading
│   ├── db/              # Database connections and models
│   ├── user/            # Feature: User (entity, model, repository, usecase, handler, middleware, tests)
│   │   ├── entity/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── tests/
│   ├── post/            # Feature: Post (entity, model, repository, usecase, handler, middleware, tests)
│   │   ├── entity/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── tests/
│   └── server/          # Fiber router setup
├── db/                  # SQL files
│   ├── migrations/      # Database migration scripts
│   └── queries/         # sqlc query definitions
├── migrations/          # SQL schema migrations
├── go.mod               # Go module definition
└── go.sum               # Go module checksum file
```

This feature-based structure enhances modularity by keeping feature logic self-contained, reducing dependencies across unrelated features, and improving maintainability as the application grows.

Implemented user authentication and post CRUD features under `internal/user` and `internal/post`. User authentication uses `golang.org/x/crypto/bcrypt` for password hashing and `github.com/golang-jwt/jwt/v4` for JWT-based sessions.

## Dependency Flow

Dependencies generally flow inwards:
`Entrypoints` → `API/Transport` → `Application/Service` → `Domain Models`
`Application/Service` depends on `Data Access` (via interfaces).
`Shared Libraries/Utilities` are depended upon by other layers.

## Key Architectural Decisions

-   **Uber fx for DI:** Provides a structured way to manage dependencies, promote modularity, and simplify application startup/shutdown.
-   **sqlc for Type-Safe SQL:** Eliminates manual SQL query writing and reduces runtime errors by generating Go code from SQL.
-   **golang-migrate/migrate:** Ensures controlled and versioned database schema evolution.
-   **JWT for Authentication:** Stateless, scalable authentication mechanism.
-   **User Auth Dependencies:** Password hashing with `golang.org/x/crypto/bcrypt` and JWT sessions via `github.com/golang-jwt/jwt/v4`.
-   **Mono-repo with Domain Separation:** Allows for shared tooling and simplified dependency management while maintaining strong domain boundaries.
-   **JSend Compliant Responses:** Standardized API response format for consistency and clarity.

## Future Considerations

-   **Event-Driven Architecture:** Introduction of a message bus (`internal/shared/bus/`) for inter-domain communication and asynchronous processing.
-   **Observability:** Further integration of tracing, metrics, and structured logging across all layers.
-   **Error Handling:** Consistent, centralized error handling strategy with clear error codes and messages.