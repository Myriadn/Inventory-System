package repository

import (
	"context"
	"regexp"
	"testing"

	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWarehouseRepository(mock)
	wh := &entity.Warehouse{Name: "Gudang A", Location: "Jkt", Description: "Desc"}

	query := regexp.QuoteMeta(`INSERT INTO warehouses (name, location, description, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`)

	mock.ExpectQuery(query).
		WithArgs(wh.Name, wh.Location, wh.Description).
		WillReturnRows(mock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(int64(1), nil, nil))

	err = repo.Create(context.Background(), wh)
	assert.NoError(t, err)
}

func TestWarehouseRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWarehouseRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, name, location, description, created_at, updated_at FROM warehouses ORDER BY id DESC LIMIT $1 OFFSET $2`)

	// Perhatikan urutan kolom MOCK harus sama dengan SELECT di repository (name dulu, baru location)
	rows := mock.NewRows([]string{"id", "name", "location", "description", "created_at", "updated_at"}).
		AddRow(int64(1), "Gudang A", "Jkt", "Desc", nil, nil)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	res, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Gudang A", res[0].Name)
}

func TestWarehouseRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWarehouseRepository(mock)
	wh := &entity.Warehouse{ID: 1, Name: "Baru", Location: "Baru", Description: "Baru"}

	// Pastikan di sini pakai $4 sesuai perbaikan kode repository di atas
	query := regexp.QuoteMeta(`UPDATE warehouses SET name = $1, location = $2, description = $3, updated_at = NOW() WHERE id = $4`)

	mock.ExpectExec(query).
		WithArgs(wh.Name, wh.Location, wh.Description, wh.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), wh)
	assert.NoError(t, err)
}

func TestWarehouseRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewWarehouseRepository(mock)
	query := regexp.QuoteMeta(`DELETE FROM warehouses WHERE id = $1`)

	mock.ExpectExec(query).WithArgs(int64(1)).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
