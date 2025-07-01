package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)
type Config struct {
	GO_ENV       string `mapstructure:"GO_ENV"`
	APP_PORT     string `mapstructure:"APP_PORT"`
	GRPC_PORT    string `mapstructure:"GRPC_PORT"`
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	JWT_SECRET_KEY string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig() (config Config, err error) {
	goEnv := os.Getenv("GO_ENV")
	if len(goEnv) == 0 {
		goEnv = "local"
	}

	v := viper.New()
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.SetConfigName(goEnv)
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	if err := v.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return
}
