package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/KAnggara75/BelajarGolang/repository"
)

// Mock transaction repository
type mockTransactionRepository struct {
	transactions map[int]models.Transaction
	nextID       int
	products     map[int]mockProduct
}

type mockProduct struct {
	id    int
	name  string
	price float64
	stock int
}

func newMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{
		transactions: make(map[int]models.Transaction),
		nextID:       1,
		products: map[int]mockProduct{
			1: {id: 1, name: "iPhone 15 Pro", price: 999.99, stock: 50},
			2: {id: 2, name: "MacBook Pro M3", price: 2499.99, stock: 25},
			3: {id: 3, name: "AirPods Pro", price: 249.99, stock: 100},
		},
	}
}

func (m *mockTransactionRepository) Checkout(ctx context.Context, req models.CheckoutRequest) (models.Transaction, error) {
	if len(req.Items) == 0 {
		return models.Transaction{}, repository.ErrEmptyCheckout
	}

	var totalAmount int
	var details []models.TransactionDetail

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return models.Transaction{}, repository.ErrInvalidQuantity
		}

		product, exists := m.products[item.ProductID]
		if !exists {
			return models.Transaction{}, repository.ErrProductNotFoundInTx
		}

		if product.stock < item.Quantity {
			return models.Transaction{}, repository.ErrInsufficientStock
		}

		subtotal := int(product.price*100) * item.Quantity
		totalAmount += subtotal

		// Update stock
		product.stock -= item.Quantity
		m.products[item.ProductID] = product

		details = append(details, models.TransactionDetail{
			ID:          len(details) + 1,
			ProductID:   item.ProductID,
			ProductName: product.name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	transaction := models.Transaction{
		ID:          m.nextID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}

	for i := range transaction.Details {
		transaction.Details[i].TransactionID = transaction.ID
	}

	m.transactions[m.nextID] = transaction
	m.nextID++

	return transaction, nil
}

func (m *mockTransactionRepository) GetAll(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	for _, tx := range m.transactions {
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (m *mockTransactionRepository) GetByID(ctx context.Context, id int) (models.Transaction, error) {
	tx, exists := m.transactions[id]
	if !exists {
		return models.Transaction{}, repository.ErrTransactionNotFound
	}
	return tx, nil
}

// Tests

func TestCheckout_Success(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 3, Quantity: 1},
		},
	}

	body, _ := json.Marshal(checkoutReq)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewReader(body))
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

	if response.Message != "Transaction created successfully" {
		t.Errorf("Expected message 'Transaction created successfully', got '%s'", response.Message)
	}

	// Verify stock was decremented
	if repo.products[1].stock != 48 {
		t.Errorf("Expected product 1 stock to be 48, got %d", repo.products[1].stock)
	}
	if repo.products[3].stock != 99 {
		t.Errorf("Expected product 3 stock to be 99, got %d", repo.products[3].stock)
	}
}

func TestCheckout_EmptyItems(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{},
	}

	body, _ := json.Marshal(checkoutReq)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCheckout_InvalidQuantity(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 0},
		},
	}

	body, _ := json.Marshal(checkoutReq)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCheckout_InsufficientStock(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 100}, // Only 50 in stock
		},
	}

	body, _ := json.Marshal(checkoutReq)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCheckout_ProductNotFound(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 999, Quantity: 1},
		},
	}

	body, _ := json.Marshal(checkoutReq)
	req := httptest.NewRequest(http.MethodPost, "/transactions/checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetAllTransactions(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	// Create a transaction first
	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 1},
		},
	}
	repo.Checkout(context.Background(), checkoutReq)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
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
}

func TestGetTransactionByID_Success(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	// Create a transaction first
	checkoutReq := models.CheckoutRequest{
		Items: []models.CheckoutItem{
			{ProductID: 1, Quantity: 1},
		},
	}
	repo.Checkout(context.Background(), checkoutReq)

	req := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
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
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	repo := newMockTransactionRepository()
	handler := NewTransactionHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/transactions/999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
