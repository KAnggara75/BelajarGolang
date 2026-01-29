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
	// Initialize database
	db, err := database.InitDB(config.GetDatabaseURL())
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
		log.Fatal("Failed to seed database:", err)
	}

	// Initialize repository
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	// Setup routes
	http.Handle("/categories", categoryHandler)
	http.Handle("/categories/", categoryHandler)

	// Start server
	port := config.GetPort()
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Println("📦 Available endpoints:")
	fmt.Println("   GET    /categories      - Get all categories")
	fmt.Println("   POST   /categories      - Create a category")
	fmt.Println("   GET    /categories/{id} - Get a category by ID")
	fmt.Println("   PUT    /categories/{id} - Update a category")
	fmt.Println("   DELETE /categories/{id} - Delete a category")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
