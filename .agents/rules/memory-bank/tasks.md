# **Tasks Documentation**

## **sqlc Migration and Implementation Tasks**

This document outlines the tasks required to migrate from raw SQL to sqlc generated code as the primary data access layer, implementing transaction state management and data aggregation in the repository layer, and adding a new post domain.

## **Task 1: Update sqlc Configuration for Domain-Based Generation**

### **Description**
Update the sqlc configuration to generate code in domain entity directories instead of a centralized package.

### **Affected Files**
- `sqlc.yaml`

### **Steps**
1. **Update sqlc.yaml Configuration**
   ```yaml
   version: "2"
   sql:
     - engine: "postgresql"
       queries: "db/queries/user.sql"
       schema: "db/migrations/"
       gen:
         go:
           package: "entity"
           out: "internal/domains/user/entity"
           sql_package: "pgx/v5"
           emit_json_tags: true
           emit_prepared_queries: false
           emit_interface: true
           emit_exact_table_names: false
     - engine: "postgresql"
       queries: "db/queries/post.sql"
       schema: "db/migrations/"
       gen:
         go:
           package: "entity"
           out: "internal/domains/post/entity"
           sql_package: "pgx/v5"
           emit_json_tags: true
           emit_prepared_queries: false
           emit_interface: true
           emit_exact_table_names: false
   ```

2. **Update Makefile Commands**
   ```makefile
   .PHONY: sqlc sqlc-user sqlc-post
   
   sqlc:
   	sqlc generate
   
   sqlc-user:
   	sqlc generate -f sqlc.yaml --sql-files=user.sql
   
   sqlc-post:
   	sqlc generate -f sqlc.yaml --sql-files=post.sql
   ```

3. **Validate Configuration**
   - Run `make sqlc` to test generation
   - Verify code is generated in correct directories
   - Ensure no conflicts between domains

### **Expected Outcome**
- sqlc code generated in domain-specific entity directories
- Clear separation of generated code by domain
- Updated build commands for domain-specific generation

### **Validation Criteria**
- [ ] User domain code generated in `internal/domains/user/entity/`
- [ ] Post domain code generated in `internal/domains/post/entity/`
- [ ] No code conflicts between domains
- [ ] All generated code compiles successfully

---

## **Task 2: Create Posts Database Migration**

### **Description**
Create database migration for the posts table with proper relationships to users.

### **Affected Files**
- `db/migrations/000002_create_posts_table.up.sql`
- `db/migrations/000002_create_posts_table.down.sql`

### **Steps**
1. **Create Up Migration**
   ```sql
   -- 000002_create_posts_table.up.sql
   CREATE TABLE posts (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
       title VARCHAR(255) NOT NULL,
       content TEXT NOT NULL,
       status VARCHAR(50) NOT NULL DEFAULT 'draft',
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   );
   
   -- Indexes for performance
   CREATE INDEX idx_posts_user_id ON posts(user_id);
   CREATE INDEX idx_posts_status ON posts(status);
   CREATE INDEX idx_posts_created_at ON posts(created_at);
   CREATE INDEX idx_posts_user_status ON posts(user_id, status);
   
   -- Trigger for updated_at
   CREATE OR REPLACE FUNCTION update_updated_at_column()
   RETURNS TRIGGER AS $$
   BEGIN
       NEW.updated_at = NOW();
       RETURN NEW;
   END;
   $$ language 'plpgsql';
   
   CREATE TRIGGER update_posts_updated_at 
       BEFORE UPDATE ON posts 
       FOR EACH ROW 
       EXECUTE FUNCTION update_updated_at_column();
   ```

2. **Create Down Migration**
   ```sql
   -- 000002_create_posts_table.down.sql
   DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
   DROP FUNCTION IF EXISTS update_updated_at_column();
   DROP TABLE IF EXISTS posts;
   ```

3. **Run Migration**
   ```bash
   make migrate-up
   ```

### **Expected Outcome**
- Posts table created with proper relationships
- Indexes optimized for common queries
- Automatic timestamp management

### **Validation Criteria**
- [ ] Migration runs successfully
- [ ] Posts table created with all columns
- [ ] Foreign key relationship to users established
- [ ] Indexes created properly
- [ ] Trigger for updated_at works correctly

---

## **Task 3: Create Post Domain SQL Queries**

