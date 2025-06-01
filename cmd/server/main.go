package main

import (
	"github.com/paaart/kavalife-erp-backend/internal"
	"github.com/paaart/kavalife-erp-backend/internal/config"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	util "github.com/paaart/kavalife-erp-backend/internal/utils"
)

func main() {

	Log := util.InitLogger()
	appConfig := config.ConfigLoader()
	
	db.Connect(appConfig) // connection to postgres

	Log.Info("Starting Kava Life ERP Backend...")
	r := internal.RunApp()

	Log.Info("Database connection established successfully")

	Log.Info("Server will run on port:", appConfig.Port)
	if err := r.Run(":" + appConfig.Port); err != nil {
		Log.Fatal("Failed to run server:", err)
	}
}
