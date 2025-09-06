package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type AppConfig struct {
	Port        string `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
}

type DatabaseConfig struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type JWTConfig struct {
	Secret    string        `mapstructure:"secret"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
	Issuer    string        `mapstructure:"issuer"`
}

type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("app.port", "3000")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.name", "gofiber-skeleton")
	viper.SetDefault("app.version", "1.0.0")

	viper.SetDefault("database.url", "postgres://user:password@localhost:5432/gofiber_skeleton?sslmode=disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", time.Hour)
	viper.SetDefault("database.conn_max_idle_time", 30*time.Minute)

	viper.SetDefault("jwt.secret", "your-secret-key-change-this-in-production")
	viper.SetDefault("jwt.expires_in", 24*time.Hour)
	viper.SetDefault("jwt.issuer", "gofiber-skeleton")

	viper.SetDefault("redis.url", "redis://localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("cors.expose_headers", []string{})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 3600)

	bindEnvVariables()
}

func bindEnvVariables() {
	viper.BindEnv("app.port", "PORT")
	viper.BindEnv("app.environment", "ENV", "ENVIRONMENT")
	viper.BindEnv("database.url", "DATABASE_URL")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expires_in", "JWT_EXPIRES_IN")
	viper.BindEnv("redis.url", "REDIS_URL")
	viper.BindEnv("cors.allow_origins", "CORS_ORIGINS")
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) IsStaging() bool {
	return c.App.Environment == "staging"
}