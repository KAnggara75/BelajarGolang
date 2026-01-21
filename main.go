package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KAnggara75/BelajarGolang/handlers"
	"github.com/KAnggara75/BelajarGolang/store"
)

func main() {
	// Initialize the in-memory store
	categoryStore := store.NewCategoryStore()
	categoryStore.SeedData()

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryStore)

	// Setup routes
	http.Handle("/categories", categoryHandler)
	http.Handle("/categories/", categoryHandler)

	// Start server
	port := ":8080"
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
