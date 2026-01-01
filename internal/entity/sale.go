package entity

import "time"

type Sale struct {
	ID              int64        `json:"id"`
	UserID          int64        `json:"user_id"` // Kasir
	TransactionDate time.Time    `json:"transaction_date"`
	TotalAmount     float64      `json:"total_amount"`
	Details         []SaleDetail `json:"details,omitempty"` // Opsional ditampilkan
	CreatedAt       time.Time    `json:"created_at"`
}

type SaleDetail struct {
	ID        int64   `json:"id"`
	SaleID    int64   `json:"sale_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"` // Harga saat transaksi
	SubTotal  float64 `json:"sub_total"`
}

// Payload Request dari Frontend
type SaleRequest struct {
	Items []SaleItemRequest `json:"items" validate:"required,min=1,dive"`
}

type SaleItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,min=1"`
}
