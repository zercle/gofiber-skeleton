# **Essential Tasks: Go Fiber Skeleton**

*Template eliminates 80-90% of initial setup work - these tasks focus on the 10-20% of business logic development*

## **1. Project Setup Tasks**

### **Task: Initialize New Project from Template**

**Description:** Clone and customize template for new project.

**Files to Modify:**
- `go.mod` (module name)
- `.env.example` (project-specific config)
- `README.md` (project documentation)
- Docker configuration files

**Steps:**
1. Clone template repository
2. Update module name in `go.mod`
3. Update project name in config files
4. Customize `.env.example` with project variables
5. Update README.md with project information
6. Initialize git repository
7. Run `make dev` to verify setup

**Prerequisites:** Go 1.25.0+, Docker, Docker Compose

---

### **Task: Start Development Environment**

**Description:** Set up local development with all services.

**Steps:**
1. Copy `.env.example` to `.env`
2. Configure database and cache settings
3. Start services: `make dev`
4. Run migrations: `make migrate-up`
5. Generate docs: `make docs`
6. Verify all services running

**Validation:** API accessible at `/swagger/`, database connected

---

## **2. Domain Development Tasks**

### **Task: Add New Business Domain**

**Description:** Create new domain following established patterns with DI and routing.

**Files to Create:**
- `internal/domains/{domain}/` (domain structure)
- `internal/domains/{domain}/router/` (domain routing)
- `db/migrations/` (domain migrations)
- `db/queries/` (domain queries)

**Steps:**
1. Create domain directory structure with router subdirectory
2. Implement entity layer with models
3. Create repository interface and SQL queries
4. Implement repository with sqlc code
5. Create usecase layer with business logic
6. Implement handler layer with HTTP endpoints
7. Create domain router implementing `shared/router.Router` interface
8. Register domain services in DI container
9. Add domain router to shared router list in DI container
10. Write comprehensive tests
11. Update API documentation

**Critical Requirements:**
- Follow user/post domain patterns exactly
- Maintain 90%+ test coverage
- Use DI container for all dependencies
- Implement self-registering domain router
- Include Swagger documentation

---

### **Task: Create Database Migration**

**Description:** Add database schema changes for new domain.

**Steps:**
1. Create migration: `make migrate-create name={migration_name}`
2. Write up migration with schema changes
3. Write down migration for rollback
4. Apply migration: `make migrate-up`
5. Generate sqlc queries: `make sqlc`
6. Verify schema changes

**Best Practices:**
- Migrations must be reversible
- Test on development database first
- Consider foreign key relationships

---

### **Task: Add New API Endpoint**

**Description:** Add endpoint to existing domain.

**Steps:**
1. Define usecase interface and implementation
2. Create handler with Swagger annotations
3. Add route to domain router
4. Write unit tests for usecase and handler
5. Update API documentation
6. Test endpoint manually

**Requirements:**
- Follow error handling patterns
- Include proper input validation
- Update Swagger documentation
- Consider authentication needs

---

## **3. Testing Tasks**

### **Task: Write Domain Tests**

**Description:** Create comprehensive tests for new domain.

**Steps:**
1. Generate mocks: `make generate-mocks`
2. Write unit tests for usecase layer
3. Write unit tests for handler layer
4. Write integration tests for repository
5. Run tests: `make test`
6. Verify coverage: `make test-coverage`

**Requirements:**
- Mock all external dependencies
- Test success and error scenarios
- Achieve 90%+ code coverage
- Use table-driven tests

---

### **Task: Run Quality Checks**

**Description:** Execute complete quality validation.

**Steps:**
1. Run unit tests: `make test`
2. Run integration tests: `make test-integration`
3. Generate coverage: `make test-coverage`
4. Run linting: `make lint`
5. Run security checks: `make security`
6. Review and fix any issues

**Requirements:**
- All tests must pass
- Coverage must remain 90%+
- All linting issues resolved
- Security vulnerabilities addressed

---

## **4. Database Tasks**

### **Task: Update Database Schema**

**Description:** Modify database schema with new tables/changes.

**Steps:**
1. Create new migration file
2. Write schema changes in up migration
3. Write rollback in down migration
4. Update entity models if needed
5. Update SQL queries
6. Regenerate sqlc code: `make sqlc`
7. Apply migration: `make migrate-up`
8. Test changes

**Considerations:**
- Test migrations on development database
- Ensure backward compatibility when possible
- Update all related code after schema changes

---

## **5. Build & Deployment Tasks**

### **Task: Build for Production**

**Description:** Create optimized production binary.

**Steps:**
1. Set production environment variables
2. Run tests: `make test`
3. Build binary: `make build`
4. Verify binary functionality
5. Check binary size and dependencies
6. Prepare for deployment

**Requirements:**
- Use production configuration
- Optimize for performance
- Include all necessary assets

---

### **Task: Docker Deployment**

**Description:** Deploy application using containers.

**Steps:**
1. Build Docker image: `make docker-build`
2. Test image locally: `make docker-run`
3. Push to registry: `make docker-push`
4. Deploy to target environment
5. Verify deployment health
6. Monitor application logs

**Best Practices:**
- Use multi-stage builds
- Configure health checks
- Set resource limits
- Manage secrets securely

---

## **6. Maintenance Tasks**

### **Task: Update Dependencies**

**Description:** Update project dependencies to latest versions.

**Steps:**
1. Check outdated dependencies: `make deps-outdated`
2. Update dependencies: `make deps-update`
3. Run tests: `make test`
4. Fix any breaking changes
5. Update documentation if needed
6. Commit changes

**Considerations:**
- Review changelogs for breaking changes
- Test thoroughly after updates
- Consider security implications

---

### **Task: Update API Documentation**

**Description:** Refresh documentation for new/modified endpoints.

**Steps:**
1. Update Swagger annotations in handlers
2. Generate documentation: `make docs`
3. Review generated documentation
4. Update README with changes
5. Test documentation accessibility
6. Commit documentation updates

**Requirements:**
- Keep documentation synchronized with code
- Include clear examples for all endpoints
- Document error responses

---

## **7. Troubleshooting Task**

### **Task: Debug Application Issues**

**Description:** Systematically resolve application problems.

**Steps:**
1. Review application logs
2. Reproduce issue consistently
3. Enable debug logging if needed
4. Isolate problematic component
5. Add temporary logging for debugging
6. Fix identified issue
7. Test fix thoroughly
8. Remove temporary debug code

**Best Practices:**
- Keep detailed debugging notes
- Use systematic approach to isolation
- Consider edge cases and error conditions
- Document root causes and solutions