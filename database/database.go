package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func InitDB(connectionString string) (*pgx.Conn, error) {
	// Parse connection config
	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// Disable prepared statement cache for compatibility with connection poolers
	// (PgBouncer, Supabase, etc.)
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Open database
	db, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	// Example query to test connection
	var version string
	if err := db.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	log.Println("Database connected successfully")
	return db, nil
}
