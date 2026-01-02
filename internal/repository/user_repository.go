package repository

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db DBExecutor
}

// Struct bantuan untuk hasil query join
type SessionResult struct {
	UserID    int64
	Role      string
	ExpiredAt time.Time
}

func NewUserRepository(db DBExecutor) *UserRepository {
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

// FindAll mengambil semua user dengan pagination
func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.User, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users ORDER BY id DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		// Scan field yang diperlukan (Password hash tidak kita kembalikan demi keamanan)
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// FindByID mencari user berdasarkan ID
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`
	var u entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// UpdateRole mengubah role user
func (r *UserRepository) UpdateRole(ctx context.Context, id int64, newRole string) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	cmd, err := r.db.Exec(ctx, query, newRole, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows // User tidak ditemukan
	}
	return nil
}

// Delete menghapus user secara permanen
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
