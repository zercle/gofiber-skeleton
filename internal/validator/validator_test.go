package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/validator"
)

type TestUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Age      int    `json:"age" validate:"required,gte=18,lte=120"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
}

type TestProduct struct {
	Name  string  `json:"name" validate:"required,min=1,max=100"`
	Price float64 `json:"price" validate:"required,gt=0"`
	SKU   string  `json:"sku" validate:"required,len=8,alphanum"`
}

type TestPassword struct {
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func TestValidator_Validate(t *testing.T) {
	v := validator.New()

	t.Run("valid user data", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "password123",
			Age:      25,
			Username: "testuser",
		}

		errors := v.Validate(user)
		assert.Empty(t, errors)
	})

	t.Run("missing required fields", func(t *testing.T) {
		user := TestUser{
			Email: "test@example.com",
			// Missing Password, Age, Username
		}

		errors := v.Validate(user)
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 3)

		// Check that all missing fields are reported
		fields := make(map[string]bool)
		for _, err := range errors {
			fields[err.Field] = true
			assert.Equal(t, "required", err.Tag)
			assert.Contains(t, err.Message, "is required")
		}

		assert.True(t, fields["password"])
		assert.True(t, fields["age"])
		assert.True(t, fields["username"])
	})

	t.Run("invalid email format", func(t *testing.T) {
		user := TestUser{
			Email:    "invalid-email",
			Password: "password123",
			Age:      25,
			Username: "testuser",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "email", errors[0].Field)
		assert.Equal(t, "email", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "valid email")
	})

	t.Run("password too short", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "short",
			Age:      25,
			Username: "testuser",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "password", errors[0].Field)
		assert.Equal(t, "min", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "at least")
		assert.Contains(t, errors[0].Message, "8")
	})

	t.Run("age below minimum", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "password123",
			Age:      15,
			Username: "testuser",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "age", errors[0].Field)
		assert.Equal(t, "gte", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "greater than or equal to")
	})

	t.Run("age above maximum", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "password123",
			Age:      150,
			Username: "testuser",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "age", errors[0].Field)
		assert.Equal(t, "lte", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "less than or equal to")
	})

	t.Run("username too short", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "password123",
			Age:      25,
			Username: "ab",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "username", errors[0].Field)
		assert.Equal(t, "min", errors[0].Tag)
	})

	t.Run("username too long", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "password123",
			Age:      25,
			Username: "thisusernameiswaytoolongforvalidation",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "username", errors[0].Field)
		assert.Equal(t, "max", errors[0].Tag)
	})

	t.Run("username not alphanumeric", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Password: "password123",
			Age:      25,
			Username: "test@user",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 1)
		assert.Equal(t, "username", errors[0].Field)
		assert.Equal(t, "alphanum", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "letters and numbers")
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		user := TestUser{
			Email:    "invalid-email",
			Password: "short",
			Age:      15,
			Username: "ab",
		}

		errors := v.Validate(user)
		assert.Len(t, errors, 4)
	})
}

func TestValidator_ValidateProduct(t *testing.T) {
	v := validator.New()

	t.Run("valid product", func(t *testing.T) {
		product := TestProduct{
			Name:  "Test Product",
			Price: 99.99,
			SKU:   "ABC12345",
		}

		errors := v.Validate(product)
		assert.Empty(t, errors)
	})

	t.Run("price must be greater than zero", func(t *testing.T) {
		product := TestProduct{
			Name:  "Test Product",
			Price: -10.5,
			SKU:   "ABC12345",
		}

		errors := v.Validate(product)
		assert.Len(t, errors, 1)
		assert.Equal(t, "price", errors[0].Field)
		assert.Equal(t, "gt", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "greater than")
	})

	t.Run("SKU exact length", func(t *testing.T) {
		product := TestProduct{
			Name:  "Test Product",
			Price: 99.99,
			SKU:   "ABC123", // Too short
		}

		errors := v.Validate(product)
		assert.Len(t, errors, 1)
		assert.Equal(t, "sku", errors[0].Field)
		assert.Equal(t, "len", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "exactly")
		assert.Contains(t, errors[0].Message, "8")
	})
}

func TestValidator_ValidatePassword(t *testing.T) {
	v := validator.New()

	t.Run("passwords match", func(t *testing.T) {
		passwords := TestPassword{
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		errors := v.Validate(passwords)
		assert.Empty(t, errors)
	})

	t.Run("passwords do not match", func(t *testing.T) {
		passwords := TestPassword{
			Password:        "password123",
			ConfirmPassword: "different123",
		}

		errors := v.Validate(passwords)
		assert.Len(t, errors, 1)
		assert.Equal(t, "confirm_password", errors[0].Field)
		assert.Equal(t, "eqfield", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "must equal")
	})
}

func TestValidator_ValidateVar(t *testing.T) {
	v := validator.New()

	t.Run("valid email variable", func(t *testing.T) {
		email := "test@example.com"
		err := v.ValidateVar(email, "required,email")
		assert.NoError(t, err)
	})

	t.Run("invalid email variable", func(t *testing.T) {
		email := "invalid-email"
		err := v.ValidateVar(email, "required,email")
		assert.Error(t, err)
	})

	t.Run("empty required variable", func(t *testing.T) {
		value := ""
		err := v.ValidateVar(value, "required")
		assert.Error(t, err)
	})

	t.Run("numeric validation", func(t *testing.T) {
		value := "12345"
		err := v.ValidateVar(value, "numeric")
		assert.NoError(t, err)
	})

	t.Run("non-numeric value", func(t *testing.T) {
		value := "abc123"
		err := v.ValidateVar(value, "numeric")
		assert.Error(t, err)
	})
}

func TestValidator_ErrorMessages(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name          string
		value         interface{}
		expectedField string
		expectedTag   string
		messageContains string
	}{
		{
			name: "required field",
			value: struct {
				Name string `json:"name" validate:"required"`
			}{},
			expectedField:   "name",
			expectedTag:     "required",
			messageContains: "is required",
		},
		{
			name: "email format",
			value: struct {
				Email string `json:"email" validate:"email"`
			}{Email: "invalid"},
			expectedField:   "email",
			expectedTag:     "email",
			messageContains: "valid email",
		},
		{
			name: "URL format",
			value: struct {
				Website string `json:"website" validate:"url"`
			}{Website: "not-a-url"},
			expectedField:   "website",
			expectedTag:     "url",
			messageContains: "valid URL",
		},
		{
			name: "UUID format",
			value: struct {
				ID string `json:"id" validate:"uuid"`
			}{ID: "not-a-uuid"},
			expectedField:   "id",
			expectedTag:     "uuid",
			messageContains: "valid UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := v.Validate(tt.value)
			assert.NotEmpty(t, errors)
			assert.Equal(t, tt.expectedField, errors[0].Field)
			assert.Equal(t, tt.expectedTag, errors[0].Tag)
			assert.Contains(t, errors[0].Message, tt.messageContains)
		})
	}
}
