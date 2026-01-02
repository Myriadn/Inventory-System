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

func TestRackRepository_FindByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewRackRepository(mock)

	query := regexp.QuoteMeta(`SELECT id, name, category, created_at, updated_at FROM racks WHERE id = $1`)

	rows := mock.NewRows([]string{"id", "name", "category", "created_at", "updated_at"}).
		AddRow(int64(5), "Rack A1", "Food", nil, nil)

	mock.ExpectQuery(query).WithArgs(int64(5)).WillReturnRows(rows)

	res, err := repo.FindByID(context.Background(), 5)

	assert.NoError(t, err)
	assert.Equal(t, "Rack A1", res.Name)
}

func TestRackRepository_Delete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewRackRepository(mock)

	query := regexp.QuoteMeta(`DELETE FROM racks WHERE id = $1`)

	mock.ExpectExec(query).
		WithArgs(int64(99)).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(context.Background(), 99)

	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}

func TestRackRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewRackRepository(mock)
	rack := &entity.Rack{Name: "R-01", Category: "General"}

	query := regexp.QuoteMeta(`INSERT INTO racks (name, category, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`)

	mock.ExpectQuery(query).
		WithArgs(rack.Name, rack.Category).
		WillReturnRows(mock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(int64(1), nil, nil))

	err = repo.Create(context.Background(), rack)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rack.ID)
}

func TestRackRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewRackRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, name, category, created_at, updated_at FROM racks ORDER BY id DESC LIMIT $1 OFFSET $2`)

	rows := mock.NewRows([]string{"id", "name", "category", "created_at", "updated_at"}).
		AddRow(int64(1), "R1", "C1", nil, nil)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	res, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestRackRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewRackRepository(mock)
	rack := &entity.Rack{ID: 1, Name: "RackUpdated", Category: "CatUpdated"}

	query := regexp.QuoteMeta(`UPDATE racks SET name = $1, category = $2, updated_at = NOW() WHERE id = $3`)

	mock.ExpectExec(query).
		WithArgs(rack.Name, rack.Category, rack.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), rack)
	assert.NoError(t, err)
}

func TestRackRepository_Delete_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewRackRepository(mock)
	query := regexp.QuoteMeta(`DELETE FROM racks WHERE id = $1`)

	mock.ExpectExec(query).
		WithArgs(int64(1)).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)) // 1 row affected

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
