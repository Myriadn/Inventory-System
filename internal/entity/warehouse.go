package entity

import "time"

type Warehouse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request Payload untuk Create/Update
type WarehouseRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Location    string `json:"location" validate:"required,min=3,max=255"`
	Description string `json:"description"`
}
