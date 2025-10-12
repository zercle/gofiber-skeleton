# **Tasks Documentation**

## **sqlc Migration and Implementation Tasks**

This document outlines the tasks required to migrate from raw SQL to sqlc generated code as the primary data access layer, implementing transaction state management and data aggregation in the repository layer.

## **Task 1: sqlc Code Generation Setup**

### **Description**
Generate initial sqlc code and validate the configuration works correctly with the existing database schema.

### **Prerequisites**
- Database migrations are up to date
- sqlc configuration is properly set up in `sqlc.yaml`
- Go dependencies are installed

### **Steps**
1. **Generate sqlc Code**
   ```bash
   make sqlc
   ```

2. **Validate Generated Code**
   - Check that `pkg/database/` directory contains generated files
   - Verify `db.go`, `models.go`, and `queries.sql.go` are created
   - Ensure no compilation errors in generated code

3. **Review Generated Interfaces**
   - Examine generated `Queries` interface
   - Verify all query methods are properly typed
   - Check that model structs match database schema

4. **Update Makefile**
   - Ensure `make sqlc` is included in the development workflow
   - Add sqlc generation to CI/CD pipeline

### **Expected Outcome**
- Successfully generated sqlc code in `pkg/database/`
- All generated code compiles without errors
- Clear understanding of generated interfaces and models

### **Validation Criteria**
- [ ] `pkg/database/` directory exists with generated files
- [ ] `go build` succeeds after code generation
- [ ] All SQL queries from `db/queries/` are represented in generated code

---

## **Task 2: Repository Layer Migration**

### **Description**
Migrate the user repository implementation from raw SQL to sqlc generated code while maintaining the same interface.

### **Affected Files**
- `internal/domains/user/repository/user_impl.go`
- `internal/domains/user/repository/user.go`

### **Steps**
1. **Analyze Current Repository Implementation**
   - Review existing methods in `user_impl.go`
   - Identify all raw SQL queries
   - Map current methods to corresponding sqlc generated methods

2. **Update Repository Structure**
   ```go
   type userRepository struct {
       db      *database.DB
       queries *database.Queries // Add sqlc generated queries
   }
   ```

3. **Migrate Basic CRUD Operations**
   - `GetByID`: Use `queries.GetUserByID`
   - `GetByEmail`: Use `queries.GetUserByEmail`
   - `Create`: Use `queries.CreateUser`
   - `Update`: Use `queries.UpdateUser`
   - `Deactivate`: Use `queries.DeactivateUser`
   - `List`: Use `queries.ListUsers`

4. **Update Constructor**
   ```go
   func NewUserRepository(injector do.Injector) (UserRepository, error) {
       db := do.MustInvoke[*database.Database](injector)
       return &userRepository{
           db:      db.GetDB(),
           queries: database.New(db.GetDB()), // Initialize sqlc queries
       }, nil
   }
   ```

5. **Add Mapping Functions**
   - Create helper functions to convert between domain entities and database models
   - Handle field mappings and data type conversions

### **Expected Outcome**
- User repository fully migrated to use sqlc generated code
- All existing functionality preserved
- Improved type safety and compile-time validation

### **Validation Criteria**
- [ ] All repository methods use sqlc generated queries
- [ ] Repository interface remains unchanged
- [ ] All existing tests pass without modification
- [ ] No runtime SQL errors

---

## **Task 3: Transaction Management Implementation**

### **Description**
Implement transaction state management in the repository layer for operations that require multiple database operations.

### **Affected Files**
- `internal/domains/user/repository/user_impl.go`
- `pkg/database/database.go`

### **Steps**
1. **Enhance Database Transaction Support**
   ```go
   // In pkg/database/database.go
   func (d *Database) WithTx(tx *sql.Tx) *database.Queries {
       return database.NewWithTx(tx)
   }
   ```

2. **Implement Transaction Pattern in Repository**
   ```go
   func (r *userRepository) CreateWithProfile(ctx context.Context, user *entity.User, profile *entity.Profile) error {
       tx, err := r.db.BeginTxx(ctx, nil)
       if err != nil {
           return err
       }
       defer tx.Rollback()
       
       qtx := r.queries.WithTx(tx)
       
       // Execute operations within transaction
       if err := qtx.CreateUser(ctx, toDBUser(user)); err != nil {
           return err
       }
       
       if err := qtx.CreateProfile(ctx, toDBProfile(profile)); err != nil {
           return err
       }
       
       return tx.Commit()
   }
   ```

