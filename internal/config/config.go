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
	Redis    RedisConfig
	JWT      JWTConfig
	App      AppConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	AccessExpiration time.Duration
	RefreshExpiration time.Duration
}

// AppConfig holds general application configuration
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
	Version     string
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	v := viper.New()

	// Set config file properties
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	// Read config file (optional, fall back to env vars)
	_ = v.ReadInConfig()

	// Bind environment variables
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	cfg := &Config{
		Server: ServerConfig{
			Port:         v.GetString("PORT"),
			ReadTimeout:  v.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  v.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetString("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			DBName:          v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSLMODE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetString("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:            v.GetString("JWT_SECRET"),
			AccessExpiration:  v.GetDuration("JWT_ACCESS_EXPIRATION"),
			RefreshExpiration: v.GetDuration("JWT_REFRESH_EXPIRATION"),
		},
		App: AppConfig{
			Name:        v.GetString("APP_NAME"),
			Environment: v.GetString("APP_ENV"),
			Debug:       v.GetBool("APP_DEBUG"),
			Version:     v.GetString("APP_VERSION"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("PORT", "3000")
	v.SetDefault("SERVER_READ_TIMEOUT", 10*time.Second)
	v.SetDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)
	v.SetDefault("SERVER_IDLE_TIMEOUT", 120*time.Second)

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "gofiber_skeleton")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Redis defaults
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	// JWT defaults
	v.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")
	v.SetDefault("JWT_ACCESS_EXPIRATION", 15*time.Minute)
	v.SetDefault("JWT_REFRESH_EXPIRATION", 7*24*time.Hour)

	// App defaults
	v.SetDefault("APP_NAME", "gofiber-skeleton")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_DEBUG", true)
	v.SetDefault("APP_VERSION", "1.0.0")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key-change-in-production" {
		if c.App.Environment == "production" {
			return fmt.Errorf("JWT secret must be set in production")
		}
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis connection address
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// IsProduction returns true if the app is running in production
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the app is running in development
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}
