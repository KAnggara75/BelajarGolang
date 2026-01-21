package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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
		// Handle collection routes: GET /categories
		switch r.Method {
		case http.MethodGet:
			h.GetAll(w, r)
		default:
			h.methodNotAllowed(w)
		}
		return
	}
}

// GetAll returns all categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories := h.store.GetAll()
	h.sendSuccess(w, http.StatusOK, "Categories retrieved successfully", categories)
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
