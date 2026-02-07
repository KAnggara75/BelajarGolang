package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KAnggara75/BelajarGolang/models"
)

// Mock report repository
type mockReportRepository struct {
	dailyReport  models.ReportResponse
	periodReport models.ReportResponse
}

func newMockReportRepository() *mockReportRepository {
	return &mockReportRepository{
		dailyReport: models.ReportResponse{
			TotalRevenue:   45000,
			TotalTransaksi: 5,
			ProdukTerlaris: &models.TopProduct{
				Nama:       "Indomie Goreng",
				QtyTerjual: 12,
			},
		},
		periodReport: models.ReportResponse{
			TotalRevenue:   150000,
			TotalTransaksi: 15,
			ProdukTerlaris: &models.TopProduct{
				Nama:       "Indomie Goreng",
				QtyTerjual: 45,
			},
		},
	}
}

func (m *mockReportRepository) GetDailyReport(ctx context.Context, date time.Time) (models.ReportResponse, error) {
	return m.dailyReport, nil
}

func (m *mockReportRepository) GetPeriodReport(ctx context.Context, startDate, endDate time.Time) (models.ReportResponse, error) {
	return m.periodReport, nil
}

// Tests

func TestGetTodayReport_Success(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report/hari-ini", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var report models.ReportResponse
	if err := json.NewDecoder(rec.Body).Decode(&report); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if report.TotalRevenue != 45000 {
		t.Errorf("Expected total_revenue 45000, got %d", report.TotalRevenue)
	}

	if report.TotalTransaksi != 5 {
		t.Errorf("Expected total_transaksi 5, got %d", report.TotalTransaksi)
	}

	if report.ProdukTerlaris == nil {
		t.Fatal("Expected produk_terlaris to be non-nil")
	}

	if report.ProdukTerlaris.Nama != "Indomie Goreng" {
		t.Errorf("Expected produk_terlaris.nama 'Indomie Goreng', got '%s'", report.ProdukTerlaris.Nama)
	}

	if report.ProdukTerlaris.QtyTerjual != 12 {
		t.Errorf("Expected produk_terlaris.qty_terjual 12, got %d", report.ProdukTerlaris.QtyTerjual)
	}
}

func TestGetPeriodReport_Success(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2026-01-01&end_date=2026-02-01", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var report models.ReportResponse
	if err := json.NewDecoder(rec.Body).Decode(&report); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if report.TotalRevenue != 150000 {
		t.Errorf("Expected total_revenue 150000, got %d", report.TotalRevenue)
	}

	if report.TotalTransaksi != 15 {
		t.Errorf("Expected total_transaksi 15, got %d", report.TotalTransaksi)
	}

	if report.ProdukTerlaris == nil {
		t.Fatal("Expected produk_terlaris to be non-nil")
	}
}

func TestGetPeriodReport_MissingParameters(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetPeriodReport_InvalidStartDate(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=invalid&end_date=2026-02-01", nil)
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

func TestGetPeriodReport_InvalidEndDate(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2026-01-01&end_date=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetPeriodReport_EndDateBeforeStartDate(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2026-02-01&end_date=2026-01-01", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !contains(response.Message, "end_date must be after start_date") {
		t.Errorf("Expected error message about date order, got '%s'", response.Message)
	}
}

func TestReportHandler_MethodNotAllowed(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/report/hari-ini", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestReportHandler_NotFound(t *testing.T) {
	repo := newMockReportRepository()
	handler := NewReportHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/report/invalid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
