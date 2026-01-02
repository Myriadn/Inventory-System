package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"project-app-inventory-restapi-golang-anas/internal/entity"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)

	user := &entity.User{
		Username:     "johndoe",
		Email:        "john@example.com",
		PasswordHash: "hashedsecret",
		Role:         "admin",
	}

	// Expect QueryRow untuk INSERT dengan RETURNING id
	query := regexp.QuoteMeta(`INSERT INTO users (username, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5) RETURNING id`)

	mock.ExpectQuery(query).
		WithArgs(user.Username, user.Email, user.PasswordHash, user.Role, pgxmock.AnyArg()).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(int64(1)))

	err = repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	email := "test@example.com"

	query := regexp.QuoteMeta(`SELECT id, username, email, password_hash, role FROM users WHERE email = $1`)

	rows := mock.NewRows([]string{"id", "username", "email", "password_hash", "role"}).
		AddRow(int64(1), "testuser", email, "hash", "user")

	mock.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)

	user, err := repo.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
}

func TestUserRepository_CreateSession(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	sess := &entity.Session{
		ID: uuid.New(), UserID: 1, Token: uuid.New(), ExpiredAt: time.Now(),
	}

	query := regexp.QuoteMeta(`INSERT INTO sessions (id, user_id, token, expired_at, created_at) VALUES ($1, $2, $3, $4, NOW())`)

	mock.ExpectExec(query).
		WithArgs(sess.ID, sess.UserID, sess.Token, sess.ExpiredAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateSession(context.Background(), sess)
	assert.NoError(t, err)
}

func TestUserRepository_GetSessionByToken(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	token := "valid-token"

	query := regexp.QuoteMeta(`SELECT s.user_id, s.expired_at, u.role FROM sessions s JOIN users u ON s.user_id = u.id WHERE s.token = $1 AND s.is_revoked = FALSE`)

	rows := mock.NewRows([]string{"user_id", "expired_at", "role"}).
		AddRow(int64(1), time.Now().Add(time.Hour), "admin")

	mock.ExpectQuery(query).WithArgs(token).WillReturnRows(rows)

	res, err := repo.GetSessionByToken(context.Background(), token)
	assert.NoError(t, err)
	assert.Equal(t, "admin", res.Role)
}

func TestUserRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, username, email, role, created_at, updated_at FROM users ORDER BY id DESC LIMIT $1 OFFSET $2`)

	rows := mock.NewRows([]string{"id", "username", "email", "role", "created_at", "updated_at"}).
		AddRow(int64(1), "user1", "email1", "admin", nil, nil)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	res, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestUserRepository_FindByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	query := regexp.QuoteMeta(`SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`)

	rows := mock.NewRows([]string{"id", "username", "email", "role", "created_at", "updated_at"}).
		AddRow(int64(1), "user1", "email1", "admin", nil, nil)

	mock.ExpectQuery(query).WithArgs(int64(1)).WillReturnRows(rows)

	res, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "user1", res.Username)
}

func TestUserRepository_UpdateRole(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	query := regexp.QuoteMeta(`UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`)

	mock.ExpectExec(query).
		WithArgs("editor", int64(1)).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateRole(context.Background(), 1, "editor")
	assert.NoError(t, err)
}

func TestUserRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewUserRepository(mock)
	query := regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)

	mock.ExpectExec(query).
		WithArgs(int64(1)).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
