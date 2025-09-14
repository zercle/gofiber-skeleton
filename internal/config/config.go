package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseDSN string
	JWTSecret   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=user password=password dbname=fiber_forum port=5432 sslmode=disable TimeZone=Asia/Shanghai"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretjwtkey"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}