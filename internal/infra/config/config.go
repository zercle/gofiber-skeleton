package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)
type Config struct {
	GoEnv       string `mapstructure:"GO_ENV"`
	AppPort     string `mapstructure:"APP_PORT"`
	GrpcPort    string `mapstructure:"GRPC_PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig() (config *Config, err error) {
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
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config = &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return config, nil
}
