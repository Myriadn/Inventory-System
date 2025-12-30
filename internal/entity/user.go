package entity

import (
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan tabel users
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" agar tidak ikut ter-render ke JSON response
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session merepresentasikan tabel sessions
type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     uuid.UUID `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

// RegisterRequest payload untuk register
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=super_admin admin staff"`
}

// LoginRequest payload untuk login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse response setelah login sukses
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
	Role      string `json:"role"`
}
