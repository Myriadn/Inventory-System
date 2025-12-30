package repository

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser menyimpan user baru ke database
func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role, time.Now()).Scan(&user.ID)
	return err
}

// GetUserByEmail mencari user berdasarkan email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, username, email, password_hash, role FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User tidak ditemukan
		}
		return nil, err
	}
	return &user, nil
}

// CreateSession menyimpan session token
func (r *UserRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token, expired_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := r.db.Exec(ctx, query, session.ID, session.UserID, session.Token, session.ExpiredAt)
	return err
}

// Struct bantuan untuk hasil query join
type SessionResult struct {
	UserID    int64
	Role      string
	ExpiredAt time.Time
}

func (r *UserRepository) GetSessionByToken(ctx context.Context, token string) (*SessionResult, error) {
	query := `
		SELECT s.user_id, s.expired_at, u.role
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.token = $1 AND s.is_revoked = FALSE
	`

	var result SessionResult
	err := r.db.QueryRow(ctx, query, token).Scan(&result.UserID, &result.ExpiredAt, &result.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Token tidak valid atau sudah direvoke
		}
		return nil, err
	}

	return &result, nil
}
