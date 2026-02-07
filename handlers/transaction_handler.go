package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/KAnggara75/BelajarGolang/repository"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	repo repository.TransactionRepository
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

// ServeHTTP implements the http.Handler interface
func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/transactions")
	path = strings.TrimPrefix(path, "/")

	// POST /transactions/checkout - Create a new transaction
	if path == "checkout" && r.Method == http.MethodPost {
		h.Checkout(w, r)
		return
	}

	// GET /transactions - Get all transactions
	if path == "" && r.Method == http.MethodGet {
		h.GetAll(w, r)
		return
	}

	// GET /transactions/{id} - Get transaction by ID
	if path != "" && r.Method == http.MethodGet {
		id, err := strconv.Atoi(path)
		if err != nil {
			h.sendError(w, http.StatusBadRequest, "Invalid transaction ID")
			return
		}
		h.GetByID(w, r, id)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// Checkout processes a checkout request
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, err := h.repo.Checkout(r.Context(), req)
	if err != nil {
		switch err {
		case repository.ErrEmptyCheckout:
			h.sendError(w, http.StatusBadRequest, "Checkout items cannot be empty")
		case repository.ErrInvalidQuantity:
			h.sendError(w, http.StatusBadRequest, "Invalid quantity: must be greater than 0")
		case repository.ErrProductNotFoundInTx:
			h.sendError(w, http.StatusNotFound, "One or more products not found")
		case repository.ErrInsufficientStock:
			h.sendError(w, http.StatusBadRequest, err.Error())
		default:
			if strings.Contains(err.Error(), "insufficient stock") {
				h.sendError(w, http.StatusBadRequest, err.Error())
			} else {
				h.sendError(w, http.StatusInternalServerError, "Failed to process checkout")
			}
		}
		return
	}

	h.sendSuccess(w, http.StatusCreated, "Transaction created successfully", transaction)
}

// GetAll returns all transactions
func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.repo.GetAll(r.Context())
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	h.sendSuccess(w, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// GetByID returns a transaction by ID
func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	transaction, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			h.sendError(w, http.StatusNotFound, "Transaction not found")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve transaction")
		return
	}

	h.sendSuccess(w, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// sendSuccess sends a success response
func (h *TransactionHandler) sendSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// sendError sends an error response
func (h *TransactionHandler) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}
