package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KAnggara75/BelajarGolang/config"
	"github.com/KAnggara75/BelajarGolang/database"
	"github.com/KAnggara75/BelajarGolang/handlers"
	"github.com/KAnggara75/BelajarGolang/repository"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}
}

func main() {
	// Get database URL
	dbURL := config.GetDatabaseURL()
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set. Please set DATABASE_URL environment variable or add it to .env file")
	}

	// Initialize database
	db, err := database.InitDB(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(context.Background())

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed initial data
	if err := database.SeedCategories(db); err != nil {
		log.Fatal("Failed to seed categories:", err)
	}
	if err := database.SeedProducts(db); err != nil {
		log.Fatal("Failed to seed products:", err)
	}

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	productHandler := handlers.NewProductHandler(productRepo)

	// Setup routes
	http.Handle("/categories", categoryHandler)
	http.Handle("/categories/", categoryHandler)
	http.Handle("/products", productHandler)
	http.Handle("/products/", productHandler)

	// Start server
	port := config.GetPort()
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Println("📦 Available endpoints:")
	fmt.Println("   GET    /categories      - Get all categories")
	fmt.Println("   POST   /categories      - Create a category")
	fmt.Println("   GET    /categories/{id} - Get a category by ID")
	fmt.Println("   PUT    /categories/{id} - Update a category")
	fmt.Println("   DELETE /categories/{id} - Delete a category")
	fmt.Println("")
	fmt.Println("   GET    /products        - Get all products")
	fmt.Println("   POST   /products        - Create a product")
	fmt.Println("   GET    /products/{id}   - Get a product by ID")
	fmt.Println("   PUT    /products/{id}   - Update a product")
	fmt.Println("   DELETE /products/{id}   - Delete a product")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
