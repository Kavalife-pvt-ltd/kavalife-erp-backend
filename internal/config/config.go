package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DB_URL     string
	JWT_SECRET string
}

func ConfigLoader() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system env")
	}
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		DB_URL:     getEnv("DB_URL", ""),
		JWT_SECRET: getEnv("JWT_SECRET", ""),
	}

	if cfg.JWT_SECRET == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	return *cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
