package entity

type DashboardReport struct {
	TotalProducts int     `json:"total_products"`
	TotalSales    int     `json:"total_sales"`   // Jumlah transaksi
	TotalRevenue  float64 `json:"total_revenue"` // Total duit masuk
}
