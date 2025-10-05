package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/pkg/validator"
)

type TestStruct struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=20"`
	Age      int    `validate:"required,min=18,max=100"`
	Password string `validate:"required,min=8"`
}

func TestValidateStruct_Valid(t *testing.T) {
	valid := TestStruct{
		Email:    "test@example.com",
		Username: "testuser",
		Age:      25,
		Password: "password123",
	}

	err := validator.ValidateStruct(valid)
	assert.NoError(t, err)
}

func TestValidateStruct_RequiredFieldMissing(t *testing.T) {
	invalid := TestStruct{
		Email:    "",
		Username: "testuser",
		Age:      25,
		Password: "password123",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	invalid := TestStruct{
		Email:    "notanemail",
		Username: "testuser",
		Age:      25,
		Password: "password123",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}

func TestValidateStruct_MinLength(t *testing.T) {
	invalid := TestStruct{
		Email:    "test@example.com",
		Username: "ab", // Too short
		Age:      25,
		Password: "password123",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}

func TestValidateStruct_MaxLength(t *testing.T) {
	invalid := TestStruct{
		Email:    "test@example.com",
		Username: "verylongusernamethatexceedsmaximum",
		Age:      25,
		Password: "password123",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}

func TestValidateStruct_MinValue(t *testing.T) {
	invalid := TestStruct{
		Email:    "test@example.com",
		Username: "testuser",
		Age:      15, // Too young
		Password: "password123",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}

func TestValidateStruct_MaxValue(t *testing.T) {
	invalid := TestStruct{
		Email:    "test@example.com",
		Username: "testuser",
		Age:      150, // Too old
		Password: "password123",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}

func TestValidateStruct_MultipleErrors(t *testing.T) {
	invalid := TestStruct{
		Email:    "invalidemail",
		Username: "ab",
		Age:      15,
		Password: "short",
	}

	err := validator.ValidateStruct(invalid)
	require.Error(t, err)
}
