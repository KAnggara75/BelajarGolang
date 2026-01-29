package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/KAnggara75/BelajarGolang/repository"
)

// mockCategoryRepository is a mock implementation of CategoryRepository for testing
type mockCategoryRepository struct {
	categories map[int]models.Category
	nextID     int
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[int]models.Category),
		nextID:     1,
	}
}

func (m *mockCategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	result := make([]models.Category, 0, len(m.categories))
	for _, cat := range m.categories {
		result = append(result, cat)
	}
	return result, nil
}

func (m *mockCategoryRepository) GetByID(ctx context.Context, id int) (models.Category, error) {
	cat, exists := m.categories[id]
	if !exists {
		return models.Category{}, repository.ErrNotFound
	}
	return cat, nil
}

func (m *mockCategoryRepository) Create(ctx context.Context, cat models.Category) (models.Category, error) {
	// Check if name already exists
	for _, existing := range m.categories {
		if existing.Name == cat.Name {
			return models.Category{}, repository.ErrNameExists
		}
	}

	cat.ID = m.nextID
	m.nextID++
	m.categories[cat.ID] = cat
	return cat, nil
}

func (m *mockCategoryRepository) Update(ctx context.Context, id int, cat models.Category) (models.Category, error) {
	if _, exists := m.categories[id]; !exists {
		return models.Category{}, repository.ErrNotFound
	}

	cat.ID = id
	m.categories[id] = cat
	return cat, nil
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id int) error {
	if _, exists := m.categories[id]; !exists {
		return repository.ErrNotFound
	}

	delete(m.categories, id)
	return nil
}

// SeedData adds sample data for testing
func (m *mockCategoryRepository) SeedData() {
	initialData := []models.Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Clothing", Description: "Apparel and fashion items"},
		{Name: "Books", Description: "Books and reading materials"},
		{Name: "Food & Beverages", Description: "Food products and drinks"},
		{Name: "Sports", Description: "Sports equipment and accessories"},
	}

	for _, cat := range initialData {
		_, _ = m.Create(context.Background(), cat)
	}
}

// setupTestHandler creates a fresh handler with an empty mock repository for testing
func setupTestHandler() *CategoryHandler {
	repo := newMockCategoryRepository()
	return NewCategoryHandler(repo)
}

// setupTestHandlerWithData creates a handler with seeded data
func setupTestHandlerWithData() *CategoryHandler {
	repo := newMockCategoryRepository()
	repo.SeedData()
	return NewCategoryHandler(repo)
}

// TestGetAllCategories_Empty tests GET /categories with empty repo
func TestGetAllCategories_Empty(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
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

	// Data should be an empty array
	data, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}
	if len(data) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(data))
	}
}

// TestGetAllCategories_WithData tests GET /categories with seeded data
func TestGetAllCategories_WithData(t *testing.T) {
	handler := setupTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
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

	data, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}
	if len(data) != 5 {
		t.Errorf("Expected 5 categories, got %d", len(data))
	}
}

// TestGetCategoryByID_Success tests GET /categories/{id} with valid ID
func TestGetCategoryByID_Success(t *testing.T) {
	handler := setupTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
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

	if response.Message != "Category retrieved successfully" {
		t.Errorf("Expected message 'Category retrieved successfully', got '%s'", response.Message)
	}

	// Check category data
	data, ok := response.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be an object, got %T", response.Data)
	}

	if data["name"] != "Electronics" {
		t.Errorf("Expected name 'Electronics', got '%v'", data["name"])
	}
}

// TestGetCategoryByID_NotFound tests GET /categories/{id} with non-existent ID
func TestGetCategoryByID_NotFound(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
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

// TestGetCategoryByID_InvalidID tests GET /categories/{id} with invalid ID
func TestGetCategoryByID_InvalidID(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Invalid category ID" {
		t.Errorf("Expected message 'Invalid category ID', got '%s'", response.Message)
	}
}

// TestCreateCategory_Success tests POST /categories with valid data
func TestCreateCategory_Success(t *testing.T) {
	handler := setupTestHandler()

	category := models.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != "Category created successfully" {
		t.Errorf("Expected message 'Category created successfully', got '%s'", response.Message)
	}

	// Check response data
	data, ok := response.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be an object, got %T", response.Data)
	}

	if data["name"] != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%v'", data["name"])
	}
}

// TestCreateCategory_EmptyName tests POST /categories with empty name
func TestCreateCategory_EmptyName(t *testing.T) {
	handler := setupTestHandler()

	category := models.Category{
		Name:        "",
		Description: "Test Description",
	}

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Name is required" {
		t.Errorf("Expected message 'Name is required', got '%s'", response.Message)
	}
}

// TestCreateCategory_DuplicateName tests POST /categories with duplicate name
func TestCreateCategory_DuplicateName(t *testing.T) {
	handler := setupTestHandlerWithData()

	category := models.Category{
		Name:        "Electronics", // Already exists in seed data
		Description: "Duplicate",
	}

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Category name already exists" {
		t.Errorf("Expected message 'Category name already exists', got '%s'", response.Message)
	}
}

// TestCreateCategory_InvalidJSON tests POST /categories with invalid JSON
func TestCreateCategory_InvalidJSON(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Invalid request body" {
		t.Errorf("Expected message 'Invalid request body', got '%s'", response.Message)
	}
}

