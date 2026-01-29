package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/KAnggara75/BelajarGolang/repository"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/products")
	path = strings.TrimPrefix(path, "/")

	// Check for query parameter to filter by category
	if path == "" && r.Method == http.MethodGet {
		categoryIDStr := r.URL.Query().Get("category_id")
		if categoryIDStr != "" {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				h.sendError(w, http.StatusBadRequest, "Invalid category_id parameter")
				return
			}
			h.GetByCategory(w, r, categoryID)
			return
		}
	}

	if path == "" {
		// Handle collection routes: GET /products, POST /products
		switch r.Method {
		case http.MethodGet:
			h.GetAll(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			h.methodNotAllowed(w)
		}
		return
	}

	// Handle single resource routes: GET/PUT/DELETE /products/{id}
	id, err := strconv.Atoi(path)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		h.methodNotAllowed(w)
	}
}

// GetAll returns all products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAll(r.Context())
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Products retrieved successfully", products)
}

// GetByCategory returns products filtered by category
func (h *ProductHandler) GetByCategory(w http.ResponseWriter, r *http.Request, categoryID int) {
	products, err := h.repo.GetByCategory(r.Context(), categoryID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Products retrieved successfully", products)
}

// GetByID returns a single product
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	product, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			h.sendError(w, http.StatusNotFound, "Product not found")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Product retrieved successfully", product)
}

// Create adds a new product
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		h.sendError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if input.Price < 0 {
		h.sendError(w, http.StatusBadRequest, "Price cannot be negative")
		return
	}

	if input.Stock < 0 {
		h.sendError(w, http.StatusBadRequest, "Stock cannot be negative")
		return
	}

	product := input.ToProduct()

	created, err := h.repo.Create(r.Context(), product)
	if err != nil {
		if err == repository.ErrProductNameExists {
			h.sendError(w, http.StatusConflict, "Product name already exists")
			return
		}
		if err == repository.ErrProductCategoryNotFound {
			h.sendError(w, http.StatusBadRequest, "Category not found")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}
	h.sendSuccess(w, http.StatusCreated, "Product created successfully", created)
}

// Update updates an existing product
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var input models.ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		h.sendError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if input.Price < 0 {
		h.sendError(w, http.StatusBadRequest, "Price cannot be negative")
		return
	}

	if input.Stock < 0 {
		h.sendError(w, http.StatusBadRequest, "Stock cannot be negative")
		return
	}

	product := input.ToProduct()

	updated, err := h.repo.Update(r.Context(), id, product)
	if err != nil {
		if err == repository.ErrProductNotFound {
			h.sendError(w, http.StatusNotFound, "Product not found")
			return
		}
		if err == repository.ErrProductCategoryNotFound {
			h.sendError(w, http.StatusBadRequest, "Category not found")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Product updated successfully", updated)
}

// Delete removes a product
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err == repository.ErrProductNotFound {
			h.sendError(w, http.StatusNotFound, "Product not found")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) sendSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *ProductHandler) sendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}

func (h *ProductHandler) methodNotAllowed(w http.ResponseWriter) {
	h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
}
