package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Valkey   ValkeyConfig   `mapstructure:"valkey"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type AppConfig struct {
	Name        string `mapstructure:"name" default:"gofiber-skeleton"`
	Port        string `mapstructure:"port" default:"3000"`
	Environment string `mapstructure:"environment" default:"development"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	Schema          string `mapstructure:"schema"`
	SSLMode         string `mapstructure:"sslmode"`
	URL             string `mapstructure:"url"` // fallback
	MaxOpenConns    int    `mapstructure:"max_open_conns" default:"25"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" default:"25"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" default:"5"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn string `mapstructure:"expires_in" default:"24h"` // e.g., "24h", "720h"
	Issuer    string `mapstructure:"issuer" default:"gofiber-skeleton"`
}

type ValkeyConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db" default:"0"`
	URL      string `mapstructure:"url"` // fallback
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins" default:"[\"*\"]"`
	AllowMethods     []string `mapstructure:"allow_methods" default:"[\"GET\",\"POST\",\"PUT\",\"DELETE\",\"OPTIONS\"]"`
	AllowHeaders     []string `mapstructure:"allow_headers" default:"[\"Origin\",\"Content-Type\",\"Accept\",\"Authorization\"]"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials" default:"false"`
	MaxAge           int      `mapstructure:"max_age" default:"0"`
}

// DatabaseURL constructs a PostgreSQL connection URL from DB_* variables if available,
// otherwise falls back to DB_URL. DB_* variables take precedence.
func (c *Config) DatabaseURL() string {
	// If individual DB_* components are provided, construct URL
	if c.Database.Host != "" && c.Database.Port != "" && c.Database.User != "" && c.Database.Name != "" {
		// Include search_path parameter for schema if provided
		searchPath := ""
		if c.Database.Schema != "" {
			searchPath = fmt.Sprintf("&search_path=%s", c.Database.Schema)
		}
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s%s",
			c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port,
			c.Database.Name, c.Database.SSLMode, searchPath)
	}
	// Fall back to DB_URL
	return c.Database.URL
}

// ValkeyURL constructs a Redis-compatible connection URL from VALKEY_* variables if available,
// otherwise falls back to VALKEY_URL. VALKEY_* variables take precedence.
func (c *Config) ValkeyURL() string {
	// If individual VALKEY_* components are provided, construct URL
	if c.Valkey.Host != "" && c.Valkey.Port != "" {
		auth := ""
		if c.Valkey.Password != "" {
			auth = fmt.Sprintf(":%s@", c.Valkey.Password)
		}
		return fmt.Sprintf("redis://%s%s:%s/%d", auth, c.Valkey.Host, c.Valkey.Port, c.Valkey.DB)
	}
	// Fall back to VALKEY_URL
	return c.Valkey.URL
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func NewConfig() *Config {
	v := viper.New()

	// Set defaults
	v.SetDefault("app.name", "gofiber-skeleton")
	v.SetDefault("app.port", "3000")
	v.SetDefault("app.environment", "development")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.schema", "public")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 25)
	v.SetDefault("database.conn_max_lifetime", 5)
	v.SetDefault("jwt.expires_in", "24h") // Default to 24 hours
	v.SetDefault("jwt.issuer", "gofiber-skeleton")
	v.SetDefault("valkey.port", "6379")
	v.SetDefault("valkey.db", 0)
	v.SetDefault("cors.allow_origins", []string{"*"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	v.SetDefault("cors.allow_credentials", false)
	v.SetDefault("cors.max_age", 0)

	// Config file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables - canonical first, fallbacks second
	v.BindEnv("app.port", "PORT")
	v.BindEnv("app.environment", "ENV")

	// Database: canonical DB_* variables take precedence
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.name", "DB_NAME")
	v.BindEnv("database.schema", "DB_SCHEMA")
	v.BindEnv("database.sslmode", "DB_SSLMODE")
	v.BindEnv("database.url", "DB_URL") // fallback

	// Valkey: canonical VALKEY_* variables take precedence
	v.BindEnv("valkey.host", "VALKEY_HOST")
	v.BindEnv("valkey.port", "VALKEY_PORT")
	v.BindEnv("valkey.password", "VALKEY_PASSWORD")
	v.BindEnv("valkey.db", "VALKEY_DB")
	v.BindEnv("valkey.url", "VALKEY_URL") // fallback

	v.BindEnv("jwt.secret", "JWT_SECRET")
	v.BindEnv("jwt.expires_in", "JWT_EXPIRES_IN")
	v.BindEnv("cors.allow_origins", "CORS_ORIGINS")

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found, using defaults and environment variables")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	// Validate required fields - check both canonical and fallback patterns
	dbURL := config.DatabaseURL()
	if dbURL == "" {
		log.Fatal("Database configuration is required. Provide either DB_* variables or DB_URL")
	}
	if config.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return &config
}
