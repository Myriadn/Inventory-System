package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWarehouseRepo struct{ mock.Mock }

func (m *MockWarehouseRepo) Create(ctx context.Context, w *entity.Warehouse) error {
	return m.Called(ctx, w).Error(0)
}
func (m *MockWarehouseRepo) FindAll(ctx context.Context, l, o int) ([]entity.Warehouse, error) {
	args := m.Called(ctx, l, o)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Warehouse), args.Error(1)
}
func (m *MockWarehouseRepo) FindByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Warehouse), args.Error(1)
}
func (m *MockWarehouseRepo) Update(ctx context.Context, w *entity.Warehouse) error {
	return m.Called(ctx, w).Error(0)
}
func (m *MockWarehouseRepo) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func TestWarehouseService_Full(t *testing.T) {
	mockRepo := new(MockWarehouseRepo)
	service := NewWarehouseService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	_, err := service.Create(context.Background(), &entity.WarehouseRequest{Name: "Gudang"})
	assert.NoError(t, err)

	mockRepo.On("FindAll", mock.Anything, 10, 0).Return([]entity.Warehouse{}, nil)
	_, err = service.GetAll(context.Background(), 1, 10)
	assert.NoError(t, err)

	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Warehouse{}, nil)
	_, err = service.GetByID(context.Background(), 1)
	assert.NoError(t, err)

	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	_, err = service.GetByID(context.Background(), 99)
	assert.Error(t, err)

	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	err = service.Update(context.Background(), 1, &entity.WarehouseRequest{})
	assert.NoError(t, err) // ID 1 sudah di-mock return Warehouse di atas

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
	err = service.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
