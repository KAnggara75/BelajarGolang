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

// GetOTLPEndpoint returns the OTLP endpoint for traces
func GetOTLPEndpoint() string {
	endpoint := viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318" // Default OTLP HTTP endpoint
	}
	return endpoint
}

// GetServiceName returns the service name for tracing
func GetServiceName() string {
	name := viper.GetString("OTEL_SERVICE_NAME")
	if name == "" {
		name = "belajar-golang-api"
	}
	return name
}

// GetServiceVersion returns the service version
func GetServiceVersion() string {
	version := viper.GetString("SERVICE_VERSION")
	if version == "" {
		version = "1.0.0"
	}
	return version
}

// GetEnvironment returns the deployment environment
func GetEnvironment() string {
	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return env
}

// IsOTelEnabled returns whether OpenTelemetry is enabled
func IsOTelEnabled() bool {
	return viper.GetBool("OTEL_ENABLED")
}
