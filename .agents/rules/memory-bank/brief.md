# **Project Brief: Go Fiber Template Repository**

## **1. The Foundation: Core Architecture**

We are building a backend template using **Go** and the **Fiber v2** framework. The foundational architecture will be a **Domain-Driven Clean Architecture** implemented within a **mono-repo structure**. This strategic choice emphasizes strict domain isolation and adherence to SOLID principles. The primary benefit is creating a highly maintainable and scalable codebase where each business domain (e.g., auth, posts) is a self-contained, independently testable module with zero dependencies on other domains. This separation of concerns is crucial for long-term project health and team collaboration.

## **2. High-Level Overview: What We're Building**

The goal is to create a robust, production-ready template that significantly accelerates the development of new Go backend services. This template will serve as a comprehensive "skeleton" project, pre-configured with established best practices and a curated set of tools for handling common backend tasks like database interaction, configuration, and authentication. The ideal outcome is a developer experience where one can clone the repository, quickly add a new business domain by following a clear guide, and immediately focus on writing valuable business logic without getting bogged down in repetitive boilerplate setup.

## **3. Core Requirements and Goals**

### **Key Technologies:**

* **Framework:** Go Fiber v2 - A high-performance, Express.js-inspired web framework for building APIs.  
* **Dependency Injection:** Uber's fx - To manage dependencies and ensure a modular, loosely coupled application structure.  
* **Configuration:** Viper - For robust configuration management that supports .env files, YAML, and environment variables with clear precedence rules.  
* **Database Migrations:** golang-migrate/migrate - To handle database schema evolution with versioned, repeatable SQL migration files.  
* **SQL Generation:** sqlc - For generating fully type-safe, idiomatic Go code directly from raw SQL queries, catching errors at compile-time.  
* **Authentication:** golang-jwt - A standard library for creating and verifying JSON Web Tokens for secure, stateless authentication.  
* **API Documentation:** swaggo/swag - To automatically generate interactive API documentation from comments in the Go source code.  
* **Development:** Hot-reloading with Air - To improve the development feedback loop by automatically rebuilding and restarting the server on file changes.  
* **Testing & Mocking:**  
  * go.uber.org/mock/mockgen: For auto-generating mock implementations of interfaces using //go:generate annotations.  
  * DATA-DOG/go-sqlmock: For mocking the SQL database layer to test data access logic without a real database.

### **Must-Have Features:**

* **Project Structure:** A clear and scalable directory structure that intuitively separates domains, shared infrastructure (like database connections and middleware), and command entry points (/cmd, /internal/domains, /db, etc.).  
* **Configuration Management:** A flexible and environment-aware config system that prioritizes runtime environment variables for production but falls back to a .env file for easy local development.  
* **Database Tooling:** Integrated tools and scripts for running migrations to update the database schema and for generating type-safe Go queries, ensuring database operations are reliable and maintainable.  
* **Authentication:** A pre-built and fully functional auth domain to serve as a practical example, demonstrating user registration with password hashing and stateless authentication using JSON Web Tokens (JWT).  
* **Containerization:** A compose.yml file that allows developers to spin up a consistent, isolated development environment with a single command, including PostgreSQL and a Valkey (or Redis) cache.  
* **Developer Experience:** A suite of helper scripts or make commands for executing common tasks, such as running the full test suite, linting the code for style consistency, building a production binary, and generating API documentation.  
* **Comprehensive Testing Strategy:** The template must demonstrate a robust testing approach. This includes:  
  * **Unit Testability:** Using go.uber.org/mock/mockgen with //go:generate annotations on all repositories and usecases interfaces to ensure business logic can be tested in complete isolation.
  * **Repository Testing:** Providing examples of repositories tests that use DATA-DOG/go-sqlmock to simulate database queries and responses, ensuring data logic is correct without requiring a live database connection.
* **Clear Instructions:** The process for adding a new business domain to the project must be simple, repeatable, and thoroughly documented to guide developers.