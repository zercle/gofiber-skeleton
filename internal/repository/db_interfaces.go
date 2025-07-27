//go:generate mockgen -source=db_interfaces.go -destination=mocks/mock_db_queries_interface.go -package=mocks DBQueriesInterface

package repository

import (
	"context"
	"gofiber-skeleton/internal/repository/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// DBQueriesInterface defines the methods from db.Queries that are used by repositories.
// This interface is used for mocking the database layer in tests.
type DBQueriesInterface interface {
	CreateURL(ctx context.Context, arg db.CreateURLParams) (db.Url, error)
	DeleteURL(ctx context.Context, id pgtype.UUID) error
	GetURLByShortCode(ctx context.Context, shortCode string) (db.Url, error)
	GetURLsByUserID(ctx context.Context, userID pgtype.UUID) ([]db.Url, error)
	UpdateURL(ctx context.Context, arg db.UpdateURLParams) (db.Url, error)
	GetURLByID(ctx context.Context, id pgtype.UUID) (db.Url, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
}