package main

import (
	"log"

	"github.com/paaart/kavalife-erp-backend/internal/config"
	"github.com/paaart/kavalife-erp-backend/internal/handlers"
)

func main() {
	config := config.ConfigLoader()
	r := handlers.SetupRouter()
	log.Println("Server will run on port:", config.Port)
	r.Run(":" + config.Port) // Listen and serve on port 8080
}
