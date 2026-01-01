package repository

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetDashboardStats(ctx context.Context) (*entity.DashboardReport, error) {
	var report entity.DashboardReport

	// Hitung Total Produk
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&report.TotalProducts)
	if err != nil {
		return nil, err
	}

	// Hitung Total Transaksi Penjualan
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM sales").Scan(&report.TotalSales)
	if err != nil {
		return nil, err
	}

	// Hitung Total Revenue (Handle jika NULL/belum ada penjualan)
	// COALESCE(SUM(...), 0) artinya jika hasilnya null, ganti jadi 0.
	err = r.db.QueryRow(ctx, "SELECT COALESCE(SUM(total_amount), 0) FROM sales").Scan(&report.TotalRevenue)
	if err != nil {
		return nil, err
	}

	return &report, nil
}
