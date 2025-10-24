package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Logging  LoggingConfig
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Host            string
	Port            int
	Environment     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	Driver          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig contains JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	AllowedCredentials bool
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Load config file if exists (but don't fail if it doesn't)
	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Host:            viper.GetString("SERVER_HOST"),
			Port:            viper.GetInt("SERVER_PORT"),
			Environment:     viper.GetString("SERVER_ENV"),
			ReadTimeout:     viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout:    viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			ShutdownTimeout: viper.GetDuration("SERVER_SHUTDOWN_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetInt("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Name:            viper.GetString("DB_NAME"),
			Driver:          viper.GetString("DB_DRIVER"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: viper.GetDuration("JWT_EXPIRATION"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			AllowedCredentials: viper.GetBool("CORS_ALLOWED_CREDENTIALS"),
		},
		Logging: LoggingConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 3000)
	viper.SetDefault("SERVER_ENV", "development")
	viper.SetDefault("SERVER_READ_TIMEOUT", "10s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("SERVER_SHUTDOWN_TIMEOUT", "30s")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "template_go_fiber")
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	viper.SetDefault("JWT_SECRET", "your-super-secret-jwt-key")
	viper.SetDefault("JWT_EXPIRATION", "3600s")

	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("CORS_ALLOWED_CREDENTIALS", false)

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
}

// Validate validates critical configuration
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.JWT.Secret == "your-super-secret-jwt-key" && c.Server.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	return nil
}
