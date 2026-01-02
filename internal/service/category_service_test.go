package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepo struct{ mock.Mock }

func (m *MockCategoryRepo) Create(ctx context.Context, c *entity.Category) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockCategoryRepo) FindAll(ctx context.Context, l, o int) ([]entity.Category, error) {
	args := m.Called(ctx, l, o)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Category), args.Error(1)
}
func (m *MockCategoryRepo) FindByID(ctx context.Context, id int64) (*entity.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}
func (m *MockCategoryRepo) Update(ctx context.Context, c *entity.Category) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockCategoryRepo) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func TestCategoryService_Create(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := NewCategoryService(mockRepo)
	req := &entity.CategoryRequest{Name: "Cat A"}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	res, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Cat A", res.Name)
}

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := NewCategoryService(mockRepo)

	mockData := []entity.Category{{ID: 1, Name: "Cat A"}}
	mockRepo.On("FindAll", mock.Anything, 10, 0).Return(mockData, nil)

	res, err := service.GetAll(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestCategoryService_GetByID(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := NewCategoryService(mockRepo)

	// Case Found
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Category{ID: 1}, nil)
	res, err := service.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Case Not Found
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	res2, err2 := service.GetByID(context.Background(), 99)
	assert.Error(t, err2)
	assert.Nil(t, res2)
}

func TestCategoryService_Update(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := NewCategoryService(mockRepo)
	req := &entity.CategoryRequest{Name: "Updated"}

	// Found
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Category{ID: 1}, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := service.Update(context.Background(), 1, req)
	assert.NoError(t, err)

	// Not Found
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	err2 := service.Update(context.Background(), 99, req)
	assert.Error(t, err2)
}

func TestCategoryService_Delete(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := NewCategoryService(mockRepo)

	// Found
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Category{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
	err := service.Delete(context.Background(), 1)
	assert.NoError(t, err)

	// Not Found
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	err2 := service.Delete(context.Background(), 99)
	assert.Error(t, err2)
}
