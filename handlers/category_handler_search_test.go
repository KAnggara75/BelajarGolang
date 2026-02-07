package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGetCategoryByName_Success tests GET /categories?name=Elect (partial match)
func TestGetCategoryByName_Success(t *testing.T) {
	handler := setupTestHandlerWithData()

	// Search for "Elect" which should match "Electronics"
	req := httptest.NewRequest(http.MethodGet, "/categories?name=Elect", nil)
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

	if response.Message != "Categories retrieved successfully" {
		t.Errorf("Expected message 'Categories retrieved successfully', got '%s'", response.Message)
	}

	// Check category data - should be an array
	dataArray, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}

	if len(dataArray) == 0 {
		t.Fatal("Expected at least one category in results")
	}

	// Check first category
	firstCat, ok := dataArray[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first item to be an object, got %T", dataArray[0])
	}

	// Should match "Electronics" when searching for "Elect"
	categoryName := firstCat["name"].(string)
	if !strings.Contains(categoryName, "Elect") {
		t.Errorf("Expected name to contain 'Elect', got '%v'", firstCat["name"])
	}
}

// TestGetCategoryByName_NotFound tests GET /categories?name=NonExistent
func TestGetCategoryByName_NotFound(t *testing.T) {
	handler := setupTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/categories?name=NonExistent", nil)
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

	if response.Message != "Category not found" {
		t.Errorf("Expected message 'Category not found', got '%s'", response.Message)
	}
}

// TestGetCategoryByName_URLEncoded tests GET /categories?name=Food%20%26%20Beverages
func TestGetCategoryByName_URLEncoded(t *testing.T) {
	handler := setupTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/categories?name=Food%20%26%20Beverages", nil)
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

	// Check category data - should be an array
	dataArray, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}

	if len(dataArray) == 0 {
		t.Fatal("Expected at least one category in results")
	}

	firstCat, ok := dataArray[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first item to be an object, got %T", dataArray[0])
	}

	if firstCat["name"] != "Food & Beverages" {
		t.Errorf("Expected name 'Food & Beverages', got '%v'", firstCat["name"])
	}
}
