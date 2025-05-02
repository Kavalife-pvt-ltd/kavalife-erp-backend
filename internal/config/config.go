package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string 
	DB_USER string
	DB_PASS string
	DB_NAME string
	DB_HOST string
	DB_PORT string
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
	config.DB_USER = os.Getenv("DB_USER")
	config.DB_HOST = os.Getenv("DB_HOST")
	config.DB_NAME = os.Getenv("DB_NAME")
	config.DB_PORT = os.Getenv("DB_PORT")
	config.DB_PASS = os.Getenv("DB_PASS")
	return config
}
