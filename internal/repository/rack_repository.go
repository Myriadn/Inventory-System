package repository

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RackRepository struct {
	db *pgxpool.Pool
}

func NewRackRepository(db *pgxpool.Pool) *RackRepository {
	return &RackRepository{db: db}
}

func (r *RackRepository) Create(ctx context.Context, rack *entity.Rack) error {
	query := `INSERT INTO racks (name, category, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, rack.Name, rack.Category).Scan(&rack.ID, &rack.CreatedAt, &rack.UpdatedAt)
}

func (r *RackRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Rack, error) {
	query := `SELECT id, name, category, created_at, updated_at FROM racks ORDER BY id DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racks []entity.Rack
	for rows.Next() {
		var rk entity.Rack
		if err := rows.Scan(&rk.ID, &rk.Name, &rk.Category, &rk.CreatedAt, &rk.UpdatedAt); err != nil {
			return nil, err
		}
		racks = append(racks, rk)
	}
	return racks, nil
}

func (r *RackRepository) FindByID(ctx context.Context, id int64) (*entity.Rack, error) {
	query := `SELECT id, name, category, created_at, updated_at FROM racks WHERE id = $1`
	var rk entity.Rack
	err := r.db.QueryRow(ctx, query, id).Scan(&rk.ID, &rk.Name, &rk.Category, &rk.CreatedAt, &rk.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rk, nil
}

func (r *RackRepository) Update(ctx context.Context, rack *entity.Rack) error {
	query := `UPDATE racks SET name = $1, category = $2, updated_at = NOW() WHERE id = $3`
	cmd, err := r.db.Exec(ctx, query, rack.Name, rack.Category, rack.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *RackRepository) Delete(ctx context.Context, id int64) error {
	// Ingat: Database pakai ON DELETE RESTRICT di produk.
	// Jadi kalau kategori ini dipakai produk, delete ini akan gagal (Error FK).
	query := `DELETE FROM racks WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
