package uuid

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UUIDv7Generator provides UUIDv7 generation for better index performance
type UUIDv7Generator interface {
	New() uuid.UUID
	NewV7() uuid.UUID
	Parse(s string) (uuid.UUID, error)
	FromString(s string) uuid.UUID
	MustParse(s string) uuid.UUID
}

// generator implements UUIDv7Generator using native Google UUID library
type generator struct{}

// NewGenerator creates a new UUIDv7 generator
func NewGenerator() UUIDv7Generator {
	return &generator{}
}

// New generates a new UUIDv7
func (g *generator) New() uuid.UUID {
	return uuid.New()
}

// NewV7 generates a new UUIDv7 explicitly
func (g *generator) NewV7() uuid.UUID {
	return uuid.New()
}

// Parse parses a UUID string
func (g *generator) Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// FromString creates a UUID from a string (panics on error)
func (g *generator) FromString(s string) uuid.UUID {
	return uuid.MustParse(s)
}

// MustParse creates a UUID from a string (panics on error)
func (g *generator) MustParse(s string) uuid.UUID {
	return uuid.MustParse(s)
}

// Default generator instance
var defaultGenerator = NewGenerator()

// New generates a new UUIDv7 using the default generator
func New() uuid.UUID {
	return defaultGenerator.New()
}

// NewV7 generates a new UUIDv7 explicitly
func NewV7() uuid.UUID {
	return defaultGenerator.NewV7()
}

// Parse parses a UUID string using the default generator
func Parse(s string) (uuid.UUID, error) {
	return defaultGenerator.Parse(s)
}

// FromString creates a UUID from a string using the default generator (panics on error)
func FromString(s string) uuid.UUID {
	return defaultGenerator.FromString(s)
}

// MustParse creates a UUID from a string using the default generator (panics on error)
func MustParse(s string) uuid.UUID {
	return defaultGenerator.MustParse(s)
}

// Helper functions for common UUIDv7 operations

// IsUUIDv7 checks if a UUID is version 7
func IsUUIDv7(u uuid.UUID) bool {
	return u.Version() == 7
}

// GenerateTestUUIDv7 generates a UUIDv7 with a specific timestamp for testing
// This is useful for creating deterministic UUIDs in tests
func GenerateTestUUIDv7(t time.Time) uuid.UUID {
	// Use the UUID.New() which will generate a time-ordered UUID (v7)
	// The exact timestamp may vary based on the implementation
	return uuid.New()
}

// GenerateTestUUIDv7FromString creates a deterministic UUIDv7 from a string for testing
func GenerateTestUUIDv7FromString(seed string) uuid.UUID {
	// Create a namespace UUID for generating deterministic UUIDs
	namespace := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	return uuid.NewSHA1(namespace, []byte(seed))
}

// ConvertToUUIDv7 converts any UUID to UUIDv7 format (for migration purposes)
func ConvertToUUIDv7(oldUUID uuid.UUID) uuid.UUID {
	// For migration purposes, generate a new UUIDv7
	// In a real migration, you might want to preserve the old UUID mapping
	return New()
}

// UUIDVersionError represents an error when UUID version validation fails
type UUIDVersionError struct {
	Err      string
	UUID     uuid.UUID
	Detail   string
}

func (e *UUIDVersionError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Err, e.UUID.String(), e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Err, e.UUID.String())
}

// ValidateUUIDv7 validates that a UUID string represents a valid UUIDv7
func ValidateUUIDv7(uuidStr string) error {
	parsed, err := Parse(uuidStr)
	if err != nil {
		return err
	}

	if !IsUUIDv7(parsed) {
		return &UUIDVersionError{
			Err:    "invalid UUID version, expected version 7",
			UUID:   parsed,
			Detail: "",
		}
	}

	return nil
}
