package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/KAnggara75/BelajarGolang/repository"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	repo repository.ReportRepository
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(repo repository.ReportRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

// ServeHTTP implements the http.Handler interface
func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/report")
	path = strings.TrimPrefix(path, "/")

	// GET /api/report/hari-ini - Get today's report
	if path == "hari-ini" {
		h.GetTodayReport(w, r)
		return
	}

	// GET /api/report?start_date=...&end_date=... - Get period report
	if path == "" {
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		if startDateStr != "" && endDateStr != "" {
			h.GetPeriodReport(w, r, startDateStr, endDateStr)
			return
		}

		h.sendError(w, http.StatusBadRequest, "Missing start_date or end_date query parameters")
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// GetTodayReport returns today's sales report
func (h *ReportHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	// Get current time in local timezone
	now := time.Now()

	report, err := h.repo.GetDailyReport(r.Context(), now)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve report")
		return
	}

	h.sendJSON(w, http.StatusOK, report)
}

// GetPeriodReport returns sales report for a date range
func (h *ReportHandler) GetPeriodReport(w http.ResponseWriter, r *http.Request, startDateStr, endDateStr string) {
	// Parse dates (format: YYYY-MM-DD)
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		h.sendError(w, http.StatusBadRequest, "end_date must be after start_date")
		return
	}

	report, err := h.repo.GetPeriodReport(r.Context(), startDate, endDate)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to retrieve report")
		return
	}

	h.sendJSON(w, http.StatusOK, report)
}

// sendJSON sends a JSON response
func (h *ReportHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func (h *ReportHandler) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}
