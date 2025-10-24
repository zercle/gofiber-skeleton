package config_test

import (
	"os"
	"testing"

	"github.com/zercle/template-go-fiber/internal/config"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Clear any env vars for clean test
	os.Clearenv()

	cfg, err := config.Load()

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	// Verify defaults
	if cfg.Server.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Server.Port)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected host 0.0.0.0, got %s", cfg.Server.Host)
	}

	if cfg.Database.Port != 3306 {
		t.Errorf("expected db port 3306, got %d", cfg.Database.Port)
	}

	if cfg.Database.Driver != "mysql" {
		t.Errorf("expected driver mysql, got %s", cfg.Database.Driver)
	}
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("SERVER_PORT", "8080")
	_ = os.Setenv("SERVER_ENV", "development") // Use development to avoid JWT secret validation
	_ = os.Setenv("DB_HOST", "prod-db.example.com")

	defer func() {
		_ = os.Unsetenv("SERVER_PORT")
		_ = os.Unsetenv("SERVER_ENV")
		_ = os.Unsetenv("DB_HOST")
	}()

	cfg, err := config.Load()

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Server.Environment != "development" {
		t.Errorf("expected environment development, got %s", cfg.Server.Environment)
	}

	if cfg.Database.Host != "prod-db.example.com" {
		t.Errorf("expected db host prod-db.example.com, got %s", cfg.Database.Host)
	}
}

func TestConfig_Validate_Success(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        3000,
			Environment: "development",
			Host:        "0.0.0.0",
		},
		Database: config.DatabaseConfig{
			Host: "localhost",
		},
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
		},
	}

	err := cfg.Validate()

	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestConfig_Validate_InvalidPort(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 99999, // Invalid port
		},
		Database: config.DatabaseConfig{
			Host: "localhost",
		},
	}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestConfig_Validate_MissingDBHost(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 3000,
		},
		Database: config.DatabaseConfig{
			Host: "", // Empty host
		},
	}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected error for missing database host")
	}
}

func TestConfig_Validate_DefaultJWTSecretInProduction(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        3000,
			Environment: "production",
			Host:        "0.0.0.0",
		},
		Database: config.DatabaseConfig{
			Host: "localhost",
		},
		JWT: config.JWTConfig{
			Secret: "your-super-secret-jwt-key", // Default secret
		},
	}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected error for default JWT secret in production")
	}
}
