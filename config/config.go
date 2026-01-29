package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func GetPort() string {
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func GetDatabaseURL() string {
	// First try DATABASE_URL (Railway's default)
	dbURL := viper.GetString("DATABASE_URL")
	if dbURL != "" {
		return dbURL
	}

	// Fallback: Try to build from individual Railway Postgres variables
	// Railway provides: PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	database := os.Getenv("PGDATABASE")

	if host != "" && user != "" && database != "" {
		if port == "" {
			port = "5432"
		}
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password, host, port, database)
	}

	return ""
}
