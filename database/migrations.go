package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

// RunMigrations creates the necessary database tables
func RunMigrations(db *pgx.Conn) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(context.Background(), migration)
		if err != nil {
			return err
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedCategories seeds initial category data if the table is empty
func SeedCategories(db *pgx.Conn) error {
	// Check if data already exists
	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Categories table already has data, skipping seed")
		return nil
	}

	// Seed initial data
	seedData := []struct {
		Name        string
		Description string
	}{
		{"Electronics", "Electronic devices and gadgets"},
		{"Clothing", "Apparel and fashion items"},
		{"Books", "Books and reading materials"},
		{"Food & Beverages", "Food products and drinks"},
		{"Sports", "Sports equipment and accessories"},
	}

	for _, data := range seedData {
		_, err := db.Exec(context.Background(),
			"INSERT INTO categories (name, description) VALUES ($1, $2)",
			data.Name, data.Description)
		if err != nil {
			return err
		}
	}

	log.Println("Database seeding completed successfully")
	return nil
}
