package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// RootHandler handles the root endpoint
type RootHandler struct{}

// NewRootHandler creates a new RootHandler
func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

// RootResponse represents the response for the root endpoint
type RootResponse struct {
	Message   string   `json:"message"`
	Version   string   `json:"version"`
	Endpoints []string `json:"endpoints"`
}

// HealthResponse represents the response for the health check endpoint
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// ServeHTTP handles HTTP requests for the root endpoint
func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := RootResponse{
		Message: "Welcome to BelajarGolang API",
		Version: "1.0.0",
		Endpoints: []string{
			"GET    /              - API information",
			"GET    /health        - Health check",
			"GET    /categories    - Get all categories",
			"POST   /categories    - Create a category",
			"GET    /categories/{id} - Get a category by ID",
			"GET    /categories?name={name} - Search categories by name",
			"PUT    /categories/{id} - Update a category",
			"DELETE /categories/{id} - Delete a category",
			"GET    /products      - Get all products",
			"POST   /products      - Create a product",
			"GET    /products/{id} - Get a product by ID",
			"GET    /products?name={name} - Search products by name",
			"GET    /products?category_id={id} - Get products by category",
			"PUT    /products/{id} - Update a product",
			"DELETE /products/{id} - Delete a product",
			"POST   /transactions/checkout - Create a new transaction (checkout)",
			"GET    /transactions    - Get all transactions",
			"GET    /transactions/{id} - Get a transaction by ID",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthHandler handles the health check endpoint
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP handles HTTP requests for the health check endpoint
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "BelajarGolang API",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
