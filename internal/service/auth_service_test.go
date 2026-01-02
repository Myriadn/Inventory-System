package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

// Test Cases
func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &entity.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "staff",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(nil, nil)

	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &entity.RegisterRequest{
		Email: "duplicate@example.com",
	}

	existingUser := &entity.User{ID: 1, Email: "duplicate@example.com"}

	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &entity.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hashed),
		Role:         "admin",
	}

	req := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(existingUser, nil)
	mockRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("*entity.Session")).Return(nil)

	resp, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "admin", resp.Role)
	assert.NotEmpty(t, resp.Token)
}
