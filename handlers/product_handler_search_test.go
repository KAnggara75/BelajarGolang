package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGetProductByName_Success tests GET /products?name=MacBook (partial match)
func TestGetProductByName_Success(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	// Search for "MacBook" which should match "MacBook Pro M3"
	req := httptest.NewRequest(http.MethodGet, "/products?name=MacBook", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != "Products retrieved successfully" {
		t.Errorf("Expected message 'Products retrieved successfully', got '%s'", response.Message)
	}

	// Check product data - should be an array
	dataArray, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}

	if len(dataArray) == 0 {
		t.Fatal("Expected at least one product in results")
	}

	// Check first product
	firstProduct, ok := dataArray[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first item to be an object, got %T", dataArray[0])
	}

	// Should match "MacBook Pro M3" when searching for "MacBook"
	productName := firstProduct["name"].(string)
	if !strings.Contains(productName, "MacBook") {
		t.Errorf("Expected name to contain 'MacBook', got '%v'", firstProduct["name"])
	}

	// Check category is included
	if firstProduct["category"] == nil {
		t.Error("Expected category to be included")
	}
}

// TestGetProductByName_NotFound tests GET /products?name=NonExistent
func TestGetProductByName_NotFound(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/products?name=NonExistent", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Product not found" {
		t.Errorf("Expected message 'Product not found', got '%s'", response.Message)
	}
}

// TestGetProductByName_CaseInsensitive tests case-insensitive partial search
func TestGetProductByName_CaseInsensitive(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	// Search with different case - "iphone" should match "iPhone 15 Pro"
	req := httptest.NewRequest(http.MethodGet, "/products?name=iphone", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	// Check that we found products containing "iPhone"
	dataArray, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}

	if len(dataArray) == 0 {
		t.Fatal("Expected at least one product in results")
	}

	firstProduct, ok := dataArray[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first item to be an object, got %T", dataArray[0])
	}

	productName := firstProduct["name"].(string)
	if !strings.Contains(strings.ToLower(productName), "iphone") {
		t.Errorf("Expected name to contain 'iphone' (case-insensitive), got '%v'", productName)
	}
}
