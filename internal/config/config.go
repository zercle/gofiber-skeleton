package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, will use environment variables
	}

	// Set defaults
	setDefaults()

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER.PORT", "8080")
	viper.SetDefault("SERVER.ENV", "development")

	// Database defaults
	viper.SetDefault("DATABASE.DSN", "host=db user=user password=password dbname=fiber_forum port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	viper.SetDefault("DATABASE.MAXOPENCONNS", 25)
	viper.SetDefault("DATABASE.MAXIDLECONNS", 5)
	viper.SetDefault("DATABASE.CONNMAXLIFETIME", 300)

	// Redis defaults
	viper.SetDefault("REDIS.ADDR", "redis:6379")
	viper.SetDefault("REDIS.PASSWORD", "")
	viper.SetDefault("REDIS.DB", 0)

	// JWT defaults
	viper.SetDefault("JWT.SECRET", "supersecretjwtkey")
	viper.SetDefault("JWT.EXPIRATIONHOURS", 72)
}

// GetString gets a string config value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt gets an integer config value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool gets a boolean config value
func GetBool(key string) bool {
	return viper.GetBool(key)
}
