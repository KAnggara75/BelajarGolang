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
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			price DECIMAL(10, 2) NOT NULL DEFAULT 0,
			stock INTEGER NOT NULL DEFAULT 0,
			category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		// Add category_id column if it doesn't exist (for existing databases)
		`DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'products' AND column_name = 'category_id'
			) THEN
				ALTER TABLE products ADD COLUMN category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL;
			END IF;
		END $$`,
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

	log.Println("Categories seeding completed successfully")
	return nil
}

// SeedProducts seeds initial product data if the table is empty
func SeedProducts(db *pgx.Conn) error {
	// Check if data already exists
	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Products table already has data, skipping seed")
		return nil
	}

	// Seed initial data with category_id (all Electronics = category_id 1)
	seedData := []struct {
		Name       string
		Price      float64
		Stock      int
		CategoryID int
	}{
		{"iPhone 15 Pro", 999.99, 50, 1},
		{"MacBook Pro M3", 2499.99, 25, 1},
		{"AirPods Pro", 249.99, 100, 1},
		{"iPad Air", 599.99, 40, 1},
		{"Apple Watch Series 9", 399.99, 60, 1},
	}

	for _, data := range seedData {
		_, err := db.Exec(context.Background(),
			"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4)",
			data.Name, data.Price, data.Stock, data.CategoryID)
		if err != nil {
			return err
		}
	}

	log.Println("Products seeding completed successfully")
	return nil
}
