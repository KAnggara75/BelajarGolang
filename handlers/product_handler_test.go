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

// mockProductRepository is a mock implementation of ProductRepository for testing
type mockProductRepository struct {
	products   map[int]models.Product
	categories map[int]models.Category
	nextID     int
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products:   make(map[int]models.Product),
		categories: make(map[int]models.Category),
		nextID:     1,
	}
}

func (m *mockProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	result := make([]models.Product, 0, len(m.products))
	for _, p := range m.products {
		// Attach category if exists
		if p.CategoryID > 0 {
			if cat, ok := m.categories[p.CategoryID]; ok {
				p.Category = &cat
			}
		}
		result = append(result, p)
	}
	return result, nil
}

func (m *mockProductRepository) GetByID(ctx context.Context, id int) (models.Product, error) {
	p, exists := m.products[id]
	if !exists {
		return models.Product{}, repository.ErrProductNotFound
	}
	// Attach category if exists
	if p.CategoryID > 0 {
		if cat, ok := m.categories[p.CategoryID]; ok {
			p.Category = &cat
		}
	}
	return p, nil
}

func (m *mockProductRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error) {
	result := make([]models.Product, 0)
	for _, p := range m.products {
		if p.CategoryID == categoryID {
			if cat, ok := m.categories[p.CategoryID]; ok {
				p.Category = &cat
			}
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockProductRepository) CategoryExists(ctx context.Context, categoryID int) (bool, error) {
	_, exists := m.categories[categoryID]
	return exists, nil
}

func (m *mockProductRepository) Create(ctx context.Context, p models.Product) (models.Product, error) {
	// Check if name already exists
	for _, existing := range m.products {
		if existing.Name == p.Name {
			return models.Product{}, repository.ErrProductNameExists
		}
	}

	// Check if category exists (if specified)
	if p.CategoryID > 0 {
		if _, exists := m.categories[p.CategoryID]; !exists {
			return models.Product{}, repository.ErrProductCategoryNotFound
		}
	}

	p.ID = m.nextID
	m.nextID++
	m.products[p.ID] = p
	return p, nil
}

func (m *mockProductRepository) Update(ctx context.Context, id int, p models.Product) (models.Product, error) {
	if _, exists := m.products[id]; !exists {
		return models.Product{}, repository.ErrProductNotFound
	}

	// Check if category exists (if specified)
	if p.CategoryID > 0 {
		if _, exists := m.categories[p.CategoryID]; !exists {
			return models.Product{}, repository.ErrProductCategoryNotFound
		}
	}

	p.ID = id
	m.products[id] = p
	return p, nil
}

func (m *mockProductRepository) Delete(ctx context.Context, id int) error {
	if _, exists := m.products[id]; !exists {
		return repository.ErrProductNotFound
	}

	delete(m.products, id)
	return nil
}

// SeedCategories adds sample categories for testing
func (m *mockProductRepository) SeedCategories() {
	m.categories[1] = models.Category{ID: 1, Name: "Electronics", Description: "Electronic devices"}
	m.categories[2] = models.Category{ID: 2, Name: "Clothing", Description: "Apparel items"}
	m.categories[3] = models.Category{ID: 3, Name: "Books", Description: "Books and reading"}
}

// SeedData adds sample data for testing
func (m *mockProductRepository) SeedData() {
	m.SeedCategories()
	initialData := []models.Product{
		{Name: "iPhone 15 Pro", Price: 999.99, Stock: 50, CategoryID: 1},
		{Name: "MacBook Pro M3", Price: 2499.99, Stock: 25, CategoryID: 1},
		{Name: "AirPods Pro", Price: 249.99, Stock: 100, CategoryID: 1},
		{Name: "iPad Air", Price: 599.99, Stock: 40, CategoryID: 1},
		{Name: "Apple Watch Series 9", Price: 399.99, Stock: 60, CategoryID: 1},
	}

	for _, p := range initialData {
		_, _ = m.Create(context.Background(), p)
	}
}

// setupProductTestHandler creates a fresh handler with an empty mock repository for testing
func setupProductTestHandler() *ProductHandler {
	repo := newMockProductRepository()
	repo.SeedCategories() // Always seed categories
	return NewProductHandler(repo)
}

// setupProductTestHandlerWithData creates a handler with seeded data
func setupProductTestHandlerWithData() *ProductHandler {
	repo := newMockProductRepository()
	repo.SeedData()
	return NewProductHandler(repo)
}

// TestGetAllProducts_Empty tests GET /products with empty repo
func TestGetAllProducts_Empty(t *testing.T) {
	handler := setupProductTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
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

	// Data should be an empty array
	data, ok := response.Data.([]any)
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", response.Data)
	}
	if len(data) != 0 {
		t.Errorf("Expected 0 products, got %d", len(data))
	}
}

// TestGetAllProducts_WithData tests GET /products with seeded data
func TestGetAllProducts_WithData(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
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
		t.Errorf("Expected 5 products, got %d", len(data))
	}
}

