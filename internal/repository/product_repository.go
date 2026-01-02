package repository

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db DBExecutor
}

func NewProductRepository(db DBExecutor) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *entity.Product) error {
	query := `
		INSERT INTO products (sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		p.SKU, p.Name, p.Description, p.Stock, p.Price,
		p.CategoryID, p.RackID, p.WarehouseID,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Product, error) {
	query := `
		SELECT id, sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at
		FROM products
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.Stock, &p.Price,
			&p.CategoryID, &p.RackID, &p.WarehouseID, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `
		SELECT id, sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	var p entity.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Stock, &p.Price,
		&p.CategoryID, &p.RackID, &p.WarehouseID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *entity.Product) error {
	query := `
		UPDATE products
		SET sku=$1, name=$2, description=$3, stock=$4, price=$5,
		    category_id=$6, rack_id=$7, warehouse_id=$8, updated_at=NOW()
		WHERE id=$9
	`
	cmd, err := r.db.Exec(ctx, query,
		p.SKU, p.Name, p.Description, p.Stock, p.Price,
		p.CategoryID, p.RackID, p.WarehouseID, p.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		// Akan error jika produk sudah pernah terjual (Constraint FK di tabel sale_details)
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) FindLowStock(ctx context.Context, threshold int) ([]entity.Product, error) {
	query := `
		SELECT id, sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at
		FROM products
		WHERE stock < $1
		ORDER BY stock ASC
	`
	rows, err := r.db.Query(ctx, query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.Stock, &p.Price,
			&p.CategoryID, &p.RackID, &p.WarehouseID, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
