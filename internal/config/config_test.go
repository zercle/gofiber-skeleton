package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/config"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Clear any existing env vars
	os.Clearenv()

	cfg, err := config.LoadConfig()

	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Check default values
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Env)
	assert.NotEmpty(t, cfg.Database.DSN)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)
	assert.Equal(t, "redis:6379", cfg.Redis.Addr)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.NotEmpty(t, cfg.JWT.Secret)
	assert.Equal(t, 72, cfg.JWT.ExpirationHours)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Clearenv()
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("SERVER_ENV", "production")
	os.Setenv("DATABASE_MAXOPENCONNS", "50")
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_DB", "1")
	os.Setenv("JWT_SECRET", "my-secret-key")
	os.Setenv("JWT_EXPIRATIONHOURS", "24")

	cfg, err := config.LoadConfig()

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, "production", cfg.Server.Env)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, 1, cfg.Redis.DB)
	assert.Equal(t, "my-secret-key", cfg.JWT.Secret)
	assert.Equal(t, 24, cfg.JWT.ExpirationHours)

	// Cleanup
	os.Clearenv()
}

func TestConfigGetters(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_STRING", "value")
	os.Setenv("TEST_INT", "42")
	os.Setenv("TEST_BOOL", "true")

	_, err := config.LoadConfig()
	require.NoError(t, err)

	// Test getter functions
	assert.Equal(t, "value", config.GetString("TEST_STRING"))
	assert.Equal(t, 42, config.GetInt("TEST_INT"))
	assert.Equal(t, true, config.GetBool("TEST_BOOL"))

	// Cleanup
	os.Clearenv()
}
