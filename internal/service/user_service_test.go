package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserMgmtRepo struct{ mock.Mock }

func (m *MockUserMgmtRepo) FindAll(ctx context.Context, l, o int) ([]entity.User, error) {
	args := m.Called(ctx, l, o)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.User), args.Error(1)
}
func (m *MockUserMgmtRepo) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserMgmtRepo) UpdateRole(ctx context.Context, id int64, r string) error {
	return m.Called(ctx, id, r).Error(0)
}
func (m *MockUserMgmtRepo) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func TestUserService_Full(t *testing.T) {
	mockRepo := new(MockUserMgmtRepo)
	service := NewUserService(mockRepo)

	// GetAll
	mockRepo.On("FindAll", mock.Anything, 10, 0).Return([]entity.User{}, nil)
	_, err := service.GetAll(context.Background(), 1, 10)
	assert.NoError(t, err)

	// UpdateRole (Success)
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.User{}, nil)
	mockRepo.On("UpdateRole", mock.Anything, int64(1), "admin").Return(nil)
	err = service.UpdateRole(context.Background(), 1, "admin")
	assert.NoError(t, err)

	// UpdateRole (Invalid Role)
	err = service.UpdateRole(context.Background(), 1, "invalid_role")
	assert.Error(t, err)

	// Delete (Success)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
	err = service.Delete(context.Background(), 1)
	assert.NoError(t, err)

	// Delete (Not Found)
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	err = service.Delete(context.Background(), 99)
	assert.Error(t, err)
}
