package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRackRepo struct{ mock.Mock }

func (m *MockRackRepo) Create(ctx context.Context, r *entity.Rack) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockRackRepo) FindAll(ctx context.Context, l, o int) ([]entity.Rack, error) {
	args := m.Called(ctx, l, o)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Rack), args.Error(1)
}
func (m *MockRackRepo) FindByID(ctx context.Context, id int64) (*entity.Rack, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Rack), args.Error(1)
}
func (m *MockRackRepo) Update(ctx context.Context, r *entity.Rack) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockRackRepo) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func TestRackService_Full(t *testing.T) {
	mockRepo := new(MockRackRepo)
	service := NewRackService(mockRepo)

	// Create
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	_, err := service.Create(context.Background(), &entity.RackRequest{})
	assert.NoError(t, err)

	// GetAll
	mockRepo.On("FindAll", mock.Anything, 10, 0).Return([]entity.Rack{}, nil)
	_, err = service.GetAll(context.Background(), 1, 10)
	assert.NoError(t, err)

	// GetByID Found
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Rack{}, nil)
	_, err = service.GetByID(context.Background(), 1)
	assert.NoError(t, err)

	// GetByID Not Found
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	_, err = service.GetByID(context.Background(), 99)
	assert.Error(t, err)

	// Update
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	err = service.Update(context.Background(), 1, &entity.RackRequest{})
	assert.NoError(t, err)

	// Delete
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
	err = service.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