3. **Add Transaction Helper Methods**
   ```go
   func (r *userRepository) WithTransaction(ctx context.Context, fn func(*database.Queries) error) error {
       tx, err := r.db.BeginTxx(ctx, nil)
       if err != nil {
           return err
       }
       defer tx.Rollback()
       
       if err := fn(r.queries.WithTx(tx)); err != nil {
           return err
       }
       
       return tx.Commit()
   }
   ```

4. **Update Repository Interface**
   - Add methods that require transaction support
   - Document transaction boundaries and behavior

### **Expected Outcome**
- Robust transaction management in repository layer
- Consistent transaction patterns across all repositories
- Proper error handling and rollback mechanisms

### **Validation Criteria**
- [ ] Transactions properly commit on success
- [ ] Transactions properly rollback on error
- [ ] Concurrent operations handle conflicts correctly
- [ ] Transaction boundaries are clearly defined

---

## **Task 4: Data Aggregation Implementation**

### **Description**
Implement data aggregation patterns in the repository layer for complex queries and reporting operations.

### **Affected Files**
- `db/queries/user.sql` (add new aggregation queries)
- `internal/domains/user/repository/user_impl.go`
- `internal/domains/user/repository/user.go`

### **Steps**
1. **Add Aggregation Queries to SQL**
   ```sql
   -- name: GetUserStats :one
   SELECT 
       COUNT(*) as total_users,
       COUNT(CASE WHEN created_at > NOW() - INTERVAL '30 days' THEN 1 END) as new_users,
       COUNT(CASE WHEN is_active = true THEN 1 END) as active_users
   FROM users;
   
   -- name: GetUsersWithActivity :many
   SELECT 
       u.id, u.email, u.full_name, u.created_at,
       COUNT(p.id) as post_count,
       MAX(p.created_at) as last_post_date
   FROM users u
   LEFT JOIN posts p ON u.id = p.user_id
   WHERE u.is_active = true
   GROUP BY u.id, u.email, u.full_name, u.created_at
   ORDER BY post_count DESC;
   ```

2. **Regenerate sqlc Code**
   ```bash
   make sqlc
   ```

3. **Implement Aggregation Methods**
   ```go
   func (r *userRepository) GetUserStats(ctx context.Context) (*entity.UserStats, error) {
       stats, err := r.queries.GetUserStats(ctx)
       if err != nil {
           return nil, err
       }
       return &entity.UserStats{
           TotalUsers:  stats.TotalUsers,
           NewUsers:    stats.NewUsers,
           ActiveUsers: stats.ActiveUsers,
       }, nil
   }
   
   func (r *userRepository) GetUsersWithActivity(ctx context.Context) ([]*entity.UserWithActivity, error) {
       users, err := r.queries.GetUsersWithActivity(ctx)
       if err != nil {
           return nil, err
       }
       return mapToUsersWithActivity(users), nil
   }
   ```

4. **Add Pagination Support for Aggregations**
   ```sql
   -- name: GetUsersWithActivityPaginated :many
   SELECT 
       u.id, u.email, u.full_name, u.created_at,
       COUNT(p.id) as post_count,
       MAX(p.created_at) as last_post_date
   FROM users u
   LEFT JOIN posts p ON u.id = p.user_id
   WHERE u.is_active = true
   GROUP BY u.id, u.email, u.full_name, u.created_at
   ORDER BY post_count DESC
   LIMIT $1 OFFSET $2;
   ```

### **Expected Outcome**
- Comprehensive data aggregation capabilities
- Efficient complex queries with proper indexing
- Consistent pagination patterns for aggregated data

### **Validation Criteria**
- [ ] Aggregation queries return expected results
- [ ] Performance is acceptable for large datasets
- [ ] Pagination works correctly with aggregated data
- [ ] Error handling is robust for complex queries

---

## **Task 5: Testing Strategy Update**

### **Description**
Update the testing strategy to work with sqlc generated code and maintain high test coverage.

### **Affected Files**
- `internal/domains/user/tests/user_test.go`
- Test fixtures and mocks

