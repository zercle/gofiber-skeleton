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
	Cache    CacheConfig
	JWT      JWTConfig
	CORS     CORSConfig
	RateLimit RateLimitConfig
	Log      LogConfig
	Swagger  SwaggerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string
	Port         int
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Max        int
	Expiration time.Duration
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
	Output string
}

// SwaggerConfig holds Swagger configuration
type SwaggerConfig struct {
	Enabled  bool
	Host     string
	BasePath string
}

// Load reads configuration from environment and files
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Read config file if exists
	_ = viper.ReadInConfig()

	// Bind environment variables
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	config := &Config{
		Server: ServerConfig{
			Host:         viper.GetString("SERVER_HOST"),
			Port:         viper.GetInt("SERVER_PORT"),
			Env:          viper.GetString("SERVER_ENV"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetInt("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Name:            viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSLMODE"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Cache: CacheConfig{
			Host:     viper.GetString("VALKEY_HOST"),
			Port:     viper.GetInt("VALKEY_PORT"),
			Password: viper.GetString("VALKEY_PASSWORD"),
			DB:       viper.GetInt("VALKEY_DB"),
			PoolSize: viper.GetInt("VALKEY_POOL_SIZE"),
		},
		JWT: JWTConfig{
			Secret:           viper.GetString("JWT_SECRET"),
			ExpiresIn:        viper.GetDuration("JWT_EXPIRES_IN"),
			RefreshExpiresIn: viper.GetDuration("JWT_REFRESH_EXPIRES_IN"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods:   viper.GetStringSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders:   viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
			AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		},
		RateLimit: RateLimitConfig{
			Max:        viper.GetInt("RATE_LIMIT_MAX"),
			Expiration: viper.GetDuration("RATE_LIMIT_EXPIRATION"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
			Output: viper.GetString("LOG_OUTPUT"),
		},
		Swagger: SwaggerConfig{
			Enabled:  viper.GetBool("SWAGGER_ENABLED"),
			Host:     viper.GetString("SWAGGER_HOST"),
			BasePath: viper.GetString("SWAGGER_BASE_PATH"),
		},
	}

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 3000)
	viper.SetDefault("SERVER_ENV", "development")
	viper.SetDefault("SERVER_READ_TIMEOUT", 10*time.Second)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "gofiber_skeleton")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Cache defaults
	viper.SetDefault("VALKEY_HOST", "localhost")
	viper.SetDefault("VALKEY_PORT", 6379)
	viper.SetDefault("VALKEY_PASSWORD", "")
	viper.SetDefault("VALKEY_DB", 0)
	viper.SetDefault("VALKEY_POOL_SIZE", 10)

	// JWT defaults
	viper.SetDefault("JWT_SECRET", "change-me-in-production")
	viper.SetDefault("JWT_EXPIRES_IN", 24*time.Hour)
	viper.SetDefault("JWT_REFRESH_EXPIRES_IN", 168*time.Hour)

	// CORS defaults
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	viper.SetDefault("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", true)

	// Rate limit defaults
	viper.SetDefault("RATE_LIMIT_MAX", 100)
	viper.SetDefault("RATE_LIMIT_EXPIRATION", 60*time.Second)

	// Log defaults
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_OUTPUT", "stdout")

	// Swagger defaults
	viper.SetDefault("SWAGGER_ENABLED", true)
	viper.SetDefault("SWAGGER_HOST", "localhost:3000")
	viper.SetDefault("SWAGGER_BASE_PATH", "/api/v1")
}

// validate validates configuration values
func validate(config *Config) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.JWT.Secret == "change-me-in-production" && config.Server.Env == "production" {
		return fmt.Errorf("JWT secret must be changed in production")
	}

	if config.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetCacheAddress returns the cache address
func (c *Config) GetCacheAddress() string {
	return fmt.Sprintf("%s:%d", c.Cache.Host, c.Cache.Port)
}
