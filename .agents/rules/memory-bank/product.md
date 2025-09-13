# Product Overview

This project is a production-ready starter template for building HTTP APIs and services in Go using Fiber and Domain-Driven Clean Architecture. It provides a modular, domain-driven foundation that follows SOLID principles and Clean Architecture to accelerate development.

## Core Features

- Clean Architecture: Modular domains with clear boundaries.
- Dependency Injection: Powered by Uber fx.
- Configuration: Viper with YAML, .env, and environment variable support.
- Database: PostgreSQL with pgx and sqlc for type-safe SQL queries.
- Authentication: JWT-based with custom claims.
- Middlewares: Logger, Recover, CORS.
- Development: Hot reload with Air and Swagger docs generation.
- Deployment: Multi-stage Docker build and Docker Compose setup.

## Goals

- Provide a standard template to bootstrap new Go API services.
- Enforce domain-driven design and SOLID principles out of the box.
- Simplify common infrastructure setup (database, configuration, logging, validation).
- Enable easy local development and production readiness with minimal configuration.