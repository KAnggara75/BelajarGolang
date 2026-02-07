package repository

import (
	"context"
	"time"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/jackc/pgx/v5"
)

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	GetDailyReport(ctx context.Context, date time.Time) (models.ReportResponse, error)
	GetPeriodReport(ctx context.Context, startDate, endDate time.Time) (models.ReportResponse, error)
}

// reportRepository implements ReportRepository using PostgreSQL
type reportRepository struct {
	db *pgx.Conn
}

// NewReportRepository creates a new ReportRepository
func NewReportRepository(db *pgx.Conn) ReportRepository {
	return &reportRepository{db: db}
}

// GetDailyReport returns report for a specific date
func (r *reportRepository) GetDailyReport(ctx context.Context, date time.Time) (models.ReportResponse, error) {
	// Set start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	return r.getReport(ctx, startOfDay, endOfDay)
}

// GetPeriodReport returns report for a date range
func (r *reportRepository) GetPeriodReport(ctx context.Context, startDate, endDate time.Time) (models.ReportResponse, error) {
	// Set to start of start date and end of end date
	startOfPeriod := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endOfPeriod := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return r.getReport(ctx, startOfPeriod, endOfPeriod)
}

// getReport is a helper function to get report for any date range
func (r *reportRepository) getReport(ctx context.Context, startDate, endDate time.Time) (models.ReportResponse, error) {
	var report models.ReportResponse

	// Get total revenue and transaction count
	revenueQuery := `
		SELECT
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2
	`

	err := r.db.QueryRow(ctx, revenueQuery, startDate, endDate).Scan(
		&report.TotalRevenue,
		&report.TotalTransaksi,
	)
	if err != nil {
		return models.ReportResponse{}, err
	}

	// Get best-selling product
	topProductQuery := `
		SELECT
			p.name,
			SUM(td.quantity) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`

	var productName string
	var qtyTerjual int

	err = r.db.QueryRow(ctx, topProductQuery, startDate, endDate).Scan(&productName, &qtyTerjual)
	if err != nil {
		// If no products found (no transactions), that's okay
		if err == pgx.ErrNoRows {
			report.ProdukTerlaris = nil
			return report, nil
		}
		return models.ReportResponse{}, err
	}

	report.ProdukTerlaris = &models.TopProduct{
		Nama:       productName,
		QtyTerjual: qtyTerjual,
	}

	return report, nil
}
