package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Helper provides common test utilities
type Helper struct {
	t *testing.T
}

// New creates a new test helper
func New(t *testing.T) *Helper {
	t.Helper()
	return &Helper{t: t}
}

// Context creates a test context with timeout
func (h *Helper) Context() context.Context {
	h.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	h.t.Cleanup(cancel)
	return ctx
}

// NewUUID generates a new UUID for testing
func (h *Helper) NewUUID() uuid.UUID {
	h.t.Helper()
	id, err := uuid.NewV7()
	require.NoError(h.t, err)
	return id
}

// RequireNoError fails the test if error is not nil
func (h *Helper) RequireNoError(err error) {
	h.t.Helper()
	require.NoError(h.t, err)
}

// RequireError fails the test if error is nil
func (h *Helper) RequireError(err error) {
	h.t.Helper()
	require.Error(h.t, err)
}

// TimeNow returns a consistent time for testing
func TimeNow() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

// TimeParse parses a time string for testing
func TimeParse(t *testing.T, layout, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(layout, value)
	require.NoError(t, err)
	return parsed
}
