package repository

import (
	"context"
	"regexp"
	"testing"

	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)
	p := &entity.Product{SKU: "S1", Name: "N1", Description: "D1", Stock: 1, Price: 1, CategoryID: 1, RackID: 1, WarehouseID: 1}

	query := regexp.QuoteMeta(`INSERT INTO products (sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) RETURNING id, created_at, updated_at`)

	mock.ExpectQuery(query).
		WithArgs(p.SKU, p.Name, p.Description, p.Stock, p.Price, p.CategoryID, p.RackID, p.WarehouseID).
		WillReturnRows(mock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(int64(1), nil, nil))

	err = repo.Create(context.Background(), p)
	assert.NoError(t, err)
}

func TestProductRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at FROM products ORDER BY id DESC LIMIT $1 OFFSET $2`)

	rows := mock.NewRows([]string{"id", "sku", "name", "description", "stock", "price", "category_id", "rack_id", "warehouse_id", "created_at", "updated_at"}).
		AddRow(int64(1), "S1", "N1", "D1", 1, 100.0, 1, 1, 1, nil, nil)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	res, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestProductRepository_FindByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at FROM products WHERE id = $1`)

	rows := mock.NewRows([]string{"id", "sku", "name", "description", "stock", "price", "category_id", "rack_id", "warehouse_id", "created_at", "updated_at"}).
		AddRow(int64(1), "S1", "N1", "D1", 1, 100.0, 1, 1, 1, nil, nil)

	mock.ExpectQuery(query).WithArgs(int64(1)).WillReturnRows(rows)

	res, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestProductRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)
	p := &entity.Product{ID: 1, SKU: "S1", Name: "N1", Description: "D1", Stock: 1, Price: 1, CategoryID: 1, RackID: 1, WarehouseID: 1}

	query := regexp.QuoteMeta(`UPDATE products SET sku=$1, name=$2, description=$3, stock=$4, price=$5, category_id=$6, rack_id=$7, warehouse_id=$8, updated_at=NOW() WHERE id=$9`)

	mock.ExpectExec(query).
		WithArgs(p.SKU, p.Name, p.Description, p.Stock, p.Price, p.CategoryID, p.RackID, p.WarehouseID, p.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), p)
	assert.NoError(t, err)
}

func TestProductRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)
	query := regexp.QuoteMeta(`DELETE FROM products WHERE id = $1`)

	mock.ExpectExec(query).WithArgs(int64(1)).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
}

func TestProductRepository_FindLowStock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, sku, name, description, stock, price, category_id, rack_id, warehouse_id, created_at, updated_at FROM products WHERE stock < $1 ORDER BY stock ASC`)

	rows := mock.NewRows([]string{"id", "sku", "name", "description", "stock", "price", "category_id", "rack_id", "warehouse_id", "created_at", "updated_at"}).
		AddRow(int64(1), "S1", "N1", "D1", 2, 100.0, 1, 1, 1, nil, nil)

	mock.ExpectQuery(query).WithArgs(5).WillReturnRows(rows)

	res, err := repo.FindLowStock(context.Background(), 5)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}
