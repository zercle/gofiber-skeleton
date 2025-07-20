package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/internal/repository/db"
)

type postgresRepository struct {
	*db.Queries
}

func NewPostgresRepository(q *db.Queries) *postgresRepository {
	return &postgresRepository{q}
}

func (r *postgresRepository) CreateUser(ctx context.Context, username, password string) (*entities.User, error) {
	u7, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	user, err := r.Queries.CreateUser(ctx, db.CreateUserParams{ID: u7, Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	return &entities.User{ID: user.ID.String(), Username: user.Username, CreatedAt: user.CreatedAt.Time}, nil
}

func (r *postgresRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := r.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &entities.User{ID: user.ID.String(), Username: user.Username, Password: user.Password, CreatedAt: user.CreatedAt.Time}, nil
}

func (r *postgresRepository) CreateURL(ctx context.Context, userID *string, shortCode, originalURL string) (*entities.URL, error) {
	u7, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	params := db.CreateURLParams{ID: u7, ShortCode: shortCode, LongUrl: originalURL}
	if userID != nil {
		uuidUserID, err := uuid.FromString(*userID)
		if err != nil {
			return nil, err
		}
		params.UserID = uuid.NullUUID{UUID: uuidUserID, Valid: true}
	}
	url, err := r.Queries.CreateURL(ctx, params)
	if err != nil {
		return nil, err
	}
	var entityUserID *string
	if url.UserID.Valid {
		strUserID := url.UserID.UUID.String()
		entityUserID = &strUserID
	}
	return &entities.URL{ID: url.ID.String(), UserID: entityUserID, ShortCode: url.ShortCode, OriginalURL: url.LongUrl, CreatedAt: url.CreatedAt.Time}, nil
}

func (r *postgresRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error) {
	url, err := r.Queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	var entityUserID *string
	if url.UserID.Valid {
		strUserID := url.UserID.UUID.String()
		entityUserID = &strUserID
	}
	return &entities.URL{ID: url.ID.String(), UserID: entityUserID, ShortCode: url.ShortCode, OriginalURL: url.LongUrl, CreatedAt: url.CreatedAt.Time}, nil
}

func (r *postgresRepository) GetURLsByUserID(ctx context.Context, userID string) ([]*entities.URL, error) {
	uuidUserID, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}
	urls, err := r.Queries.GetURLsByUserID(ctx, uuid.NullUUID{UUID: uuidUserID, Valid: true})
	if err != nil {
		return nil, err
	}
	var result []*entities.URL
	for _, url := range urls {
		var entityUserID *string
		if url.UserID.Valid {
			strUserID := url.UserID.UUID.String()
			entityUserID = &strUserID
		}
		result = append(result, &entities.URL{ID: url.ID.String(), UserID: entityUserID, ShortCode: url.ShortCode, OriginalURL: url.LongUrl, CreatedAt: url.CreatedAt.Time})
	}
	return result, nil
}

func (r *postgresRepository) DeleteURL(ctx context.Context, id, userID string) error {
	uuidID, err := uuid.FromString(id)
	if err != nil {
		return err
	}
	uuidUserID, err := uuid.FromString(userID)
	if err != nil {
		return err
	}
	return r.Queries.DeleteURL(ctx, db.DeleteURLParams{ID: uuidID, UserID: uuid.NullUUID{UUID: uuidUserID, Valid: true}})
}
