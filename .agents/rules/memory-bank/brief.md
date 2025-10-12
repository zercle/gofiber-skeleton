# **Project Brief: Go Fiber Production-Ready Template**

## **1. The Foundation: Template Repository**

This is a **production-ready template repository** providing a complete backend foundation using **Go** and the **Fiber v2** framework. The architecture implements **Domain-Driven Clean Architecture** within a **mono-repo structure**, emphasizing strict domain isolation and SOLID principles. This template eliminates 80-90% of initial project setup work, allowing developers to focus on business logic rather than infrastructure configuration.

## **2. High-Level Overview: What This Template Provides**

The goal is to provide a robust, production-ready starting point that significantly accelerates the development of new Go backend services. Developers can clone this repository and immediately start adding business domains by following clear documentation, with all infrastructure, tooling, and best practices already configured. The template includes a complete reference implementation (user/auth domain) demonstrating every architectural pattern, testing strategy, and integration point.

## **3. Core Requirements and Goals**

### **Key Technologies:**

* **Framework:** Go Fiber v2 - A high-performance, Express.js-inspired web framework for building APIs.  
* **Dependency Injection:** Samber's do - A modern, lightweight dependency injection framework based on Go 1.18+ Generics for type-safe dependency management.
* **Configuration:** Viper - For robust configuration management that supports .env files, YAML, and environment variables with clear precedence rules.  
* **Database Migrations:** golang-migrate/migrate - To handle database schema evolution with versioned, repeatable SQL migration files.  
* **SQL Generation:** sqlc - For generating fully type-safe, idiomatic Go code directly from raw SQL queries, catching errors at compile-time.  
* **Authentication:** golang-jwt - A standard library for creating and verifying JSON Web Tokens for secure, stateless authentication.  
* **API Documentation:** swaggo/swag - To automatically generate interactive API documentation from comments in the Go source code.  
* **Development:** Hot-reloading with Air - To improve the development feedback loop by automatically rebuilding and restarting the server on file changes.  
* **Testing & Mocking:**
  * uber-go/mock/mockgen: For auto-generating mock implementations of interfaces using //go:generate annotations.
  * DATA-DOG/go-sqlmock: For mocking SQL driver in testing.
  * Complete test coverage examples in the reference domain with comprehensive mocking strategies.

### **Template Features:**

* **Project Structure:** A clear and scalable directory structure that intuitively separates domains, shared infrastructure (like database connections and middleware), and command entry points (/cmd, /internal/domains, /db, etc.).  
* **Configuration Management:** A flexible and environment-aware config system that prioritizes runtime environment variables for production but falls back to a .env file for easy local development.  
* **Database Tooling:** Integrated tools and scripts for running migrations to update the database schema and for generating type-safe Go queries, ensuring database operations are reliable and maintainable.  
* **Reference Implementation:** A complete user/auth domain demonstrating user registration with password hashing, JWT authentication, testing patterns, and all architectural layers.  
* **Containerization:** A compose.yml file that allows developers to spin up a consistent, isolated development environment with a single command, including PostgreSQL and Valkey (Redis-compatible) cache.  
* **Developer Experience:** A comprehensive Makefile with commands for all common tasks: running tests, linting code, building binaries, generating documentation, running migrations, and generating mocks.  
* **Comprehensive Testing Strategy:** The template demonstrates robust testing approaches:  
  * **Unit Testability:** Using go.uber.org/mock/mockgen with //go:generate annotations on all repository and usecase interfaces to ensure business logic can be tested in complete isolation.
  * **Complete Examples:** Full test coverage in the user domain showing mock-based testing, validation testing, and error handling.
* **Clear Instructions:** 
  * **TEMPLATE_SETUP.md:** Complete guide for initializing a new project from this template
  * **docs/ADDING_NEW_DOMAIN.md:** Step-by-step guide for adding new business domains with the user domain as reference

## **4. Template Value Proposition**

### **What's Included (80-90% of Setup Work)**

