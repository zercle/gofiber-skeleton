package configs

import (
	"strings"

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
		Host     string `mapstructure:"host"`
		Port     uint16 `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
	JWT struct {
		Secret     string `mapstructure:"secret"`
		Expiration int    `mapstructure:"expiration"`
	} `mapstructure:"jwt"`
	Cache struct {
		Host     string `mapstructure:"host"`
		Port     uint16 `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       uint8  `mapstructure:"db"`
	} `mapstructure:"cache"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(environment string, configPath ...string) (config Config, err error) {
	viper.Reset()

	if len(configPath) > 0 && configPath[0] != "" {
		viper.AddConfigPath(configPath[0])
		if environment == "" {
			environment = "local"
		}
		viper.SetConfigName(environment)
	} else {
		if environment == "" {
			environment = "local"
		}
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
		viper.SetConfigName(environment)
	}

	viper.SetConfigType("yaml")

	// Bind environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind config file settings
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}