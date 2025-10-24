package errors_test

import (
	"net/http"
	"testing"

	"github.com/zercle/template-go-fiber/internal/errors"
)

func TestNewNotFoundError(t *testing.T) {
	err := errors.NewNotFoundError("user not found")

	if err.Code != errors.ErrCodeNotFound {
		t.Errorf("expected code %s, got %s", errors.ErrCodeNotFound, err.Code)
	}

	if err.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, err.StatusCode)
	}

	if err.Message != "user not found" {
		t.Errorf("expected message 'user not found', got '%s'", err.Message)
	}
}

func TestNewBadRequestError(t *testing.T) {
	err := errors.NewBadRequestError("invalid input")

	if err.Code != errors.ErrCodeBadRequest {
		t.Errorf("expected code %s, got %s", errors.ErrCodeBadRequest, err.Code)
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
}

func TestNewValidationError(t *testing.T) {
	innerErr := errors.NewValidationError("email is required", nil)

	if innerErr.Code != errors.ErrCodeValidation {
		t.Errorf("expected code %s, got %s", errors.ErrCodeValidation, innerErr.Code)
	}

	if innerErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, innerErr.StatusCode)
	}
}

func TestNewDuplicateEntryError(t *testing.T) {
	err := errors.NewDuplicateEntryError("email already exists")

	if err.Code != errors.ErrCodeDuplicateEntry {
		t.Errorf("expected code %s, got %s", errors.ErrCodeDuplicateEntry, err.Code)
	}

	if err.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, err.StatusCode)
	}
}

func TestIsAPIError(t *testing.T) {
	apiErr := errors.NewNotFoundError("not found")

	if !errors.IsAPIError(apiErr) {
		t.Error("expected IsAPIError to return true for APIError")
	}

	notAPIErr := errors.NewInternalError("internal", nil)
	if !errors.IsAPIError(notAPIErr) {
		t.Error("expected IsAPIError to return true for InternalError")
	}
}

func TestAsAPIError(t *testing.T) {
	apiErr := errors.NewNotFoundError("not found")
	converted := errors.AsAPIError(apiErr)

	if converted.Code != errors.ErrCodeNotFound {
		t.Errorf("expected code %s, got %s", errors.ErrCodeNotFound, converted.Code)
	}
}

func TestErrorImplementsError(t *testing.T) {
	err := errors.NewNotFoundError("test")
	errorMsg := err.Error()

	if errorMsg == "" {
		t.Error("expected error string, got empty")
	}

	if errorMsg != "[NOT_FOUND] test" {
		t.Errorf("expected '[NOT_FOUND] test', got '%s'", errorMsg)
	}
}
