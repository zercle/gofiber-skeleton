package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database/sqlc"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

type userRepository struct {
	queries Querier
}

type Querier interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	ListUsers(ctx context.Context, arg sqlc.ListUsersParams) ([]sqlc.User, error)
}

func NewUserRepository(queries Querier) UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	var idBytes [16]byte
	copy(idBytes[:], user.ID[:])

	params := sqlc.CreateUserParams{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		CreatedAt:    pgtype.Timestamptz{Time: user.CreatedAt, Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt, Valid: true},
	}
	_, err := r.queries.CreateUser(ctx, params)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	dbUser, err := r.queries.GetUserByID(ctx, pgID)
	if err != nil {
		return nil, err
	}
	return r.dbUserToEntity(&dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.dbUserToEntity(&dbUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	var idBytes [16]byte
	copy(idBytes[:], user.ID[:])

	params := sqlc.UpdateUserParams{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt, Valid: true},
	}
	_, err := r.queries.UpdateUser(ctx, params)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	return r.queries.DeleteUser(ctx, pgID)
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	params := sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	dbUsers, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.dbUserToEntity(&dbUser)
	}
	return users, nil
}

func (r *userRepository) dbUserToEntity(dbUser *sqlc.User) *entity.User {
	userID, _ := uuid.FromBytes(dbUser.ID.Bytes[:])
	return &entity.User{
		ID:           userID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		FullName:     dbUser.FullName,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}
}