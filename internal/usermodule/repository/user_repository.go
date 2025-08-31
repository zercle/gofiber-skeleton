package userrepository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/usermodule"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

type userRepository struct {
	q     *sqlc.Queries // The generated Queries struct (holds methods)
	rawDB *sql.DB       // The underlying DB connection (passed as DBTX)
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) usermodule.UserRepository {
	return &userRepository{
		q:     sqlc.New(), // Call the parameterless New()
		rawDB: db,         // Store the actual DB connection
	}
}

func (r *userRepository) Create(user *usermodule.User) (*usermodule.User, error) {
	ctx := context.Background()

	dbUser, err := r.q.CreateUser(ctx, r.rawDB, sqlc.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return nil, err
	}

	user.ID = dbUser.ID.String()
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time
	return user, nil
}

func (r *userRepository) GetByID(id string) (*usermodule.User, error) {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbUser, err := r.q.GetUserByID(ctx, r.rawDB, parsedUUID)
	if err != nil {
		return nil, err
	}

	return &usermodule.User{
		ID:           dbUser.ID.String(),
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByUsername(username string) (*usermodule.User, error) {
	ctx := context.Background()

	dbUser, err := r.q.GetUserByUsername(ctx, r.rawDB, username)
	if err != nil {
		return nil, err
	}

	return &usermodule.User{
		ID:           dbUser.ID.String(),
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Update(user *usermodule.User) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return err
	}

	dbUser, err := r.q.UpdateUser(ctx, r.rawDB, sqlc.UpdateUserParams{
		ID:           parsedUUID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return err
	}

	user.UpdatedAt = dbUser.UpdatedAt.Time
	return nil
}

func (r *userRepository) Delete(id string) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.q.DeleteUser(ctx, r.rawDB, parsedUUID)
}