### **Description**
Create SQL queries for the post domain following sqlc naming conventions.

### **Affected Files**
- `db/queries/post.sql`

### **Steps**
1. **Create Post SQL Queries**
   ```sql
   -- name: CreatePost :one
   INSERT INTO posts (user_id, title, content, status)
   VALUES ($1, $2, $3, $4)
   RETURNING *;
   
   -- name: GetPostByID :one
   SELECT * FROM posts
   WHERE id = $1
   LIMIT 1;
   
   -- name: UpdatePost :one
   UPDATE posts
   SET title = $2, content = $3, status = $4, updated_at = NOW()
   WHERE id = $1
   RETURNING *;
   
   -- name: DeletePost :exec
   DELETE FROM posts WHERE id = $1;
   
   -- name: GetPostsByUserID :many
   SELECT * FROM posts
   WHERE user_id = $1
   ORDER BY created_at DESC
   LIMIT $2 OFFSET $3;
   
   -- name: ListPosts :many
   SELECT p.*, u.email as user_email, u.full_name as user_full_name
   FROM posts p
   JOIN users u ON p.user_id = u.id
   WHERE p.status = $1
   ORDER BY p.created_at DESC
   LIMIT $2 OFFSET $3;
   
   -- name: CountPostsByUser :one
   SELECT COUNT(*) FROM posts
   WHERE user_id = $1 AND status = 'published';
   
   -- name: GetPostsWithAuthor :many
   SELECT p.*, u.email as user_email, u.full_name as user_full_name
   FROM posts p
   JOIN users u ON p.user_id = u.id
   WHERE p.status = 'published'
   ORDER BY p.created_at DESC
   LIMIT $1 OFFSET $2;
   
   -- name: GetUserPostStats :one
   SELECT 
       COUNT(*) as total_posts,
       COUNT(CASE WHEN status = 'published' THEN 1 END) as published_posts,
       COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft_posts,
       MAX(created_at) as last_post_date
   FROM posts
   WHERE user_id = $1;
   ```

2. **Generate sqlc Code**
   ```bash
   make sqlc-post
   ```

3. **Validate Generated Code**
   - Check `internal/domains/post/entity/` directory
   - Verify all query methods are generated
   - Ensure proper types are generated

### **Expected Outcome**
- Comprehensive SQL queries for post operations
- Type-safe Go code generation
- Support for user-post relationships

### **Validation Criteria**
- [ ] All CRUD operations covered
- [ ] User relationship queries included
- [ ] Aggregation queries for statistics
- [ ] Generated code compiles without errors

---

## **Task 4: Create Post Domain Structure**

### **Description**
Create the complete post domain structure following the user domain pattern.

### **Affected Files**
- `internal/domains/post/entity/post.go`
- `internal/domains/post/repository/post.go`
- `internal/domains/post/repository/post_impl.go`
- `internal/domains/post/usecase/post.go`
- `internal/domains/post/delivery/post.go`
- `internal/domains/post/tests/post_test.go`

### **Steps**
1. **Create Post Entity**
   ```go
   // internal/domains/post/entity/post.go
   package entity
   
   import (
       "time"
       "github.com/google/uuid"
   )
   
   type Post struct {
       ID        uuid.UUID `json:"id" db:"id"`
       UserID    uuid.UUID `json:"user_id" db:"user_id"`
       Title     string    `json:"title" db:"title"`
       Content   string    `json:"content" db:"content"`
       Status    string    `json:"status" db:"status"`
       CreatedAt time.Time `json:"created_at" db:"created_at"`
       UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
   }
   
   type PostWithAuthor struct {
       ID           uuid.UUID `json:"id" db:"id"`
       UserID       uuid.UUID `json:"user_id" db:"user_id"`
       Title        string    `json:"title" db:"title"`
       Content      string    `json:"content" db:"content"`
       Status       string    `json:"status" db:"status"`
       CreatedAt    time.Time `json:"created_at" db:"created_at"`
       UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
       UserEmail    string    `json:"user_email" db:"user_email"`
       UserFullName string    `json:"user_full_name" db:"user_full_name"`
   }
   
   type PostStats struct {
       TotalPosts     int       `json:"total_posts"`
       PublishedPosts int       `json:"published_posts"`
       DraftPosts     int       `json:"draft_posts"`
       LastPostDate   time.Time `json:"last_post_date"`
   }
   ```

