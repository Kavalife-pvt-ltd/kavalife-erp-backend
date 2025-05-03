package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DB_URL string
}

func ConfigLoader() Config {
	//
	config := Config{}
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access environment variables
	config.Port = os.Getenv("PORT")
	config.DB_URL = os.Getenv("DB_URL")
	return config
}
