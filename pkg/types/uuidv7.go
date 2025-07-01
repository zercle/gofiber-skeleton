package uuidgen

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// UUID is a custom type for UUIDs to ensure type safety.
type UUID uuid.UUID

// NewUUIDv7 creates a new UUID version 7.
func NewUUIDv7() (UUID, error) {
	id, err := uuid.NewV7()
	return UUID(id), err
}

// UUIDFromString parses a string into a UUID.
// func UUIDFromString(s string) (UUID, error) {
// 	id, err := uuid.FromString(s)
// 	return UUID(id), err
// }

// String returns the string representation of the UUID.
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// MarshalJSON implements the json.Marshaler interface.
// func (u UUID) MarshalJSON() ([]byte, error) {
// 	return uuid.UUID(u).MarshalJSON()
// }

// UnmarshalJSON implements the json.Unmarshaler interface.
// func (u *UUID) UnmarshalJSON(b []byte) error {
// 	var id uuid.UUID
// 	if err := id.UnmarshalJSON(b); err != nil {
// 		return err
// 	}
// 	*u = UUID(id)
// 	return nil
// }

// GormDataType tells GORM what the database column type should be.
// For PostgreSQL, it's "uuid". For MySQL, "binary(16)".
// GORM's default dialect handling should manage this correctly.
func (u UUID) GormDataType() string {
	return "uuid"
}

// Value implements the driver.Valuer interface for database writes.
func (u UUID) Value() (driver.Value, error) {
	if u == UUID(uuid.Nil) {
		return nil, nil
	}
	return uuid.UUID(u).String(), nil
}

// Scan implements the sql.Scanner interface for database reads.
func (u *UUID) Scan(src any) error {
	if src == nil {
		*u = UUID(uuid.Nil)
		return nil
	}

	var id uuid.UUID
	switch src := src.(type) {
	case string:
		err := id.Scan(src)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported scan type for UUID: %T", src)
	}
	*u = UUID(id)
	return nil
}
