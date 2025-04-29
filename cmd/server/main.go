package main

import (
	"log"

	"github.com/paaart/kavalife-erp-backend/internal/handlers"
)

func main() {
	r := handlers.SetupRouter()

	log.Println("Starting server on http://localhost:8080")
	r.Run(":8080") // Listen and serve on port 8080
}