// TestUpdateCategory_Success tests PUT /categories/{id} with valid data
func TestUpdateCategory_Success(t *testing.T) {
	handler := setupTestHandlerWithData()

	category := models.Category{
		Name:        "Updated Electronics",
		Description: "Updated Description",
	}

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
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

	if response.Message != "Category updated successfully" {
		t.Errorf("Expected message 'Category updated successfully', got '%s'", response.Message)
	}

	data, ok := response.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be an object, got %T", response.Data)
	}

	if data["name"] != "Updated Electronics" {
		t.Errorf("Expected name 'Updated Electronics', got '%v'", data["name"])
	}
}

// TestUpdateCategory_NotFound tests PUT /categories/{id} with non-existent ID
func TestUpdateCategory_NotFound(t *testing.T) {
	handler := setupTestHandler()

	category := models.Category{
		Name:        "New Category",
		Description: "Description",
	}

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
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

// TestUpdateCategory_EmptyName tests PUT /categories/{id} with empty name
func TestUpdateCategory_EmptyName(t *testing.T) {
	handler := setupTestHandlerWithData()

	category := models.Category{
		Name:        "",
		Description: "Description",
	}

	body, _ := json.Marshal(category)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Name is required" {
		t.Errorf("Expected message 'Name is required', got '%s'", response.Message)
	}
}

// TestUpdateCategory_InvalidJSON tests PUT /categories/{id} with invalid JSON
func TestUpdateCategory_InvalidJSON(t *testing.T) {
	handler := setupTestHandlerWithData()

	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}
}

// TestDeleteCategory_Success tests DELETE /categories/{id} with valid ID
func TestDeleteCategory_Success(t *testing.T) {
	handler := setupTestHandlerWithData()

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
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

	if response.Message != "Category deleted successfully" {
		t.Errorf("Expected message 'Category deleted successfully', got '%s'", response.Message)
	}

	// Verify deletion - try to get the deleted category
	req2 := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusNotFound {
		t.Errorf("Expected deleted category to return %d, got %d", http.StatusNotFound, rec2.Code)
	}
}

// TestDeleteCategory_NotFound tests DELETE /categories/{id} with non-existent ID
func TestDeleteCategory_NotFound(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/categories/999", nil)
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

// TestMethodNotAllowed_Collection tests unsupported methods on /categories
func TestMethodNotAllowed_Collection(t *testing.T) {
	handler := setupTestHandler()

	unsupportedMethods := []string{http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range unsupportedMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/categories", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for method %s, got %d", http.StatusMethodNotAllowed, method, rec.Code)
			}

			var response Response
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Success {
				t.Error("Expected success to be false")
			}

			if response.Message != "Method not allowed" {
				t.Errorf("Expected message 'Method not allowed', got '%s'", response.Message)
			}
		})
	}
}

// TestMethodNotAllowed_Resource tests unsupported methods on /categories/{id}
func TestMethodNotAllowed_Resource(t *testing.T) {
	handler := setupTestHandlerWithData()

	unsupportedMethods := []string{http.MethodPost, http.MethodPatch}

	for _, method := range unsupportedMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/categories/1", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for method %s, got %d", http.StatusMethodNotAllowed, method, rec.Code)
			}

			var response Response
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Success {
				t.Error("Expected success to be false")
			}
		})
	}
}

// TestContentType tests that response has correct Content-Type header
func TestContentType(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestCRUDFlow tests a complete CRUD flow
func TestCRUDFlow(t *testing.T) {
	handler := setupTestHandler()

	// 1. Create a category
	createBody, _ := json.Marshal(models.Category{
		Name:        "Test Category",
		Description: "Test Description",
	})
	createReq := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	handler.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("Create failed: expected status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// 2. Get the created category
	getReq := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	getRec := httptest.NewRecorder()

	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("Get failed: expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	// 3. Update the category
	updateBody, _ := json.Marshal(models.Category{
		Name:        "Updated Category",
		Description: "Updated Description",
	})
	updateReq := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()

	handler.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("Update failed: expected status %d, got %d", http.StatusOK, updateRec.Code)
	}

	// 4. Verify the update
	verifyReq := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	verifyRec := httptest.NewRecorder()

	handler.ServeHTTP(verifyRec, verifyReq)

	var verifyResponse Response
	if err := json.NewDecoder(verifyRec.Body).Decode(&verifyResponse); err != nil {
		t.Fatalf("Failed to decode verify response: %v", err)
	}

	data := verifyResponse.Data.(map[string]any)
	if data["name"] != "Updated Category" {
		t.Errorf("Update not persisted: expected 'Updated Category', got '%v'", data["name"])
	}

	// 5. Delete the category
	deleteReq := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	deleteRec := httptest.NewRecorder()

	handler.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("Delete failed: expected status %d, got %d", http.StatusOK, deleteRec.Code)
	}

	// 6. Verify deletion
	finalReq := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	finalRec := httptest.NewRecorder()

	handler.ServeHTTP(finalRec, finalReq)

	if finalRec.Code != http.StatusNotFound {
		t.Errorf("Delete not persisted: expected status %d, got %d", http.StatusNotFound, finalRec.Code)
	}
}
