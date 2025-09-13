# Architecture: Go Fiber Template Repository

## Core Architectural Principles

This project is built upon a **Domain-Driven Clean Architecture** within a **mono-repo structure**. This design choice prioritizes:

1.  **Domain Isolation:** Each business domain (e.g., `auth`, `posts`) is a self-contained module, minimizing dependencies between domains. This enhances modularity and reduces the blast radius of changes.
2.  **SOLID Principles:** Strict adherence to Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion principles.
3.  **Testability:** Components are designed for independent testing, with interfaces and dependency injection facilitating easy mocking.
4.  **Scalability & Maintainability:** The clear separation of concerns and modular design support long-term project health and team collaboration.
5.  **Abstraction-Driven Development:** Public components interact with abstractions (interfaces) rather than concrete implementations, promoting flexibility and testability.
6.  **Composition over Inheritance:** Favoring smaller, purpose-specific abstractions.

## Layered Architecture Overview

The architecture generally follows a layered approach, though the mono-repo structure allows for domain-specific layering.

### 1. Entrypoints/Applications (`cmd/`)
-   Contains the main application entry points (e.g., `cmd/server/main.go`, `cmd/migrate/main.go`).
-   Responsible for bootstrapping the application, setting up dependency injection, and starting the server or executing commands.

### 2. API/Transport Layer (`internal/app/http/`, `internal/shared/server/`)
-   Defines API routes, request/response models, and marshaling/unmarshaling.
-   Includes middleware for common concerns like authentication, logging, and error recovery.
-   Interacts with the `handlers` Layer.
-   Handles external communication, primarily HTTP requests using the Fiber v2 framework.
-   Defines API routes, request/response models, and marshaling/unmarshaling.
-   Includes middleware for common concerns like authentication, logging, and error recovery.
-   Interacts with the Application/Service Layer.

### 3. Application/Service Layer (`internal/domains/<domain>/usecases/`)
-   Orchestrates business logic and use cases for specific domains.
-   Contains domain-specific services that coordinate operations across multiple domain entities and interact with repositories.
-   Does not contain framework-specific code.
-   Defines interfaces for external dependencies (e.g., repositories, external services).

### 4. Domain Models (`internal/domains/<domain>/entities/`, `internal/domains/<domain>/models/`)
-   Represents the core business concepts and rules.
-   `entities/`: Contains the rich domain objects and their behaviors.
-   `models/`: Contains data transfer objects (DTOs) for input/output, often used in the API layer.

### 5. Data Access Layer (`internal/domains/<domain>/repository/`, `internal/infrastructure/database/`)
-   Provides an abstraction over data persistence mechanisms.
-   `repository/`: Contains repository implementations for specific domains, adhering to interfaces defined in the Application/Service Layer.
-   `internal/infrastructure/database/`: Handles generic database connection management, SQL query generation (via sqlc), and migration logic.

### 6. Shared Libraries/Utilities (`internal/shared/`, `internal/infrastructure/`)
-   Contains reusable components and utilities that are not specific to any single domain.
-   Examples:
    -   `internal/shared/di/`: Dependency Injection helpers (Uber fx).
    -   `internal/shared/jsend/`: JSend-compliant API response formatting.
    -   `internal/infrastructure/config/`: Configuration loading (Viper).
    -   `internal/infrastructure/logging/`: Centralized logging.
    -   `internal/infrastructure/trace/`: Distributed tracing.
    -   `internal/infrastructure/validation/`: Input validation.
    -   `internal/infrastructure/middleware/`: Common Fiber middleware.

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
-   **Mono-repo with Domain Separation:** Allows for shared tooling and simplified dependency management while maintaining strong domain boundaries.
-   **JSend Compliant Responses:** Standardized API response format for consistency and clarity.

## Future Considerations

-   **Event-Driven Architecture:** Introduction of a message bus (`internal/shared/bus/`) for inter-domain communication and asynchronous processing.
-   **Observability:** Further integration of tracing, metrics, and structured logging across all layers.
-   **Error Handling:** Consistent, centralized error handling strategy with clear error codes and messages.