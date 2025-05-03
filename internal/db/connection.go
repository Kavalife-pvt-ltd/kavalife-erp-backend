package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/paaart/kavalife-erp-backend/internal/config"
)

func Connect(config config.Config) {
	conn, err := pgx.Connect(context.Background(), config.DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Postgres Connected to:", version)
}