2. **Create Repository Interface**
   ```go
   // internal/domains/post/repository/post.go
   package repository
   
   import (
       "context"
       "github.com/google/uuid"
       "github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
   )
   
   type PostRepository interface {
       Create(ctx context.Context, post *entity.Post) error
       GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
       Update(ctx context.Context, post *entity.Post) error
       Delete(ctx context.Context, id uuid.UUID) error
       GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Post, error)
       ListPosts(ctx context.Context, status string, limit, offset int) ([]*entity.PostWithAuthor, error)
       GetPostsWithAuthor(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error)
       CountPostsByUser(ctx context.Context, userID uuid.UUID) (int, error)
       GetUserPostStats(ctx context.Context, userID uuid.UUID) (*entity.PostStats, error)
   }
   ```

3. **Create Repository Implementation**
   ```go
   // internal/domains/post/repository/post_impl.go
   package repository
   
   import (
       "context"
       "database/sql"
       "github.com/google/uuid"
       "github.com/samber/do/v2"
       "github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
       "github.com/zercle/gofiber-skeleton/internal/domains/post/entity" // sqlc generated
       "github.com/zercle/gofiber-skeleton/pkg/database"
   )
   
   type postRepository struct {
       db      *database.DB
       queries *entity.Queries // sqlc generated in same domain
   }
   
   func NewPostRepository(injector do.Injector) (PostRepository, error) {
       db := do.MustInvoke[*database.Database](injector)
       return &postRepository{
           db:      db.GetDB(),
           queries: entity.New(db.GetDB()),
       }, nil
   }
   
   // Implement all repository methods using sqlc generated queries
   ```

4. **Create Use Case Layer**
5. **Create Delivery Layer**
6. **Create Tests**

### **Expected Outcome**
- Complete post domain structure
- Consistent with user domain patterns
- Full CRUD operations with user relationships

### **Validation Criteria**
- [ ] All domain layers created
- [ ] Repository uses sqlc generated code
- [ ] Use cases implement business logic
- [ ] Delivery layer handles HTTP requests
- [ ] Tests provide good coverage

---

## **Task 5: Update User Repository for Domain-Based sqlc**

### **Description**
Update the user repository implementation to use the new domain-based sqlc generated code.

### **Affected Files**
- `internal/domains/user/repository/user_impl.go`

### **Steps**
1. **Update Repository Import**
   ```go
   import (
       "github.com/zercle/gofiber-skeleton/internal/domains/user/entity" // sqlc generated
   )
   ```

2. **Update Repository Structure**
   ```go
   type userRepository struct {
       db      *database.DB
       queries *entity.Queries // sqlc generated in same domain
   }
   ```

3. **Update Constructor**
   ```go
   func NewUserRepository(injector do.Injector) (UserRepository, error) {
       db := do.MustInvoke[*database.Database](injector)
       return &userRepository{
           db:      db.GetDB(),
           queries: entity.New(db.GetDB()), // Initialize from domain entity
       }, nil
   }
   ```

4. **Update All Methods to Use Domain-Based Queries**
   - Update import paths
   - Ensure all method calls use correct package references
   - Verify type compatibility

### **Expected Outcome**
- User repository fully migrated to domain-based sqlc
- Clean import structure within domain boundaries
- Maintained functionality with improved organization

### **Validation Criteria**
- [ ] All repository methods use domain-based sqlc
- [ ] No cross-domain imports for generated code
- [ ] All existing tests pass
- [ ] Code compiles without errors

---

## **Task 6: Implement Cross-Domain Features**

### **Description**
Implement features that require interaction between user and post domains.

### **Affected Files**
- `internal/domains/user/usecase/user.go`
- `internal/domains/post/usecase/post.go`
- `internal/domains/post/repository/post_impl.go`

### **Steps**
1. **Add Post Count to User Profile**
   ```go
   // In user usecase
   func (u *userUsecase) GetProfileWithPostStats(ctx context.Context, userID uuid.UUID) (*entity.UserProfileWithStats, error) {
       // Get user profile
       user, err := u.userRepo.GetByID(ctx, userID)
       if err != nil {
           return nil, err
       }
       
       // Get post repository
       postRepo := do.MustInvoke[repository.PostRepository](u.injector)
       
       // Get post statistics
       stats, err := postRepo.GetUserPostStats(ctx, userID)
       if err != nil {
           return nil, err
       }
       
       return &entity.UserProfileWithStats{
           User:  user.ToResponse(),
           Stats: stats,
       }, nil
   }
   ```

