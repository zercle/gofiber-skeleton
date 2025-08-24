package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App AppConfig `mapstructure:"app" validate:"required"`
	DB  DBConfig  `mapstructure:"db" validate:"required"`
	JWT JWTConfig `mapstructure:"jwt" validate:"required"`
}

type AppConfig struct {
	Env  string `mapstructure:"env" validate:"required"`
	Host string `mapstructure:"host" validate:"required"`
	Port string `mapstructure:"port" validate:"required"`
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

// LoadConfig loads configuration from environment variables
func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env") // Load .env file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(".env file not found, loading from environment variables only.")
		} else {
			return config, fmt.Errorf("failed to read .env file: %w", err)
		}
	}

	viper.AutomaticEnv() // Read from environment variables

	// Set default values if not found in .env or environment variables
	viper.SetDefault("APP_PORT", "8080")

	if err = viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration using go-playground/validator
	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return config, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}
