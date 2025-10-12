# **Common Tasks and Workflows: Go Fiber Skeleton**

## **1. Project Setup and Initialization**

### **Task: New Project Setup from Template**

**Description:** Initialize a new project from the Go Fiber Skeleton template with proper customization.

**Affected Files:**
- `go.mod` (module name)
- `.env.example` (project-specific configuration)
- `README.md` (project documentation)
- Docker configuration files

**Step-by-Step Procedure:**
1. Clone the template repository
2. Update module name in `go.mod`
3. Update project name in configuration files
4. Customize `.env.example` with project-specific variables
5. Update README.md with project information
6. Initialize git repository
7. Run `make dev` to verify setup

**Dependencies and Considerations:**
- Go 1.25.0+ required
- Docker and Docker Compose required
- Ensure all environment variables are properly configured
- Verify database connectivity before proceeding

---

### **Task: Development Environment Setup**

**Description:** Set up local development environment with all required services.

**Affected Files:**
- `compose.yml` (development services)
- `.env` (local configuration)
- Makefile (development commands)

**Step-by-Step Procedure:**
1. Copy `.env.example` to `.env`
2. Configure database and cache settings
3. Start development services: `make dev`
4. Run database migrations: `make migrate-up`
5. Generate API documentation: `make docs`
6. Verify all services are running

**Dependencies and Considerations:**
- PostgreSQL and Valkey services must be healthy
- Database migrations must be applied before testing
- API documentation should be accessible at `/swagger/`

---

## **2. Domain Development**

### **Task: Add New Business Domain**

**Description:** Create a new business domain following the established patterns.

**Affected Files:**
- `internal/domains/{domain}/` (new domain structure)
- `db/migrations/` (domain migrations)
- `db/queries/` (domain queries)
- Router configuration

**Step-by-Step Procedure:**
1. Create domain directory structure
2. Implement entity layer with models and validation
3. Create repository interface and SQL queries
4. Implement repository with sqlc-generated code
5. Create usecase layer with business logic
6. Implement handler layer with HTTP endpoints
7. Add domain routes to router
8. Write comprehensive tests
9. Update API documentation

**Dependencies and Considerations:**
- Follow user domain patterns exactly
- Ensure proper dependency injection setup
- Maintain test coverage above 90%
- Update documentation with new endpoints

---

### **Task: Database Migration for Domain**

**Description:** Create and apply database migrations for new domain entities.

**Affected Files:**
- `db/migrations/` (migration files)
- Database schema

**Step-by-Step Procedure:**
1. Create migration file: `make migrate-create name={migration_name}`
2. Write up migration with schema changes
3. Write down migration for rollback
4. Apply migration: `make migrate-up`
5. Generate sqlc queries: `make sqlc`
6. Verify schema changes

**Dependencies and Considerations:**
- Migrations must be reversible
- Test migrations on development database first
- Ensure sqlc generation succeeds after schema changes
- Consider foreign key relationships

---

### **Task: Add New API Endpoint**

**Description:** Add new API endpoint to existing domain.

**Affected Files:**
- `internal/domains/{domain}/handler/` (new handler)
- `internal/domains/{domain}/usecase/` (new usecase)
- Router configuration
- API documentation

**Step-by-Step Procedure:**
1. Define usecase interface and implementation
2. Create handler with Swagger annotations
3. Add route to domain router
4. Write unit tests for usecase and handler
5. Update API documentation
6. Test endpoint manually

**Dependencies and Considerations:**
- Follow established error handling patterns
- Ensure proper input validation
- Update Swagger documentation with examples
- Consider authentication requirements

---

## **3. Testing and Quality Assurance**

### **Task: Write Domain Tests**

**Description:** Create comprehensive tests for new domain or functionality.

**Affected Files:**
- `internal/domains/{domain}/` (test files)
- Test configuration

**Step-by-Step Procedure:**
1. Create mock interfaces: `make generate-mocks`
2. Write unit tests for usecase layer
3. Write unit tests for handler layer
4. Write integration tests for repository layer
5. Run tests: `make test`
6. Verify coverage: `make test-coverage`

**Dependencies and Considerations:**
- Mock all external dependencies
- Test both success and error scenarios
- Achieve minimum 90% code coverage
- Use table-driven tests for multiple scenarios

---

### **Task: Run Full Test Suite**

**Description:** Execute complete test suite with coverage and quality checks.

**Affected Files:**
- All source files
- Test configuration

**Step-by-Step Procedure:**
1. Run unit tests: `make test`
2. Run integration tests: `make test-integration`
3. Generate coverage report: `make test-coverage`
4. Run linting: `make lint`
5. Run security checks: `make security`
6. Review results and fix issues

**Dependencies and Considerations:**
- All tests must pass before deployment
- Coverage should remain above 90%
- Linting issues must be resolved
- Security vulnerabilities must be addressed

---

## **4. Database Operations**

### **Task: Database Schema Update**

**Description:** Update database schema with new tables or modifications.

