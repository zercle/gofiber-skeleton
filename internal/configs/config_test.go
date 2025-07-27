package configs_test

import (
	"os"
	"path/filepath"
	"testing"

	"gofiber-skeleton/internal/configs" // Updated import

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test case 1: Load config with a specified environment
	t.Run("Load with specified environment", func(t *testing.T) {
		tempDir := t.TempDir()
		localConfigFile := filepath.Join(tempDir, "any_env_name.yaml")
		localContent := `server:
  host: "127.0.0.1"
  port: 8080
db:
  user: "db_user"
jwt:
  secret: "jwt_secret"
cache:
  db: 0
`
		err := os.WriteFile(localConfigFile, []byte(localContent), 0644)
		assert.NoError(t, err)

		cfg, err := configs.LoadConfig("any_env_name", tempDir)
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1", cfg.Server.Host)
		assert.Equal(t, uint16(8080), cfg.Server.Port)
		assert.Equal(t, "db_user", cfg.Database.User)
		assert.Equal(t, "jwt_secret", cfg.JWT.Secret)
		assert.Equal(t, uint8(0), cfg.Cache.DB)
	})

	// Test case 2: Load config with empty environment (should still load from env vars if present)
	t.Run("Load with empty environment and specific env vars", func(t *testing.T) {
		tempDir := t.TempDir()
		localConfigFile := filepath.Join(tempDir, "local.yaml")
		localContent := `server:
  host: "localhost_env"
  port: 9001
`
		err := os.WriteFile(localConfigFile, []byte(localContent), 0644)
		assert.NoError(t, err)

		cfg, err := configs.LoadConfig("", tempDir)
		assert.NoError(t, err)
		assert.Equal(t, "localhost_env", cfg.Server.Host)
		assert.Equal(t, uint16(9001), cfg.Server.Port)
	})

	// Test case 3: Config file not found (when no env vars are set)
	t.Run("Config file not found and no env vars", func(t *testing.T) {
		_, err := configs.LoadConfig("nonexistent_config", "/tmp/nonexistent_path") // No specific file, and no env vars
		assert.NoError(t, err)                                                      // Should not return an error, but an empty config
	})

	// Test case 4: Load with default "local" environment file (no env vars, no configPath)
	t.Run("Load with default local environment file", func(t *testing.T) {
		tempDir := t.TempDir()
		localConfigFile := filepath.Join(tempDir, "local.yaml")
		localContent := `server:
  host: "local_file_host"
  port: 9002
`
		err := os.WriteFile(localConfigFile, []byte(localContent), 0644)
		assert.NoError(t, err)

		cfg, err := configs.LoadConfig("local", tempDir)
		assert.NoError(t, err)
		assert.Equal(t, "local_file_host", cfg.Server.Host)
		assert.Equal(t, uint16(9002), cfg.Server.Port)
	})
}
