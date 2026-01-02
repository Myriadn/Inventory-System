package repository

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5"
)

type WarehouseRepository struct {
	db DBExecutor
}

func NewWarehouseRepository(db DBExecutor) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) Create(ctx context.Context, warehouse *entity.Warehouse) error {
	query := `INSERT INTO warehouses (name, location, description, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, warehouse.Name, warehouse.Location, warehouse.Description).Scan(&warehouse.ID, &warehouse.CreatedAt, &warehouse.UpdatedAt)
}

func (r *WarehouseRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Warehouse, error) {
	query := `SELECT id, name, location, description, created_at, updated_at FROM warehouses ORDER BY id DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []entity.Warehouse
	for rows.Next() {
		var w entity.Warehouse
		if err := rows.Scan(&w.ID, &w.Name, &w.Location, &w.Description, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}
	return warehouses, nil
}

func (r *WarehouseRepository) FindByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	query := `SELECT id, name, location, description, created_at, updated_at FROM warehouses WHERE id = $1`
	var w entity.Warehouse
	err := r.db.QueryRow(ctx, query, id).Scan(&w.ID, &w.Name, &w.Location, &w.Description, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *WarehouseRepository) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	query := `UPDATE warehouses SET name = $1, location = $2, description = $3, updated_at = NOW() WHERE id = $4`
	cmd, err := r.db.Exec(ctx, query, warehouse.Name, warehouse.Location, warehouse.Description, warehouse.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *WarehouseRepository) Delete(ctx context.Context, id int64) error {
	// Ingat: Database pakai ON DELETE RESTRICT di produk.
	// Jadi kalau kategori ini dipakai produk, delete ini akan gagal (Error FK).
	query := `DELETE FROM warehouses WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
