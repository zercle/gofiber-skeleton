package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/gofrs/uuid"
)

// UUIDv7 represents a UUIDv7 type.
type UUIDv7 uuid.UUID

// Value implements the driver.Valuer interface.
func (u UUIDv7) Value() (driver.Value, error) {
	return uuid.UUID(u).String(), nil
}

// Scan implements the sql.Scanner interface.
func (u *UUIDv7) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch s := src.(type) {
	case []byte:
		parsedUUID, err := uuid.FromString(string(s))
		if err != nil {
			return fmt.Errorf("failed to parse UUID from bytes: %w", err)
		}
		*u = UUIDv7(parsedUUID)
		return nil
	case string:
		parsedUUID, err := uuid.FromString(s)
		if err != nil {
			return fmt.Errorf("failed to parse UUID from string: %w", err)
		}
		*u = UUIDv7(parsedUUID)
		return nil
	default:
		return fmt.Errorf("unsupported Scan type for UUIDv7: %T", src)
	}
}

// String returns the string representation of the UUIDv7.
func (u UUIDv7) String() string {
	return uuid.UUID(u).String()
}

// NewUUIDv7 generates a new UUIDv7.
func NewUUIDv7() UUIDv7 {
	return UUIDv7(uuid.Must(uuid.NewV7()))
}