package config

import (
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"db"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Host        string        `mapstructure:"host" default:"localhost"`
	Port        uint16        `mapstructure:"port" default:"8080"`
	Env         string        `mapstructure:"env" default:"development"`
	ReadTimeout time.Duration `mapstructure:"read_timeout" default:"1m"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host" validate:"required" default:"localhost"`
	Port            string        `mapstructure:"port" validate:"required" default:"5432"`
	User            string        `mapstructure:"user" default:"postgres"`
	Password        string        `mapstructure:"password" default:"postgres"`
	Name            string        `mapstructure:"name" default:"gofiber_skeleton"`
	SSLMode         string        `mapstructure:"ssl_mode" default:"disable"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"5m"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" default:"2m"`
}

type JWTConfig struct {
	Secret    string        `mapstructure:"secret" default:"your_jwt_secret_here_change_in_production"`
	ExpiresIn time.Duration `mapstructure:"expiration_in" default:"1h"`
}

func Load() *Config {
	setupViper()

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return config
}

func setupViper() {
	// Set defaults first
	setDefaults()

	// Set default environment
	env := getEnv("APP_ENV", "development")

	// Set up Viper to read from config files
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../../configs") // For test environment

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind specific environment variables to config keys
	bindEnvVars()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found, using environment variables and defaults")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
}

func setDefaults() {
	config := &Config{}
	setDefaultsFromTags(reflect.ValueOf(config).Elem(), "")
}

func setDefaultsFromTags(v reflect.Value, prefix string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get mapstructure tag for key name
		mapstructureTag := field.Tag.Get("mapstructure")
		if mapstructureTag == "" {
			mapstructureTag = strings.ToLower(field.Name)
		}

		// Build the full key path
		var key string
		if prefix == "" {
			key = mapstructureTag
		} else {
			key = prefix + "." + mapstructureTag
		}

		// Handle nested structs
		if fieldValue.Kind() == reflect.Struct {
			setDefaultsFromTags(fieldValue, key)
			continue
		}

		// Set default value if default tag exists
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			viper.SetDefault(key, defaultValue)
		}
	}

}

func bindEnvVars() {
	// Use Viper's automatic environment binding with reflection
	config := &Config{}
	bindEnvVarsRecursively(reflect.ValueOf(config).Elem(), "")
}

func bindEnvVarsRecursively(v reflect.Value, prefix string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get mapstructure tag for key name
		mapstructureTag := field.Tag.Get("mapstructure")
		if mapstructureTag == "" {
			mapstructureTag = strings.ToLower(field.Name)
		}

		// Build the full key path
		var key string
		if prefix == "" {
			key = mapstructureTag
		} else {
			key = prefix + "." + mapstructureTag
		}

		// Handle nested structs
		if fieldValue.Kind() == reflect.Struct {
			bindEnvVarsRecursively(fieldValue, key)
			continue
		}

		// Automatically bind environment variable
		// Viper will look for environment variables that match the key with underscores
		if err := viper.BindEnv(key); err != nil {
			log.Printf("Failed to bind env var %s: %v", key, err)
		}
	}
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password +
		"@" + c.Database.Host + ":" + c.Database.Port +
		"/" + c.Database.Name + "?sslmode=" + c.Database.SSLMode
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
