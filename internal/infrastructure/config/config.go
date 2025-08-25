package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
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
	Env     string        `mapstructure:"env" validate:"required" default:"local"`
	Host    string        `mapstructure:"host" validate:"required" default:"localhost"`
	Port    string        `mapstructure:"port" validate:"required" default:"8080"`
	Timeout time.Duration `mapstructure:"timeout" default:"1m"`
}

type DBConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            string        `mapstructure:"port" validate:"required"`
	Username        string        `mapstructure:"username" validate:"required"`
	Password        string        `mapstructure:"password"` // Password can be empty
	Name            string        `mapstructure:"name" validate:"required"`
	SSLMode         string        `mapstructure:"sslmode" validate:"required" default:"disable"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" default:"20"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime" default:"1h"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time" default:"30m"`
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

	// Set default values from struct tags
	setDefaultsFromTags(reflect.TypeOf(Config{}), "")

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

// setDefaultsFromTags recursively sets default values in viper from "default" struct tags.
func setDefaultsFromTags(t reflect.Type, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip if not an exported field or if it's the `Other` field (mapstructure:",remain")
		if field.PkgPath != "" || field.Name == "Other" {
			continue
		}

		mapstructureTag := field.Tag.Get("mapstructure")
		var key string
		if mapstructureTag != "" && !strings.Contains(mapstructureTag, ",") {
			key = mapstructureTag
		} else {
			key = strings.ToLower(field.Name)
		}

		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if field.Type.Kind() == reflect.Struct {
			setDefaultsFromTags(field.Type, fullKey)
			continue
		}

		defaultValue := field.Tag.Get("default")
		if defaultValue == "" {
			continue
		}

		// Only set if viper does not already have a value for this key
		if !viper.IsSet(fullKey) {
			switch field.Type.Kind() {
			case reflect.String:
				viper.SetDefault(fullKey, defaultValue)
			case reflect.Int:
				if val, err := strconv.Atoi(defaultValue); err == nil {
					viper.SetDefault(fullKey, val)
				} else {
					log.Printf("Warning: Could not parse default int value '%s' for key '%s': %v", defaultValue, fullKey, err)
				}
			case reflect.Int64:
				if field.Type == reflect.TypeOf(time.Duration(0)) {
					if val, err := time.ParseDuration(defaultValue); err == nil {
						viper.SetDefault(fullKey, val)
					} else {
						log.Printf("Warning: Could not parse default time.Duration value '%s' for key '%s': %v", defaultValue, fullKey, err)
					}
				} else {
					if val, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
						viper.SetDefault(fullKey, val)
					} else {
						log.Printf("Warning: Could not parse default int64 value '%s' for key '%s': %v", defaultValue, fullKey, err)
					}
				}
			case reflect.Bool:
				if val, err := strconv.ParseBool(defaultValue); err == nil {
					viper.SetDefault(fullKey, val)
				} else {
					log.Printf("Warning: Could not parse default bool value '%s' for key '%s': %v", defaultValue, fullKey, err)
				}
			default:
				log.Printf("Warning: Unsupported type for default tag on key '%s': %s", fullKey, field.Type.Kind())
			}
		}
	}
}
