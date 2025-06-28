package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	GO_ENV       string `mapstructure:"GO_ENV"`
	APP_PORT     string `mapstructure:"APP_PORT"`
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	JWT_SECRET   string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (config Config, err error) {
	goEnv :=os.Getenv("GO_ENV")
	if len(goEnv) == 0 {
		goEnv = "local"
	}

	v := viper.New()
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.SetConfigName("local")
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return
	}

	err = v.Unmarshal(&config)
	return
}
