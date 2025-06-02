package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paaart/kavalife-erp-backend/internal/config"
)

var DB *pgxpool.Pool

func Connect(config config.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Postgres Connected to:", version)
	DB = conn
	return conn, err
}
