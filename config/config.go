package config

import (
    "fmt"

    "github.com/spf13/viper"
)

type Config struct {
    Port           string
    DatabaseURL    string
    JWTPrivateKey  string
    JWTPublicKey   string
}

func LoadConfig() *Config {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    viper.SetDefault("PORT", "8080")
    viper.BindEnv("DATABASE_URL")
    viper.BindEnv("JWT_PRIVATE_KEY")
    viper.BindEnv("JWT_PUBLIC_KEY")

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            fmt.Println("Config file not found, using environment variables")
        } else {
            panic(fmt.Errorf("failed to read config file: %w", err))
        }
    }

    return &Config{
        Port:          viper.GetString("PORT"),
        DatabaseURL:   viper.GetString("DATABASE_URL"),
        JWTPrivateKey: viper.GetString("JWT_PRIVATE_KEY"),
        JWTPublicKey:  viper.GetString("JWT_PUBLIC_KEY"),
    }
}