* ✅ Complete Clean Architecture implementation with working examples
* ✅ Production-ready configuration management
* ✅ Database integration with migrations and type-safe queries
* ✅ Authentication system with JWT and password hashing
* ✅ Comprehensive middleware stack (CORS, security, rate limiting, logging)
* ✅ Development environment with Docker Compose
* ✅ Testing framework with mocks and examples
* ✅ API documentation generation with Swagger
* ✅ CI/CD foundation with health checks
* ✅ Deployment-ready containerization

### **What Developers Add (10-20% Business Logic)**

* Business-specific domains following the reference pattern
* Custom validation rules and business logic
* Domain-specific database migrations
* API endpoints for their use cases
* Additional third-party integrations

## **5. Reference Implementation: User/Auth Domain**

The template includes a **complete, production-ready user/auth domain** serving as the reference implementation. This domain demonstrates:

* **Entity Layer:** User model with proper domain modeling
* **Repository Layer:** PostgreSQL integration with sqlc, mock generation
* **Usecase Layer:** Registration and login logic with comprehensive tests
* **Handler Layer:** HTTP endpoints with Swagger documentation
* **Middleware:** JWT authentication middleware
* **Testing:** 90%+ test coverage with mocks and isolation
* **Security:** Password hashing with Argon2id, JWT token generation
* **Validation:** Input validation and error handling patterns

Developers can use this domain as a template for their own domains, following the same patterns and structure.

## **6. Getting Started Workflow**

### **For New Projects:**

1. **Clone and Initialize** (5 minutes)
   - Clone the repository
   - Follow TEMPLATE_SETUP.md to customize module name and configuration
   - Start development environment with `make dev`

2. **Explore Reference Implementation** (15 minutes)
   - Review user domain structure and patterns
   - Run tests with `make test`
   - Explore API documentation at `/swagger/`

3. **Add First Business Domain** (30-60 minutes)
   - Follow docs/ADDING_NEW_DOMAIN.md guide
   - Create domain structure following user domain pattern
   - Implement business logic with tests
   - Add domain routes to router

4. **Deploy to Production** (varies)
   - Configure environment variables
   - Review Dockerfile and deployment configuration
   - Set up health check monitoring
   - Deploy containerized application

## **7. Quality Standards and Best Practices**

### **Enforced by Template:**

* **Architecture:** Strict Clean Architecture with domain isolation
* **Testing:** Mock-based unit testing with 90%+ coverage target
* **Security:** JWT authentication, password hashing, input validation
* **Performance:** Connection pooling, caching, efficient queries
* **Documentation:** Auto-generated Swagger docs for all endpoints
* **Code Quality:** golangci-lint with comprehensive rules
* **Type Safety:** sqlc for compile-time SQL validation

### **Developer Benefits:**

* No need to research best practices - they're built-in
* No need to set up tooling - it's pre-configured
* No need to write boilerplate - reference implementation provided
* No need to configure deployment - Docker setup included
* No need to implement auth - production-ready auth included

## **8. Template Maintenance and Evolution**

This template is maintained with:

* Regular dependency updates for security and features
* Best practice improvements based on Go ecosystem evolution
* Community feedback incorporation
* Documentation updates and clarifications
* New feature examples and patterns

## **9. Success Criteria**

A developer should be able to:

* ✅ Clone and run the template in under 5 minutes
* ✅ Understand the architecture by reviewing the user domain
* ✅ Add a new CRUD domain in under 1 hour
* ✅ Deploy to production with minimal configuration
* ✅ Maintain high code quality with built-in tooling
* ✅ Scale the application by adding more domains following the same pattern

## **10. Target Audience**

* **Development Teams** starting new Go backend projects
* **Individual Developers** building production applications
* **Organizations** standardizing on Go for microservices
* **Startups** needing rapid development without sacrificing quality
* **Enterprises** requiring consistent architecture across teams
* **Developers** learning Go best practices with real-world examples