2. **Implement Transaction for User-Post Operations**
   ```go
   // In post repository
   func (r *postRepository) CreatePostWithUserUpdate(ctx context.Context, post *entity.Post) error {
       tx, err := r.db.BeginTxx(ctx, nil)
       if err != nil {
           return err
       }
       defer tx.Rollback()
       
       qtx := r.queries.WithTx(tx)
       
       // Create post
       if err := qtx.CreatePost(ctx, toDBPost(post)); err != nil {
           return err
       }
       
       // Update user's last activity (if needed)
       // This would require user queries in post domain or a shared transaction
       
       return tx.Commit()
   }
   ```

3. **Add Authorization for Post Operations**
   ```go
   // In post usecase
   func (p *postUsecase) UpdatePost(ctx context.Context, userID, postID uuid.UUID, req *entity.UpdatePostRequest) (*entity.PostResponse, error) {
       // Get post
       post, err := p.postRepo.GetByID(ctx, postID)
       if err != nil {
           return nil, err
       }
       
       // Check ownership
       if post.UserID != userID {
           return nil, entity.ErrUnauthorized
       }
       
       // Update post
       post.Title = req.Title
       post.Content = req.Content
       post.Status = req.Status
       
       if err := p.postRepo.Update(ctx, post); err != nil {
           return nil, err
       }
       
       return post.ToResponse(), nil
   }
   ```

### **Expected Outcome**
- Seamless integration between user and post domains
- Proper authorization and ownership checks
- Transaction support for cross-domain operations

### **Validation Criteria**
- [ ] User profiles show post statistics
- [ ] Post ownership verification works
- [ ] Cross-domain transactions function correctly
- [ ] Authorization is properly enforced

---

## **Task 7: Update API Routes and Documentation**

### **Description**
Add post-related API routes and update documentation.

### **Affected Files**
- `cmd/server/main.go` (router setup)
- `internal/domains/post/delivery/post.go`
- Swagger documentation

### **Steps**
1. **Add Post Routes**
   ```go
   // In main.go router setup
   app.Get("/api/v1/posts", postHandler.ListPosts)
   app.Get("/api/v1/posts/:id", postHandler.GetPost)
   app.Post("/api/v1/posts", authMiddleware, postHandler.CreatePost)
   app.Put("/api/v1/posts/:id", authMiddleware, postHandler.UpdatePost)
   app.Delete("/api/v1/posts/:id", authMiddleware, postHandler.DeletePost)
   app.Get("/api/v1/users/:id/posts", postHandler.GetUserPosts)
   ```

2. **Add Swagger Documentation**
   ```go
   // @Summary Create a new post
   // @Description Create a new post for the authenticated user
   // @Tags posts
   // @Accept json
   // @Produce json
   // @Param request body entity.CreatePostRequest true "Post data"
   // @Success 201 {object} entity.PostResponse
   // @Failure 400 {object} response.ErrorResponse
   // @Failure 401 {object} response.ErrorResponse
   // @Router /api/v1/posts [post]
   ```

3. **Generate API Documentation**
   ```bash
   make swag
   ```

### **Expected Outcome**
- Complete API endpoints for post operations
- Updated Swagger documentation
- Proper authentication and authorization

### **Validation Criteria**
- [ ] All post endpoints work correctly
- [ ] Swagger documentation is complete
- [ ] Authentication is enforced where needed
- [ ] API responses are consistent

---

## **Task 8: Update Testing Strategy**

### **Description**
Update testing strategy to work with domain-based sqlc and post domain.

### **Affected Files**
- `internal/domains/user/tests/user_test.go`
- `internal/domains/post/tests/post_test.go`
- Test utilities and fixtures

### **Steps**
1. **Update User Tests**
   ```go
   func TestUserRepository(t *testing.T) {
       db := NewTestDB(t)
       queries := entity.New(db) // From user entity package
       repo := &userRepository{db: db, queries: queries}
       // ... test implementation
   }
   ```

