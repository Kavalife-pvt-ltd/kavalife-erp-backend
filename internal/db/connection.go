package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paaart/kavalife-erp-backend/config"
)

var DB *pgxpool.Pool

func Connect(cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DB_URL)
	if err != nil {
		log.Fatalf("Failed to parse DB URL: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// ✅ Use pgx-specific "simple protocol" directive
	var version string
	err = conn.QueryRow(context.Background(), `simple protocol:SELECT version()`).Scan(&version)
	if err != nil {
		log.Fatalf("Version check failed: %v", err)
	}

	log.Println("✅ Connected to Postgres:", version)
	DB = conn
	return conn, nil
}
