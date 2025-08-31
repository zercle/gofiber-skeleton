package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadDefaults(t *testing.T) {
	viper.Reset()
	config := Load()

	// Test App defaults
	assert.Equal(t, uint16(8080), config.App.Port)
	assert.Equal(t, "development", config.App.Env)

	// Test Database defaults
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "postgres", config.Database.Password)
	assert.Equal(t, "gofiber_skeleton", config.Database.Name)
	assert.Equal(t, "disable", config.Database.SSLMode)

	// Test JWT defaults
	assert.Equal(t, "your_jwt_secret_here_change_in_production", config.JWT.Secret)
	assert.Equal(t, 24*time.Hour, config.JWT.ExpiresIn)
}

func TestConfig_EnvironmentOverride(t *testing.T) {
	viper.Reset()
	// Set environment variables
	t.Setenv("APP_PORT", "3000")
	t.Setenv("APP_ENV", "production")
	t.Setenv("DB_HOST", "remote-db")
	t.Setenv("JWT_SECRET", "super-secret")

	config := Load()

	// Test that environment variables override defaults
	assert.Equal(t, uint16(3000), config.App.Port)
	assert.Equal(t, "production", config.App.Env)
	assert.Equal(t, "remote-db", config.Database.Host)
	assert.Equal(t, "super-secret", config.JWT.Secret)

	// Test that unset env vars still use defaults
	assert.Equal(t, "5432", config.Database.Port) // Should still be default
}

func TestConfig_DatabaseURL(t *testing.T) {
	viper.Reset()
	config := Load()

	expectedURL := "postgres://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable"
	assert.Equal(t, expectedURL, config.DatabaseURL())
}

func TestConfig_IsProduction(t *testing.T) {
	viper.Reset()
	config := Load()

	// Default should be development
	assert.False(t, config.IsProduction())

	// Test production environment
	t.Setenv("APP_ENV", "production")
	config = Load()
	assert.True(t, config.IsProduction())
}

func TestConfig_ViperEnvironmentBinding(t *testing.T) {
	viper.Reset()
	// Test that all environment variables are properly bound through Viper
	envVars := map[string]string{
		"APP_PORT":       "9000",
		"APP_ENV":        "testing",
		"DB_HOST":        "testdb.example.com",
		"DB_PORT":        "5433",
		"DB_USER":        "testuser",
		"DB_PASSWORD":    "testpass",
		"DB_NAME":        "testdb",
		"DB_SSL_MODE":    "require",
		"JWT_SECRET":     "test-secret",
		"JWT_EXPIRES_IN": "48h",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}

	config := Load()

	// Verify all environment variables are loaded correctly via Viper
	assert.Equal(t, uint16(9000), config.App.Port)
	assert.Equal(t, "testing", config.App.Env)
	assert.Equal(t, "testdb.example.com", config.Database.Host)
	assert.Equal(t, "5433", config.Database.Port)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "testdb", config.Database.Name)
	assert.Equal(t, "require", config.Database.SSLMode)
	assert.Equal(t, "test-secret", config.JWT.Secret)
	assert.Equal(t, 48*time.Hour, config.JWT.ExpiresIn)
}