2. **Create Post Tests**
   ```go
   func TestPostRepository(t *testing.T) {
       db := NewTestDB(t)
       queries := entity.New(db) // From post entity package
       repo := &postRepository{db: db, queries: queries}
       // ... test implementation
   }
   ```

3. **Add Integration Tests**
   ```go
   func TestUserPostIntegration(t *testing.T) {
       // Test user-post relationships
       // Test cross-domain operations
       // Test transaction scenarios
   }
   ```

4. **Update Mock Generation**
   ```bash
   make mocks
   ```

### **Expected Outcome**
- Comprehensive test coverage for both domains
- Integration tests for cross-domain features
- Updated mocks for domain-based interfaces

### **Validation Criteria**
- [ ] All repository tests pass
- [ ] Use case tests provide good coverage
- [ ] Integration tests work correctly
- [ ] Mock generation is successful

---

## **Task 9: Update CI/CD Pipeline**

### **Description**
Update CI/CD pipeline to handle domain-based sqlc generation and post domain testing.

### **Affected Files**
- `.github/workflows/ci.yml`
- `.github/workflows/go-ci.yml`

### **Steps**
1. **Update sqlc Generation in CI**
   ```yaml
   - name: Generate sqlc code
     run: |
       make sqlc-user
       make sqlc-post
   
   - name: Check for changes
     run: |
       if [[ -n $(git status --porcelain) ]]; then
         echo "sqlc generated code has changed"
         git diff
         exit 1
       fi
   ```

2. **Add Domain-Specific Tests**
   ```yaml
   - name: Run user domain tests
     run: go test ./internal/domains/user/...
   
   - name: Run post domain tests
     run: go test ./internal/domains/post/...
   
   - name: Run integration tests
     run: go test ./tests/...
   ```

3. **Update Build Process**
   - Ensure all domains are included in build
   - Validate domain-based dependencies
   - Check for cross-domain violations

### **Expected Outcome**
- CI/CD pipeline handles domain-based generation
- All domains tested independently and together
- Automated validation of domain boundaries

### **Validation Criteria**
- [ ] CI/CD pipeline runs successfully
- [ ] All domain tests are executed
- [ ] sqlc generation is validated
- [ ] No cross-domain violations

---

## **Implementation Timeline**

### **Phase 1: Foundation Updates (Days 1-2)**
- Task 1: Update sqlc configuration for domain-based generation
- Task 2: Create posts database migration
- Task 3: Create post domain SQL queries

### **Phase 2: Domain Implementation (Days 3-4)**
- Task 4: Create post domain structure
- Task 5: Update user repository for domain-based sqlc

### **Phase 3: Integration (Days 5-6)**
- Task 6: Implement cross-domain features
- Task 7: Update API routes and documentation

### **Phase 4: Testing and CI/CD (Days 7-8)**
- Task 8: Update testing strategy
- Task 9: Update CI/CD pipeline

## **Success Metrics**

### **Technical Metrics**
- 100% repository methods use domain-based sqlc
- Complete post domain implementation
- 90%+ test coverage across all domains
- Zero cross-domain violations

### **Domain Isolation Metrics**
- Clear domain boundaries maintained
- Minimal cross-domain dependencies
- Domain-specific code generation
- Independent domain testing

### **Feature Metrics**
- User can create, read, update, delete posts
- User profiles show post statistics
- Post ownership and authorization enforced
- Cross-domain transactions work correctly

## **Rollback Plan**

If issues arise during implementation:

1. **Configuration Rollback**
   - Revert sqlc.yaml to centralized generation
   - Update import paths accordingly
   - Regenerate code with old configuration

2. **Domain Rollback**
   - Remove post domain files
   - Revert database migrations
   - Update routes and documentation

3. **Partial Rollback**
   - Keep domain structure with centralized sqlc
   - Migrate gradually to domain-based approach
   - Maintain functionality while fixing issues

## **Resources**

### **Documentation**
- [sqlc Official Documentation](https://docs.sqlc.dev/)
- [Domain-Driven Design Best Practices](https://github.com/ddd-crew/ddd-starter-modern-monolith)
- [Clean Architecture Patterns](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

### **Tools**
- sqlc CLI for domain-specific code generation
- PostgreSQL client for query testing
- Go IDE plugins for sqlc support
- Database migration tools

### **Support**
- sqlc GitHub repository for issues
- Community forums for best practices
- Internal documentation for patterns
- Domain-specific examples and guides