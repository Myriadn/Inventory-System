package entity

import "time"

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request Payload untuk Create/Update
type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description"`
}
