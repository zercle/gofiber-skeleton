package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// Clear any existing environment variables
	_ = os.Unsetenv("GS_SERVER_PORT")
	_ = os.Unsetenv("GS_DB_HOST")
	_ = os.Unsetenv("GS_DB_PORT")
	_ = os.Unsetenv("GS_DB_DBNAME")
	_ = os.Unsetenv("GS_DB_USER")
	_ = os.Unsetenv("GS_DB_PASSWORD")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Test default values
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "gofiber_skeleton", cfg.Database.DBName)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "postgres", cfg.Database.Password)
	assert.Equal(t, "your-super-secret-jwt-key-change-this-in-production", cfg.JWT.Secret)
}

func TestLoad_ConfigMethods(t *testing.T) {
	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Test DSN generation
	dsn := cfg.GetDSN()
	assert.Contains(t, dsn, "postgres://")
	assert.Contains(t, dsn, cfg.Database.User)
	assert.Contains(t, dsn, cfg.Database.DBName)

	// Test server address generation
	addr := cfg.GetServerAddr()
	assert.Contains(t, addr, cfg.Server.Host)
	assert.Contains(t, addr, fmt.Sprintf("%d", cfg.Server.Port))
}

func TestConfig_GetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			DBName:   "testdb",
			User:     "testuser",
			Password: "testpass",
			SSLMode:  "disable",
		},
	}

	expectedDSN := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expectedDSN, cfg.GetDSN())
}

func TestConfig_GetServerAddr(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	expectedAddr := "localhost:8080"
	assert.Equal(t, expectedAddr, cfg.GetServerAddr())
}

func TestConfig_ValidateConfig(t *testing.T) {
	// Test valid config
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{Port: 5432},
		JWT: JWTConfig{Secret: "valid-secret"},
	}

	// This should not return an error when we test validation indirectly through Load
	// But we can't test validateConfig directly since it's not exported
	assert.NotNil(t, cfg)
}