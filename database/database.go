package database

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

var ErrEmptyConnectionString = errors.New("DATABASE_URL environment variable is empty")

func InitDB(connectionString string) (*pgx.Conn, error) {
	// Check if connection string is provided
	if connectionString == "" {
		log.Println("ERROR: DATABASE_URL is empty or not set")
		return nil, ErrEmptyConnectionString
	}

	log.Printf("Connecting to database...")

	// Parse connection config
	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Printf("ERROR: Failed to parse connection string: %v", err)
		return nil, err
	}

	// Disable prepared statement cache for compatibility with connection poolers
	// (PgBouncer, Supabase, Railway, etc.)
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Open database
	db, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Printf("ERROR: Failed to connect to database: %v", err)
		return nil, err
	}

	// Test connection
	err = db.Ping(context.Background())
	if err != nil {
		log.Printf("ERROR: Failed to ping database: %v", err)
		return nil, err
	}

	// Example query to test connection
	var version string
	if err := db.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Printf("ERROR: Query failed: %v", err)
		return nil, err
	}

	log.Println("Connected to:", version)
	log.Println("Database connected successfully")
	return db, nil
}
