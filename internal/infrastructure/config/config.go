package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App   AppConfig      `mapstructure:"app" validate:"required"`
	DB    DBConfig       `mapstructure:"db" validate:"required"`
	JWT   JWTConfig      `mapstructure:"jwt" validate:"required"`
	Other map[string]any `mapstructure:",remain"`
}

type AppConfig struct {
	Env     string        `mapstructure:"env" validate:"required"`
	Host    string        `mapstructure:"host" validate:"required"`
	Port    string        `mapstructure:"port" validate:"required"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password"` // Password can be empty
	Name     string `mapstructure:"name" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" validate:"required"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret" validate:"required"`
}

// LoadConfig loads configuration from YAML files and environment variables
func LoadConfig() (config Config, err error) {
	// Determine environment
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = "local"
	}

	viper.SetDefault("app.timeout", 1*time.Minute)

	// Load environment-specific YAML config
	viper.SetConfigFile(fmt.Sprintf("configs/%s.yaml", env))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found for environment %s, relying on environment variables.", env)
		} else {
			return config, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Load .env file (optional) and merge
	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(".env file not found, skipping.")
		} else {
			return config, fmt.Errorf("failed to read .env file: %w", err)
		}
	}

	// Override with environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal into struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration using go-playground/validator
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return config, fmt.Errorf("configuration validation failed: %w", err)
	}

	log.Printf("config: %+v", config)

	return config, nil
}
