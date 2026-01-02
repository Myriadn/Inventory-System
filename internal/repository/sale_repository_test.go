package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestSaleRepository_CreateTransaction_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSaleRepository(mock)

	sale := &entity.Sale{UserID: 1}
	items := []entity.SaleItemRequest{
		{ProductID: 10, Quantity: 2},
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO sales (user_id, transaction_date, total_amount, created_at)`)).
		WithArgs(sale.UserID).
		WillReturnRows(mock.NewRows([]string{"id", "transaction_date"}).AddRow(int64(100), time.Now()))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`)).
		WithArgs(int64(10)).
		WillReturnRows(mock.NewRows([]string{"stock", "price"}).AddRow(50, 10000.0))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(2, int64(10)).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sale_details (sale_id, product_id, quantity, unit_price, subtotal)`)).
		WithArgs(int64(100), int64(10), 2, 10000.0, 20000.0).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE sales SET total_amount = $1 WHERE id = $2`)).
		WithArgs(20000.0, int64(100)).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectCommit()

	err = repo.CreateTransaction(context.Background(), sale, items)

	assert.NoError(t, err)
	assert.Equal(t, 20000.0, sale.TotalAmount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaleRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSaleRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, user_id, transaction_date, total_amount, created_at FROM sales ORDER BY created_at DESC LIMIT $1 OFFSET $2`)

	rows := mock.NewRows([]string{"id", "user_id", "transaction_date", "total_amount", "created_at"}).
		AddRow(int64(1), int64(1), time.Now(), 10000.0, nil)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	res, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestSaleRepository_FindByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSaleRepository(mock)
	saleID := int64(100)

	querySale := regexp.QuoteMeta(`SELECT id, user_id, transaction_date, total_amount, created_at FROM sales WHERE id = $1`)
	rowsSale := mock.NewRows([]string{"id", "user_id", "transaction_date", "total_amount", "created_at"}).
		AddRow(saleID, int64(1), time.Now(), 50000.0, nil)

	mock.ExpectQuery(querySale).WithArgs(saleID).WillReturnRows(rowsSale)

	queryDetails := regexp.QuoteMeta(`SELECT id, sale_id, product_id, quantity, unit_price, subtotal FROM sale_details WHERE sale_id = $1`)
	rowsDetails := mock.NewRows([]string{"id", "sale_id", "product_id", "quantity", "unit_price", "subtotal"}).
		AddRow(int64(1), saleID, int64(10), 2, 25000.0, 50000.0)

	mock.ExpectQuery(queryDetails).WithArgs(saleID).WillReturnRows(rowsDetails)

	res, err := repo.FindByID(context.Background(), saleID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, saleID, res.ID)
	assert.Len(t, res.Details, 1)
	assert.Equal(t, 50000.0, res.Details[0].SubTotal)
}
