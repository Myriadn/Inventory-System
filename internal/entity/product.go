package entity

import "time"

type Product struct {
	ID          int64     `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	Price       float64   `json:"price"`
	CategoryID  int64     `json:"category_id"`
	RackID      int64     `json:"rack_id"`
	WarehouseID int64     `json:"warehouse_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRequest struct {
	SKU         string  `json:"sku" validate:"required,min=3,max=50"`
	Name        string  `json:"name" validate:"required,min=3,max=150"`
	Description string  `json:"description"`
	Stock       int     `json:"stock" validate:"min=0"` // Stock gak boleh minus
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryID  int64   `json:"category_id" validate:"required"`
	RackID      int64   `json:"rack_id" validate:"required"`
	WarehouseID int64   `json:"warehouse_id" validate:"required"`
}