### **Steps**
1. **Update Mock Generation**
   ```bash
   make mocks
   ```

2. **Create sqlc Test Utilities**
   ```go
   // testutils/database.go
   func NewTestDB(t *testing.T) *sql.DB {
       db := sqltest.NewDB(t, sqltest.WithMigrations(migrations))
       return db
   }
   
   func NewTestQueries(t *testing.T) *database.Queries {
       db := NewTestDB(t)
       return database.New(db)
   }
   ```

3. **Update Repository Tests**
   ```go
   func TestUserRepository_GetByID(t *testing.T) {
       db := NewTestDB(t)
       queries := NewTestQueries(t)
       repo := &userRepository{db: db, queries: queries}
       
       // Setup test data
       user := &entity.User{
           Email:    "test@example.com",
           FullName: "Test User",
       }
       
       err := repo.Create(context.Background(), user)
       require.NoError(t, err)
       
       // Test the method
       found, err := repo.GetByID(context.Background(), user.ID)
       require.NoError(t, err)
       assert.Equal(t, user.Email, found.Email)
   }
   ```

4. **Add Transaction Tests**
   ```go
   func TestUserRepository_TransactionRollback(t *testing.T) {
       db := NewTestDB(t)
       queries := NewTestQueries(t)
       repo := &userRepository{db: db, queries: queries}
       
       // Test transaction rollback
       err := repo.WithTransaction(context.Background(), func(q *database.Queries) error {
           // Create user
           user := toDBUser(&entity.User{Email: "test@example.com"})
           if err := q.CreateUser(context.Background(), user); err != nil {
               return err
           }
           
           // Force an error to test rollback
           return errors.New("intentional error")
       })
       
       assert.Error(t, err)
       
       // Verify user was not created due to rollback
       _, err = repo.GetByEmail(context.Background(), "test@example.com")
       assert.Error(t, err)
       assert.Equal(t, entity.ErrUserNotFound, err)
   }
   ```

### **Expected Outcome**
- Comprehensive test coverage for sqlc-based repository
- Reliable transaction testing
- Mock generation for all interfaces

### **Validation Criteria**
- [ ] All repository tests pass with sqlc implementation
- [ ] Transaction tests cover commit and rollback scenarios
- [ ] Mock generation works correctly
- [ ] Test coverage remains at 90%+

---

## **Task 6: Performance Optimization**

### **Description**
Optimize the sqlc-based implementation for performance and ensure efficient database operations.

### **Affected Files**
- `db/queries/user.sql`
- `internal/domains/user/repository/user_impl.go`
- Database migration files for indexes

### **Steps**
1. **Add Database Indexes**
   ```sql
   -- In migration files
   CREATE INDEX idx_users_email_active ON users(email, is_active);
   CREATE INDEX idx_users_created_at ON users(created_at);
   CREATE INDEX idx_users_active_created ON users(is_active, created_at);
   ```

2. **Optimize SQL Queries**
   - Use appropriate indexes in WHERE clauses
   - Optimize JOIN operations
   - Use proper data types
   - Add query hints if necessary

