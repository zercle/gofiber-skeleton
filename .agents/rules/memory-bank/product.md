# Product Overview

## Purpose
Go Fiber Skeleton is a production-ready backend template that accelerates development of Go-based REST APIs. It eliminates repetitive boilerplate setup by providing a fully-configured, best-practice foundation for building scalable backend services.

## Problems Solved
- **Time-to-First-Line-of-Business-Logic**: Developers spend significant time setting up infrastructure (database connections, migrations, authentication, configuration). This template reduces that to minutes.
- **Architectural Consistency**: Teams struggle with maintaining consistent patterns across services. This template enforces Domain-Driven Clean Architecture from day one.
- **Production Readiness**: Many projects launch without proper logging, health checks, or containerization. This template includes all production essentials pre-configured.

## Target Audience
- **Backend developers** building new Go microservices or REST APIs
- **Teams** needing consistent architecture across multiple services
- **Startups/Projects** requiring rapid backend development without sacrificing code quality
- **Developers transitioning** from other frameworks (Node.js, Python) to Go

## User Experience Goals
1. **Clone-and-Go**: Developer clones repository, runs `docker-compose up`, and has a working API with database in under 5 minutes
2. **Clear Domain Addition**: Following documented steps, adding a new business domain (e.g., "products", "orders") takes 15-20 minutes
3. **Type-Safe Database Operations**: All database queries are type-checked at compile time via sqlc, eliminating runtime SQL errors
4. **Self-Documenting API**: Swagger documentation auto-generates from code comments, staying always in sync

## Key Workflows
1. **Initial Setup**: Clone → Copy `.env.example` to `.env` → Run `docker-compose up` → API accessible at localhost:8080
2. **Add Domain**: Create domain directory structure → Define entity → Write SQL queries → Run `make sqlc` → Implement usecase → Add handler → Register routes
3. **Development Cycle**: Edit code → Air hot-reloads server → Test via Swagger UI → Commit changes
4. **Testing**: Run `make test` for unit tests with mocked dependencies

## Success Metrics
- **Adoption**: Developer successfully creates first custom domain within 30 minutes of cloning
- **Confidence**: All database operations are compile-time type-safe
- **Completeness**: Template includes auth, migrations, logging, containerization, and API docs out-of-box
- **Maintainability**: Clear separation of concerns allows independent testing and modification of each layer