package entity

import "time"

type Rack struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request Payload untuk Create/Update
type RackRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Category string `json:"category" validate:"required,min=3,max=50"`
}
