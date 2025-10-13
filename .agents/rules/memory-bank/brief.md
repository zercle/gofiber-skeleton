# **Project Brief: Go Fiber Production-Ready Template**

## **1. Template Purpose**

Production-ready **template repository** providing a complete Go backend foundation using **Fiber v2** with **Domain-Driven Clean Architecture**. Eliminates **80-90% of initial project setup** work, enabling developers to focus on business logic from day one.

## **2. Value Proposition**

### **What Developers Get (80-90% Complete)**
* ✅ Complete Clean Architecture implementation
* ✅ Production-ready configuration management
* ✅ Database integration with migrations + type-safe queries
* ✅ JWT authentication system with password hashing
* ✅ Comprehensive middleware stack
* ✅ Development environment with Docker Compose
* ✅ Testing framework with mocks and examples
* ✅ Auto-generated API documentation
* ✅ CI/CD foundation with health checks
* ✅ Deployment-ready containerization

### **What Developers Add (10-20% Business Logic)**
* Business-specific domains following reference pattern
* Custom validation rules and business logic
* Domain-specific database migrations
* API endpoints for their use cases
* Additional third-party integrations

## **3. Core Technologies**

* **Framework:** Go Fiber v2 (high-performance web framework)
* **Architecture:** Domain-Driven Clean Architecture
* **Dependency Injection:** Samber's do (type-safe, generics-based)
* **Database:** PostgreSQL + pgx driver + sqlc (type-safe SQL)
* **Authentication:** golang-jwt + Argon2id password hashing
* **Configuration:** Viper (environment-aware)
* **Testing:** uber-go/mock + comprehensive examples
* **Documentation:** swaggo/swag (auto-generated Swagger/OpenAPI)
* **Development:** Air (hot-reloading) + Docker Compose

## **4. Reference Implementation**

Complete **user/auth domain** demonstrating all patterns:
* Entity layer with proper domain modeling
* Repository layer with PostgreSQL integration
* Usecase layer with business logic and tests
* Handler layer with HTTP endpoints and docs
* 90%+ test coverage with proper mocking
* Production-ready security practices

## **5. Getting Started Workflow**

### **5-Minute Setup**
1. Clone repository
2. Customize module name and configuration
3. Run `make dev` - everything works out of the box
4. Explore running application and API docs

### **1-Hour Domain Addition**
1. Follow `docs/ADDING_NEW_DOMAIN.md` guide
2. Create domain structure following user pattern
3. Implement business logic with tests
4. Add routes - everything integrates automatically

## **6. Template Benefits**

### **For Development Teams**
* **Time Savings:** Eliminate weeks of initial setup
* **Consistency:** Standardized architecture across projects
* **Quality:** Built-in best practices and code quality
* **Scalability:** Proven patterns that grow with teams

### **For Individual Developers**
* **Speed:** Start building business logic immediately
* **Learning:** Real-world Go best practices examples
* **Confidence:** Production-ready foundation
* **Growth:** Understanding of enterprise architecture

### **For Organizations**
* **Standardization:** Consistent patterns across teams
* **Onboarding:** Faster team member integration
* **Maintenance:** Reduced overhead with proven patterns
* **Innovation:** More time for business logic

## **7. Success Criteria**

A developer should be able to:
* ✅ Clone and run template in **under 5 minutes**
* ✅ Understand architecture by reviewing user domain
* ✅ Add new CRUD domain in **under 1 hour**
* ✅ Deploy to production with minimal configuration
* ✅ Maintain high code quality with built-in tooling
* ✅ Scale application by adding more domains

## **8. Target Audience**

* **Development Teams** starting new Go backend projects
* **Individual Developers** building production applications
* **Organizations** standardizing on Go for microservices
* **Startups** needing rapid development without quality sacrifice
* **Developers** learning Go best practices with real examples

## **9. Quality Standards**

### **Enforced by Template**
* **Architecture:** Strict Clean Architecture with domain isolation
* **Testing:** Mock-based unit testing with 90%+ coverage target
* **Security:** JWT authentication, password hashing, input validation
* **Performance:** Connection pooling, caching, efficient queries
* **Documentation:** Auto-generated Swagger docs for all endpoints
* **Code Quality:** golangci-lint with comprehensive rules
* **Type Safety:** sqlc for compile-time SQL validation

### **Developer Experience**
* No research needed - best practices built-in
* No tooling setup - everything pre-configured
* No boilerplate writing - reference implementation provided
* No deployment configuration - Docker setup included
* No authentication implementation - production-ready auth included

## **10. Template Maintenance**

* Regular dependency updates for security and features
* Best practice improvements based on Go ecosystem evolution
* Community feedback incorporation
* Documentation updates and clarifications
* New feature examples and patterns