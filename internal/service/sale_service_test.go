package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Sale Repository
type MockSaleRepo struct {
	mock.Mock
}

func (m *MockSaleRepo) CreateTransaction(ctx context.Context, sale *entity.Sale, items []entity.SaleItemRequest) error {
	args := m.Called(ctx, sale, items)
	return args.Error(0)
}

func (m *MockSaleRepo) FindAll(ctx context.Context, limit, offset int) ([]entity.Sale, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Sale), args.Error(1)
}

func (m *MockSaleRepo) FindByID(ctx context.Context, id int64) (*entity.Sale, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Sale), args.Error(1)
}

func TestSaleService_CreateSale(t *testing.T) {
	mockRepo := new(MockSaleRepo)
	service := NewSaleService(mockRepo)

	userID := int64(1)
	req := &entity.SaleRequest{
		Items: []entity.SaleItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}

	// Expectation
	mockRepo.On("CreateTransaction", mock.Anything, mock.Anything, req.Items).Return(nil)

	// Action
	result, err := service.CreateSale(context.Background(), userID, req)

	// Assertion
	assert.NoError(t, err)
	assert.Equal(t, userID, result.UserID)
	mockRepo.AssertExpectations(t)
}

func TestSaleService_GetAll(t *testing.T) {
	mockRepo := new(MockSaleRepo)
	service := NewSaleService(mockRepo)

	mockData := []entity.Sale{
		{ID: 1, TotalAmount: 10000},
	}

	mockRepo.On("FindAll", mock.Anything, 10, 0).Return(mockData, nil)

	result, err := service.GetAll(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, float64(10000), result[0].TotalAmount)
}

func TestSaleService_GetByID_Found(t *testing.T) {
	mockRepo := new(MockSaleRepo)
	service := NewSaleService(mockRepo)

	mockSale := &entity.Sale{ID: 1, TotalAmount: 50000}

	mockRepo.On("FindByID", mock.Anything, int64(1)).Return(mockSale, nil)

	result, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, float64(50000), result.TotalAmount)
}

func TestSaleService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockSaleRepo)
	service := NewSaleService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)

	result, err := service.GetByID(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "transaction not found", err.Error())
}
