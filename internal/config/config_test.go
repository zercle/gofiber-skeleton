package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "testdb")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "testdb", cfg.Database.Host)
	assert.Equal(t, "testdb", cfg.Database.DBName)
	assert.Equal(t, "test-secret", cfg.JWT.Secret)
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{Port: "3000"},
				Database: DatabaseConfig{
					Host:   "localhost",
					DBName: "test",
				},
				JWT: JWTConfig{Secret: "test-secret"},
				App: AppConfig{Environment: "development"},
			},
			wantErr: false,
		},
		{
			name: "missing port",
			config: &Config{
				Server: ServerConfig{Port: ""},
				Database: DatabaseConfig{
					Host:   "localhost",
					DBName: "test",
				},
				JWT: JWTConfig{Secret: "test"},
			},
			wantErr: true,
			errMsg:  "server port is required",
		},
		{
			name: "missing database host",
			config: &Config{
				Server:   ServerConfig{Port: "3000"},
				Database: DatabaseConfig{DBName: "test"},
				JWT:      JWTConfig{Secret: "test"},
			},
			wantErr: true,
			errMsg:  "database host is required",
		},
		{
			name: "missing database name",
			config: &Config{
				Server:   ServerConfig{Port: "3000"},
				Database: DatabaseConfig{Host: "localhost"},
				JWT:      JWTConfig{Secret: "test"},
			},
			wantErr: true,
			errMsg:  "database name is required",
		},
		{
			name: "production without proper JWT secret",
			config: &Config{
				Server: ServerConfig{Port: "3000"},
				Database: DatabaseConfig{
					Host:   "localhost",
					DBName: "test",
				},
				JWT: JWTConfig{Secret: "your-secret-key-change-in-production"},
				App: AppConfig{Environment: "production"},
			},
			wantErr: true,
			errMsg:  "JWT secret must be set in production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, cfg.GetDSN())
}

func TestRedisConfig_GetAddr(t *testing.T) {
	cfg := RedisConfig{
		Host: "localhost",
		Port: "6379",
	}

	assert.Equal(t, "localhost:6379", cfg.GetAddr())
}

func TestAppConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{"production", "production", true},
		{"development", "development", false},
		{"staging", "staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{Environment: tt.env}
			assert.Equal(t, tt.want, cfg.IsProduction())
		})
	}
}

func TestAppConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{"development", "development", true},
		{"production", "production", false},
		{"staging", "staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{Environment: tt.env}
			assert.Equal(t, tt.want, cfg.IsDevelopment())
		})
	}
}

func TestDefaultValues(t *testing.T) {
	// Clear any existing env vars
	os.Clearenv()

	cfg, err := Load()
	require.NoError(t, err)

	// Test server defaults
	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, 10*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, cfg.Server.IdleTimeout)

	// Test database defaults
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "gofiber_skeleton", cfg.Database.DBName)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)

	// Test app defaults
	assert.Equal(t, "gofiber-skeleton", cfg.App.Name)
	assert.Equal(t, "development", cfg.App.Environment)
	assert.True(t, cfg.App.Debug)
}
