package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/paaart/kavalife-erp-backend/internal/config"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
	util "github.com/paaart/kavalife-erp-backend/internal/utils"
)

func main() {

	Log := util.InitLogger()
	appConfig := config.ConfigLoader()

	pool, err := db.Connect(appConfig) // connection to postgres
	if err != nil {
		Log.Fatal("Database connection failed:", err)
	}
	defer pool.Close()

	Log.Info("Starting Kava Life ERP Backend...")
	r := handlers.RunApp()

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
