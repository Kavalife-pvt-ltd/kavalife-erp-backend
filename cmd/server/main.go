package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/paaart/kavalife-erp-backend/config"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/routes"
	util "github.com/paaart/kavalife-erp-backend/internal/utils"
)

// @title           KavaLife ERP API
// @version         1.0
// @description     Backend API for KavaLife ERP
// @host            localhost:8080
// @BasePath        /
func main() {

	Log := util.InitLogger()
	appConfig := config.ConfigLoader()
	godotenv.Load()

	pool, err := db.Connect(appConfig) // connection to postgres
	if err != nil {
		Log.Fatal("Database connection failed:", err)
	}
	defer pool.Close()

	Log.Info("Starting Kava Life ERP Backend...")
	r := routes.RunApp()

	// Graceful shutdown on interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		Log.Info("Shutting down server...")

		// Close DB pool
		pool.Close()

		// If you have other cleanup (e.g., HTTP shutdown), do it here

		os.Exit(0)
	}()

	Log.Info("Database connection established successfully")

	Log.Info("Server will run on port:", appConfig.Port)
	if err := r.Run(":" + appConfig.Port); err != nil {
		Log.Fatal("Failed to run server:", err)
	}
}