3. **Implement Connection Pooling**
   ```go
   // In database setup
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(25)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

4. **Add Query Performance Monitoring**
   ```go
   func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
       start := time.Now()
       defer func() {
           duration := time.Since(start)
           // Log query performance
           log.Printf("GetByID query took %v", duration)
       }()
       
       user, err := r.queries.GetUserByID(ctx, id)
       if err != nil {
           return nil, r.handleError(err)
       }
       return r.toDomainEntity(user), nil
   }
   ```

### **Expected Outcome**
- Optimized database performance
- Efficient query execution
- Proper resource utilization

### **Validation Criteria**
- [ ] Query execution times are within acceptable limits
- [ ] Database indexes are properly utilized
- [ ] Connection pooling is effective
- [ ] Memory usage is optimized

---

## **Task 7: Documentation and Examples**

### **Description**
Create comprehensive documentation and examples for using sqlc in the template.

### **Affected Files**
- `docs/SQLC_USAGE.md`
- `docs/ADDING_NEW_DOMAIN.md` (update for sqlc)
- Code examples and tutorials

### **Steps**
1. **Create sqlc Usage Guide**
   ```markdown
   # sqlc Usage Guide
   
   ## Overview
   This template uses sqlc as the primary data access layer for type-safe database operations.
   
   ## Adding New Queries
   1. Add SQL queries to `db/queries/[domain].sql`
   2. Generate code with `make sqlc`
   3. Use generated code in repository implementation
   
   ## Transaction Management
   Transactions are handled at the repository layer...
   ```

2. **Update Domain Addition Guide**
   - Include sqlc configuration steps
   - Show repository implementation patterns
   - Provide transaction management examples

3. **Create Code Examples**
   - Basic CRUD operations with sqlc
   - Transaction patterns
   - Data aggregation examples
   - Testing patterns

4. **Add Best Practices**
   - Query organization
   - Naming conventions
   - Performance considerations
   - Error handling patterns

### **Expected Outcome**
- Comprehensive documentation for sqlc usage
- Clear examples and patterns
- Easy onboarding for new developers

### **Validation Criteria**
- [ ] Documentation is complete and accurate
- [ ] Examples are tested and functional
- [ ] Best practices are clearly defined
- [ ] New developers can follow the guide successfully

---

## **Task 8: CI/CD Integration**

### **Description**
Integrate sqlc code generation and validation into the CI/CD pipeline.

### **Affected Files**
- `.github/workflows/ci.yml`
- `.github/workflows/go-ci.yml`

### **Steps**
1. **Add sqlc Generation to CI**
   ```yaml
   - name: Generate sqlc code
     run: make sqlc
   
   - name: Check for changes
     run: |
       if [[ -n $(git status --porcelain) ]]; then
         echo "sqlc generated code has changed"
         git diff
         exit 1
       fi
   ```

2. **Add sqlc Validation**
   ```yaml
   - name: Validate sqlc queries
     run: |
       sqlc compile
       go build ./pkg/database/...
   ```

3. **Update Test Workflow**
   - Ensure tests run with latest generated code
   - Add performance regression tests
   - Include database migration testing

### **Expected Outcome**
- Automated sqlc validation in CI/CD
- Consistent code generation across environments
- Early detection of SQL query issues

### **Validation Criteria**
- [ ] CI/CD pipeline runs sqlc generation
- [ ] Generated code is validated
- [ ] Tests pass with generated code
- [ ] Performance is monitored

---

## **Implementation Timeline**

### **Phase 1: Foundation (Days 1-2)**
- Task 1: sqlc Code Generation Setup
- Task 2: Repository Layer Migration (Basic CRUD)

### **Phase 2: Advanced Features (Days 3-4)**
- Task 3: Transaction Management Implementation
- Task 4: Data Aggregation Implementation

### **Phase 3: Quality Assurance (Days 5-6)**
- Task 5: Testing Strategy Update
- Task 6: Performance Optimization

### **Phase 4: Documentation and Integration (Days 7-8)**
- Task 7: Documentation and Examples
- Task 8: CI/CD Integration

## **Success Metrics**

### **Technical Metrics**
- 100% repository methods using sqlc generated code
- 90%+ test coverage maintained
- Zero runtime SQL errors
- Query performance within acceptable limits

### **Developer Experience Metrics**
- Clear documentation and examples
- Easy onboarding for new developers
- Consistent patterns across domains
- Improved IDE support and autocomplete

### **Quality Metrics**
- All tests passing
- CI/CD pipeline stable
- Code generation automated
- Performance benchmarks met

## **Rollback Plan**

If issues arise during migration:

1. **Immediate Rollback**
   - Revert to previous repository implementation
   - Disable sqlc code generation
   - Continue with raw SQL implementation

2. **Partial Rollback**
   - Keep sqlc for new domains
   - Maintain raw SQL for existing domains
   - Gradual migration approach

3. **Issue Resolution**
   - Fix specific sqlc issues
   - Update configuration
   - Retry migration with fixes

## **Resources**

### **Documentation**
- [sqlc Official Documentation](https://docs.sqlc.dev/)
- [sqlc Go Tutorial](https://docs.sqlc.dev/en/latest/tutorials/getting-started-go.html)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Performance_Optimizations)

### **Tools**
- sqlc CLI for code generation
- PostgreSQL client for query testing
- Go IDE plugins for sqlc support

### **Support**
- sqlc GitHub repository for issues
- Community forums for best practices
- Internal documentation for patterns