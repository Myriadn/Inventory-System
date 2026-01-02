package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestReportRepository_GetDashboardStats(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewReportRepository(mock)

	// Urutan query harus sesuai dengan logic di repository

	// 1. Count Products
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM products`)).
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(50))

	// 2. Count Sales
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM sales`)).
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(20))

	// 3. Sum Revenue
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(total_amount), 0) FROM sales`)).
		WillReturnRows(mock.NewRows([]string{"sum"}).AddRow(1500000.0))

	stats, err := repo.GetDashboardStats(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 50, stats.TotalProducts)
	assert.Equal(t, 20, stats.TotalSales)
	assert.Equal(t, 1500000.0, stats.TotalRevenue)
}
