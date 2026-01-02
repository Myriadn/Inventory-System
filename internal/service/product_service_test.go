package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Product Repository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, p *entity.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockProductRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Product, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Product), args.Error(1)
}
func (m *MockProductRepository) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}
func (m *MockProductRepository) Update(ctx context.Context, p *entity.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockProductRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockProductRepository) FindLowStock(ctx context.Context, threshold int) ([]entity.Product, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Product), args.Error(1)
}

func TestProductService_Create(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	req := &entity.ProductRequest{
		SKU: "TEST-001", Name: "Product Test", Price: 1000, Stock: 10,
		CategoryID: 1, RackID: 1, WarehouseID: 1,
	}

	// Expectation: Repository Create dipanggil tanpa error
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil)

	res, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "TEST-001", res.SKU)
}

func TestProductService_GetByID_Found(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockProduct := &entity.Product{ID: 1, Name: "Product Test"}

	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(mockProduct, nil)

	res, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Product Test", res.Name)
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)

	res, err := service.GetByID(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "product not found", err.Error())
}

func TestProductService_GetAll(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindAll", mock.Anything, 10, 0).Return([]entity.Product{}, nil)
	_, err := service.GetAll(context.Background(), 1, 10)
	assert.NoError(t, err)
}

func TestProductService_Update(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	// Found & Update
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Product{}, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	err := service.Update(context.Background(), 1, &entity.ProductRequest{Price: 1000})
	assert.NoError(t, err)

	// Not Found
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	err = service.Update(context.Background(), 99, &entity.ProductRequest{})
	assert.Error(t, err)
}

func TestProductService_Delete(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	// Found & Delete
	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&entity.Product{}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)
	err := service.Delete(context.Background(), 1)
	assert.NoError(t, err)

	// Not Found
	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
	err = service.Delete(context.Background(), 99)
	assert.Error(t, err)
}

func TestProductService_GetLowStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindLowStock", mock.Anything, 5).Return([]entity.Product{}, nil)
	_, err := service.GetLowStock(context.Background())
	assert.NoError(t, err)
}
