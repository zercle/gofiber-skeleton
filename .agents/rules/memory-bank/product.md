# Product Overview: Go Fiber Template Repository

## Goal
The primary goal of this project is to establish a robust, production-ready backend template using Go and the Fiber v2 framework. This template aims to significantly accelerate the development of new Go backend services by providing a comprehensive "skeleton" project. It comes pre-configured with established best practices and a curated set of tools, handling common backend tasks such as database interaction, configuration, and authentication.

## Key Benefits
- **Rapid Development:** Developers can clone the repository and quickly add new business domains, focusing immediately on business logic rather than boilerplate setup.
- **Maintainability:** The architecture emphasizes strict domain isolation and adherence to SOLID principles, leading to a highly maintainable and scalable codebase.
- **Scalability:** Each business domain (e.g., auth, posts) is designed as a self-contained, independently testable module with zero dependencies on other domains, facilitating long-term project health and team collaboration.
- **Consistency:** Provides a consistent foundation for all new Go backend services, ensuring uniformity in structure, tooling, and practices.

## Target Audience
This template is intended for developers and teams looking to build high-performance, scalable, and maintainable backend services in Go, particularly those who appreciate a clean architecture approach and efficient development workflows.

## Core Features
- **Clear Project Structure:** Intuitive separation of domains, shared infrastructure, and command entry points.
- **Flexible Configuration:** Environment-aware configuration system with support for `.env` files, YAML, and environment variables.
- **Integrated Database Tooling:** Tools and scripts for database migrations and type-safe SQL query generation (sqlc).
- **Authentication:** A pre-built, functional authentication domain demonstrating user registration, password hashing, and JWT-based stateless authentication.
- **Containerization:** `compose.yml` for a consistent development environment with PostgreSQL and Valkey/Redis.
- **Developer Experience:** Helper scripts/Make commands for common tasks (testing, linting, building, API documentation).
- **Robust Testing Strategy:** Comprehensive unit and repositories testing examples using `go.uber.org/mock/mockgen` and `DATA-DOG/go-sqlmock`.
- **Comprehensive Documentation:** Clear instructions for extending the template with new business domains.