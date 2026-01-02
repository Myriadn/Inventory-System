package repository

import (
	"context"
	"errors"
	"fmt"
	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5"
)

type SaleRepository struct {
	db DBExecutor
}

func NewSaleRepository(db DBExecutor) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) CreateTransaction(ctx context.Context, sale *entity.Sale, items []entity.SaleItemRequest) error {
	// Mulai Transaksi
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer Rollback: Jika fungsi return error di tengah jalan, otomatis rollback.
	// Jika sukses commit, rollback tidak akan melakukan apa-apa.
	defer tx.Rollback(ctx)

	// Buat Header Sale dulu (Total Amount sementara 0, nanti diupdate)
	querySale := `
		INSERT INTO sales (user_id, transaction_date, total_amount, created_at)
		VALUES ($1, NOW(), 0, NOW())
		RETURNING id, transaction_date
	`
	err = tx.QueryRow(ctx, querySale, sale.UserID).Scan(&sale.ID, &sale.TransactionDate)
	if err != nil {
		return fmt.Errorf("failed to create sale header: %w", err)
	}

	var grandTotal float64

	// Loop setiap barang belanjaan
	for _, item := range items {
		var currentStock int
		var price float64

		// AAmbil Info Produk & Lock Row-nya (FOR UPDATE) biar gak balapan sama user lain
		queryProduct := `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`
		err = tx.QueryRow(ctx, queryProduct, item.ProductID).Scan(&currentStock, &price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("product id %d not found", item.ProductID)
			}
			return err
		}

		// Cek Stok Cukup Gak?
		if currentStock < item.Quantity {
			return fmt.Errorf("insufficient stock for product id %d (stock: %d, req: %d)", item.ProductID, currentStock, item.Quantity)
		}

		// Kurangi Stok
		_, err = tx.Exec(ctx, `UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}

		// Hitung Subtotal
		subTotal := price * float64(item.Quantity)
		grandTotal += subTotal

		// Insert ke Sale Details
		queryDetail := `
			INSERT INTO sale_details (sale_id, product_id, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.Exec(ctx, queryDetail, sale.ID, item.ProductID, item.Quantity, price, subTotal)
		if err != nil {
			return err
		}
	}

	// Update Total Amount di Header Sale
	_, err = tx.Exec(ctx, `UPDATE sales SET total_amount = $1 WHERE id = $2`, grandTotal, sale.ID)
	if err != nil {
		return err
	}

	sale.TotalAmount = grandTotal

	// Commit Transaksi (Simpan Permanen)
	return tx.Commit(ctx)
}

// FindAll mengambil daftar penjualan (Header only)
func (r *SaleRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Sale, error) {
	query := `
		SELECT id, user_id, transaction_date, total_amount, created_at
		FROM sales
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []entity.Sale
	for rows.Next() {
		var s entity.Sale
		if err := rows.Scan(&s.ID, &s.UserID, &s.TransactionDate, &s.TotalAmount, &s.CreatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, nil
}

// FindByID mengambil detail penjualan lengkap dengan item barangnya
func (r *SaleRepository) FindByID(ctx context.Context, id int64) (*entity.Sale, error) {
	// 1. Ambil Header Sale
	querySale := `SELECT id, user_id, transaction_date, total_amount, created_at FROM sales WHERE id = $1`
	var s entity.Sale
	err := r.db.QueryRow(ctx, querySale, id).Scan(&s.ID, &s.UserID, &s.TransactionDate, &s.TotalAmount, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// 2. Ambil Detail Barang (Sale Details)
	queryDetails := `
		SELECT id, sale_id, product_id, quantity, unit_price, subtotal
		FROM sale_details
		WHERE sale_id = $1
	`
	rows, err := r.db.Query(ctx, queryDetails, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d entity.SaleDetail
		if err := rows.Scan(&d.ID, &d.SaleID, &d.ProductID, &d.Quantity, &d.UnitPrice, &d.SubTotal); err != nil {
			return nil, err
		}
		s.Details = append(s.Details, d)
	}

	return &s, nil
}
