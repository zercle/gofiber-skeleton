package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
	CORS   CORSConfig   `mapstructure:"cors"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	Environment    string        `mapstructure:"environment"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// LoggerConfig holds logger-specific configuration
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// CORSConfig holds CORS-specific configuration
type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	AllowHeaders []string `mapstructure:"allow_headers"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Set default values
	setDefaults()

	// Configure viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// .env file not found, which is okay
			fmt.Println("No .env file found, using environment variables only")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Parse environment-specific values
	if err := parseValues(&config); err != nil {
		return nil, fmt.Errorf("error parsing config values: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.shutdown_timeout", "30s")

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")

	// CORS defaults
	viper.SetDefault("cors.allow_origins", []string{"http://localhost:3000"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
}

// parseValues parses and validates configuration values
func parseValues(config *Config) error {
	// Parse shutdown timeout
	if config.Server.ShutdownTimeout == 0 {
		duration, err := time.ParseDuration("30s")
		if err != nil {
			return fmt.Errorf("invalid shutdown timeout: %w", err)
		}
		config.Server.ShutdownTimeout = duration
	}

	// Set port from environment if not set
	if config.Server.Port == 0 {
		if port := os.Getenv("PORT"); port != "" {
			if _, err := fmt.Sscanf(port, "%d", &config.Server.Port); err != nil {
				return fmt.Errorf("invalid PORT value: %s", port)
			}
		}
	}

	return nil
}

// IsProduction returns true if the application is running in production
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production"
}

// IsDevelopment returns true if the application is running in development
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Environment) == "development"
}

// GetAddress returns the server address in format host:port
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}