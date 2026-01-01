package repository

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	query := `INSERT INTO categories (name, description, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, category.Name, category.Description).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (r *CategoryRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY id DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*entity.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	var c entity.Category
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	query := `UPDATE categories SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`
	cmd, err := r.db.Exec(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	// Ingat: Database pakai ON DELETE RESTRICT di produk.
	// Jadi kalau kategori ini dipakai produk, delete ini akan gagal (Error FK).
	query := `DELETE FROM categories WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
