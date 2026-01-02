package repository

import (
	"context"
	"regexp"
	"testing"

	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository_FindByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	expectedID := int64(1)
	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`)

	rows := mock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(expectedID, "Elektronik", "Barang elektronik", nil, nil)

	mock.ExpectQuery(query).
		WithArgs(expectedID).
		WillReturnRows(rows)

	category, err := repo.FindByID(context.Background(), expectedID)

	// Assert hasilnya
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Elektronik", category.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_Create_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	newCategory := &entity.Category{
		Name:        "Fashion",
		Description: "Baju dan Celana",
	}

	query := regexp.QuoteMeta(`INSERT INTO categories (name, description, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`)

	mock.ExpectQuery(query).
		WithArgs(newCategory.Name, newCategory.Description).
		WillReturnRows(mock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(int64(10), nil, nil))

	err = repo.Create(context.Background(), newCategory)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), newCategory.ID) // ID harus terupdate
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`)

	// Simulasi tidak ada baris yang kembali
	mock.ExpectQuery(query).
		WithArgs(int64(999)).
		WillReturnError(pgx.ErrNoRows) // Atau WillReturnRows(pgxmock.NewRows(...)) kosong

	category, err := repo.FindByID(context.Background(), 999)

	assert.NoError(t, err)  // Harusnya tidak error
	assert.Nil(t, category) // Tapi nil
}

func TestCategoryRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, name, description, created_at, updated_at FROM categories ORDER BY id DESC LIMIT $1 OFFSET $2`)

	rows := mock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(int64(1), "Cat1", "Desc1", nil, nil).
		AddRow(int64(2), "Cat2", "Desc2", nil, nil)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	res, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestCategoryRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)
	cat := &entity.Category{ID: 1, Name: "NewName", Description: "NewDesc"}

	query := regexp.QuoteMeta(`UPDATE categories SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`)

	mock.ExpectExec(query).
		WithArgs(cat.Name, cat.Description, cat.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), cat)
	assert.NoError(t, err)
}

func TestCategoryRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)
	query := regexp.QuoteMeta(`DELETE FROM categories WHERE id = $1`)

	mock.ExpectExec(query).
		WithArgs(int64(1)).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
