package configs

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port uint16 `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     uint16 `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		DBName   string `mapstructure:"DB_NAME"`
		SSLMode  string `mapstructure:"DB_SSLMODE"`
	} `mapstructure:"database"`
	JWT struct {
		Secret     string `mapstructure:"JWT_SECRET"`
		Expiration int    `mapstructure:"JWT_EXPIRATION"`
	} `mapstructure:"jwt"`
	Cache struct {
		Host     string `mapstructure:"host"`
		Port     uint16 `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       uint8  `mapstructure:"db"`
	} `mapstructure:"cache"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(environment string) (config Config, err error) {
	if environment == "" {
		environment = "local"
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(environment)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
