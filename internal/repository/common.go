package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// parseUUID converts a string to UUID, returning an error if invalid
func parseUUID(id string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %s", id)
	}
	return parsed, nil
}

// sqlNullStringToString converts sql.NullString to string
func sqlNullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// stringToSqlNullString converts string to sql.NullString
func stringToSqlNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}