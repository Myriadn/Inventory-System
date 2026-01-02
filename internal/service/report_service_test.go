package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Report Repository
type MockReportRepo struct {
	mock.Mock
}

func (m *MockReportRepo) GetDashboardStats(ctx context.Context) (*entity.DashboardReport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.DashboardReport), args.Error(1)
}

func TestReportService_GetDashboard(t *testing.T) {
	mockRepo := new(MockReportRepo)
	service := NewReportService(mockRepo)

	expectedReport := &entity.DashboardReport{
		TotalProducts: 10,
		TotalSales:    5,
		TotalRevenue:  500000,
	}

	// Expectation
	mockRepo.On("GetDashboardStats", mock.Anything).Return(expectedReport, nil)

	// Action
	result, err := service.GetDashboard(context.Background())

	// Assertion
	assert.NoError(t, err)
	assert.Equal(t, 10, result.TotalProducts)
	assert.Equal(t, 500000.0, result.TotalRevenue)
}
