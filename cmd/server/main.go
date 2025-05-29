package main

import (
	"github.com/paaart/kavalife-erp-backend/internal/config"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
)

func main() {
	//Initialize logger
	config.InitLogger()
	config.Log.Info("Starting Kava Life ERP Backend...")

	// Load configuration and connect to the database
	appConfig := config.ConfigLoader()
	config.Log.Info("Configuration loaded successfully", appConfig.DB_URL)

	// Connect to the database
	db.Connect(appConfig) // connection to postgres
	config.Log.Info("Database connection established successfully")

	// Set up the router and start the server
	r := handlers.SetupRouter()

	config.Log.Info("Server will run on port:", appConfig.Port)
	if err := r.Run(":" + appConfig.Port); err != nil {
		config.Log.Fatal("Failed to run server:", err)
	}
}