**Affected Files:**
- `db/migrations/` (new migration)
- `db/queries/` (updated queries)
- Entity models

**Step-by-Step Procedure:**
1. Create new migration file
2. Write schema changes in up migration
3. Write rollback in down migration
4. Update entity models if needed
5. Update SQL queries
6. Regenerate sqlc code: `make sqlc`
7. Apply migration: `make migrate-up`
8. Test changes

**Dependencies and Considerations:**
- Test migrations on development database first
- Ensure backward compatibility when possible
- Update all related code after schema changes
- Consider data migration for existing records

---

### **Task: Database Backup and Restore**

**Description:** Backup and restore database for development or deployment.

**Affected Files:**
- Database data
- Backup files

**Step-by-Step Procedure:**
1. Stop application services
2. Create database backup: `make db-backup`
3. Store backup securely
4. Restore from backup: `make db-restore {backup_file}`
5. Verify data integrity
6. Restart application services

**Dependencies and Considerations:**
- Ensure database is quiescent during backup
- Verify backup file integrity
- Test restore process on non-production database
- Consider encryption for sensitive data backups

---

## **5. Build and Deployment**

### **Task: Build Production Binary**

**Description:** Build optimized production binary for deployment.

**Affected Files:**
- Binary output
- Build configuration

**Step-by-Step Procedure:**
1. Set production environment variables
2. Run tests: `make test`
3. Build binary: `make build`
4. Verify binary functionality
5. Check binary size and dependencies
6. Prepare for deployment

**Dependencies and Considerations:**
- Use production configuration
- Optimize binary size and performance
- Include all necessary assets
- Test binary in production-like environment

---

### **Task: Docker Deployment**

**Description:** Deploy application using Docker containers.

**Affected Files:**
- Docker image
- Container configuration

**Step-by-Step Procedure:**
1. Build Docker image: `make docker-build`
2. Test image locally: `make docker-run`
3. Push to registry: `make docker-push`
4. Deploy to target environment
5. Verify deployment health
6. Monitor application logs

**Dependencies and Considerations:**
- Use multi-stage builds for optimization
- Configure proper health checks
- Set appropriate resource limits
- Manage secrets securely

---

## **6. Maintenance and Updates**

### **Task: Dependency Updates**

**Description:** Update project dependencies to latest stable versions.

**Affected Files:**
- `go.mod` and `go.sum`
- Potentially affected code

**Step-by-Step Procedure:**
1. Check for outdated dependencies: `make deps-outdated`
2. Update dependencies: `make deps-update`
3. Run tests: `make test`
4. Fix any breaking changes
5. Update documentation if needed
6. Commit changes

**Dependencies and Considerations:**
- Review changelogs for breaking changes
- Test thoroughly after updates
- Update documentation for API changes
- Consider security implications

---

### **Task: Security Audit**

**Description:** Perform security audit and address vulnerabilities.

**Affected Files:**
- Dependencies
- Configuration
- Code with security implications

**Step-by-Step Procedure:**
1. Run security scan: `make security`
2. Review dependency vulnerabilities
3. Check for common security issues
4. Address identified vulnerabilities
5. Update security configurations
6. Document security improvements

**Dependencies and Considerations:**
- Prioritize critical vulnerabilities
- Test security fixes thoroughly
- Keep security documentation updated
- Consider security best practices

---

## **7. Documentation and Communication**

### **Task: Update API Documentation**

**Description:** Update API documentation for new or modified endpoints.

**Affected Files:**
- Source code comments
- Generated documentation
- README files

**Step-by-Step Procedure:**
1. Update Swagger annotations in handlers
2. Generate documentation: `make docs`
3. Review generated documentation
4. Update README with changes
5. Test documentation accessibility
6. Commit documentation updates

**Dependencies and Considerations:**
- Keep documentation synchronized with code
- Include clear examples for all endpoints
- Document error responses
- Consider API versioning implications

---

### **Task: Performance Optimization**

**Description:** Optimize application performance and resource usage.

**Affected Files:**
- Database queries
- Application code
- Configuration

**Step-by-Step Procedure:**
1. Profile application: `make profile`
2. Identify performance bottlenecks
3. Optimize database queries
4. Implement caching where appropriate
5. Optimize memory usage
6. Benchmark improvements

**Dependencies and Considerations:**
- Measure before and after optimization
- Focus on critical paths first
- Consider caching strategies
- Monitor for regressions

---

## **8. Troubleshooting**

### **Task: Debug Application Issues**

**Description:** Systematically debug and resolve application issues.

**Affected Files:**
- Error logs
- Problematic code
- Configuration

**Step-by-Step Procedure:**
1. Review application logs
2. Reproduce issue consistently
3. Enable debug logging if needed
4. Isolate problematic component
5. Add temporary logging for debugging
6. Fix identified issue
7. Test fix thoroughly
8. Remove temporary debug code

**Dependencies and Considerations:**
- Keep detailed debugging notes
- Use systematic approach to isolation
- Consider edge cases and error conditions
- Document root causes and solutions