// TestGetProductsByCategory tests GET /products?category_id=1
func TestGetProductsByCategory(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/products?category_id=1", nil)
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
		t.Errorf("Expected 5 products in category 1, got %d", len(data))
	}
}

// TestGetProductsByCategory_InvalidCategoryID tests GET /products with invalid category_id
func TestGetProductsByCategory_InvalidCategoryID(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/products?category_id=abc", nil)
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

	if response.Message != "Invalid category_id parameter" {
		t.Errorf("Expected message 'Invalid category_id parameter', got '%s'", response.Message)
	}
}

// TestGetProductByID_Success tests GET /products/{id} with valid ID
func TestGetProductByID_Success(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
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

	if response.Message != "Product retrieved successfully" {
		t.Errorf("Expected message 'Product retrieved successfully', got '%s'", response.Message)
	}

	// Check product data
	data, ok := response.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be an object, got %T", response.Data)
	}

	if data["name"] != "iPhone 15 Pro" {
		t.Errorf("Expected name 'iPhone 15 Pro', got '%v'", data["name"])
	}

	// Check category is included
	if data["category"] == nil {
		t.Error("Expected category to be included")
	}
}

// TestGetProductByID_NotFound tests GET /products/{id} with non-existent ID
func TestGetProductByID_NotFound(t *testing.T) {
	handler := setupProductTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
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

// TestGetProductByID_InvalidID tests GET /products/{id} with invalid ID
func TestGetProductByID_InvalidID(t *testing.T) {
	handler := setupProductTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
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

	if response.Message != "Invalid product ID" {
		t.Errorf("Expected message 'Invalid product ID', got '%s'", response.Message)
	}
}

// TestCreateProduct_Success tests POST /products with valid data including category
func TestCreateProduct_Success(t *testing.T) {
	handler := setupProductTestHandler()

	product := models.ProductInput{
		Name:       "Test Product",
		Price:      99.99,
		Stock:      10,
		CategoryID: 1, // Electronics
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
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

	if response.Message != "Product created successfully" {
		t.Errorf("Expected message 'Product created successfully', got '%s'", response.Message)
	}

	// Check response data
	data, ok := response.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be an object, got %T", response.Data)
	}

	if data["name"] != "Test Product" {
		t.Errorf("Expected name 'Test Product', got '%v'", data["name"])
	}
}

// TestCreateProduct_InvalidCategory tests POST /products with non-existent category
func TestCreateProduct_InvalidCategory(t *testing.T) {
	handler := setupProductTestHandler()

	product := models.ProductInput{
		Name:       "Test Product",
		Price:      99.99,
		Stock:      10,
		CategoryID: 999, // Non-existent category
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
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

	if response.Message != "Category not found" {
		t.Errorf("Expected message 'Category not found', got '%s'", response.Message)
	}
}

// TestCreateProduct_EmptyName tests POST /products with empty name
func TestCreateProduct_EmptyName(t *testing.T) {
	handler := setupProductTestHandler()

	product := models.ProductInput{
		Name:  "",
		Price: 99.99,
		Stock: 10,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
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

// TestCreateProduct_NegativePrice tests POST /products with negative price
func TestCreateProduct_NegativePrice(t *testing.T) {
	handler := setupProductTestHandler()

	product := models.ProductInput{
		Name:  "Test Product",
		Price: -10.00,
		Stock: 10,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
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

	if response.Message != "Price cannot be negative" {
		t.Errorf("Expected message 'Price cannot be negative', got '%s'", response.Message)
	}
}

// TestCreateProduct_NegativeStock tests POST /products with negative stock
func TestCreateProduct_NegativeStock(t *testing.T) {
	handler := setupProductTestHandler()

	product := models.ProductInput{
		Name:  "Test Product",
		Price: 99.99,
		Stock: -5,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
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

	if response.Message != "Stock cannot be negative" {
		t.Errorf("Expected message 'Stock cannot be negative', got '%s'", response.Message)
	}
}

// TestCreateProduct_DuplicateName tests POST /products with duplicate name
func TestCreateProduct_DuplicateName(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	product := models.ProductInput{
		Name:       "iPhone 15 Pro", // Already exists in seed data
		Price:      999.99,
		Stock:      10,
		CategoryID: 1,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
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

	if response.Message != "Product name already exists" {
		t.Errorf("Expected message 'Product name already exists', got '%s'", response.Message)
	}
}

// TestCreateProduct_InvalidJSON tests POST /products with invalid JSON
func TestCreateProduct_InvalidJSON(t *testing.T) {
	handler := setupProductTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("{invalid json}"))
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

// TestUpdateProduct_Success tests PUT /products/{id} with valid data
func TestUpdateProduct_Success(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	product := models.ProductInput{
		Name:       "Updated iPhone",
		Price:      1099.99,
		Stock:      75,
		CategoryID: 2, // Change to Clothing
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
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

	if response.Message != "Product updated successfully" {
		t.Errorf("Expected message 'Product updated successfully', got '%s'", response.Message)
	}

	data, ok := response.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be an object, got %T", response.Data)
	}

	if data["name"] != "Updated iPhone" {
		t.Errorf("Expected name 'Updated iPhone', got '%v'", data["name"])
	}
}

// TestUpdateProduct_InvalidCategory tests PUT /products/{id} with invalid category
func TestUpdateProduct_InvalidCategory(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	product := models.ProductInput{
		Name:       "Updated iPhone",
		Price:      1099.99,
		Stock:      75,
		CategoryID: 999, // Non-existent
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
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

	if response.Message != "Category not found" {
		t.Errorf("Expected message 'Category not found', got '%s'", response.Message)
	}
}

// TestUpdateProduct_NotFound tests PUT /products/{id} with non-existent ID
func TestUpdateProduct_NotFound(t *testing.T) {
	handler := setupProductTestHandler()

	product := models.ProductInput{
		Name:  "New Product",
		Price: 99.99,
		Stock: 10,
	}

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/products/999", bytes.NewBuffer(body))
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

	if response.Message != "Product not found" {
		t.Errorf("Expected message 'Product not found', got '%s'", response.Message)
	}
}

// TestDeleteProduct_Success tests DELETE /products/{id} with valid ID
func TestDeleteProduct_Success(t *testing.T) {
	handler := setupProductTestHandlerWithData()

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
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

	if response.Message != "Product deleted successfully" {
		t.Errorf("Expected message 'Product deleted successfully', got '%s'", response.Message)
	}

	// Verify deletion - try to get the deleted product
	req2 := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusNotFound {
		t.Errorf("Expected deleted product to return %d, got %d", http.StatusNotFound, rec2.Code)
	}
}

// TestDeleteProduct_NotFound tests DELETE /products/{id} with non-existent ID
func TestDeleteProduct_NotFound(t *testing.T) {
	handler := setupProductTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/products/999", nil)
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

// TestProductMethodNotAllowed_Collection tests unsupported methods on /products
func TestProductMethodNotAllowed_Collection(t *testing.T) {
	handler := setupProductTestHandler()

	unsupportedMethods := []string{http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range unsupportedMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/products", nil)
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

// TestProductCRUDFlow tests a complete CRUD flow for products with category
func TestProductCRUDFlow(t *testing.T) {
	handler := setupProductTestHandler()

	// 1. Create a product with category
	createBody, _ := json.Marshal(models.ProductInput{
		Name:       "Test Product",
		Price:      99.99,
		Stock:      10,
		CategoryID: 1,
	})
	createReq := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	handler.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("Create failed: expected status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// 2. Get the created product
	getReq := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	getRec := httptest.NewRecorder()

	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("Get failed: expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	// 3. Update the product with new category
	updateBody, _ := json.Marshal(models.ProductInput{
		Name:       "Updated Product",
		Price:      199.99,
		Stock:      20,
		CategoryID: 2, // Change category
	})
	updateReq := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()

	handler.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("Update failed: expected status %d, got %d", http.StatusOK, updateRec.Code)
	}

	// 4. Verify the update
	verifyReq := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	verifyRec := httptest.NewRecorder()

	handler.ServeHTTP(verifyRec, verifyReq)

	var verifyResponse Response
	if err := json.NewDecoder(verifyRec.Body).Decode(&verifyResponse); err != nil {
		t.Fatalf("Failed to decode verify response: %v", err)
	}

	data := verifyResponse.Data.(map[string]any)
	if data["name"] != "Updated Product" {
		t.Errorf("Update not persisted: expected 'Updated Product', got '%v'", data["name"])
	}

	// 5. Delete the product
	deleteReq := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	deleteRec := httptest.NewRecorder()

	handler.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("Delete failed: expected status %d, got %d", http.StatusOK, deleteRec.Code)
	}

	// 6. Verify deletion
	finalReq := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	finalRec := httptest.NewRecorder()

	handler.ServeHTTP(finalRec, finalReq)

	if finalRec.Code != http.StatusNotFound {
		t.Errorf("Delete not persisted: expected status %d, got %d", http.StatusNotFound, finalRec.Code)
	}
}
