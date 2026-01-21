package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/KAnggara75/BelajarGolang/store"
)

type CategoryHandler struct {
	store *store.CategoryStore
}

func NewCategoryHandler(s *store.CategoryStore) *CategoryHandler {
	return &CategoryHandler{store: s}
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/categories")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		// Handle collection routes: GET /categories, POST /categories
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

	// Handle single resource routes: GET/PUT/DELETE /categories/{id}
	id, err := strconv.Atoi(path)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid category ID")
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

// GetAll returns all categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories := h.store.GetAll()
	h.sendSuccess(w, http.StatusOK, "Categories retrieved successfully", categories)
}

// GetByID returns a single category
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.store.GetByID(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Category not found")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Category retrieved successfully", category)
}

// Create adds a new category
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if cat.Name == "" {
		h.sendError(w, http.StatusBadRequest, "Name is required")
		return
	}

	created, err := h.store.Create(cat)
	if err != nil {
		if err == store.ErrNameExists {
			h.sendError(w, http.StatusConflict, "Category name already exists")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}
	h.sendSuccess(w, http.StatusCreated, "Category created successfully", created)
}

// Update updates an existing category
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if cat.Name == "" {
		h.sendError(w, http.StatusBadRequest, "Name is required")
		return
	}

	updated, err := h.store.Update(id, cat)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Category not found")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Category updated successfully", updated)
}

// Delete removes a category
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.store.Delete(id); err != nil {
		h.sendError(w, http.StatusNotFound, "Category not found")
		return
	}
	h.sendSuccess(w, http.StatusOK, "Category deleted successfully", nil)
}

func (h *CategoryHandler) sendSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *CategoryHandler) sendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}

func (h *CategoryHandler) methodNotAllowed(w http.ResponseWriter) {
	h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
}
