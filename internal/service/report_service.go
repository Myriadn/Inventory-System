package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
)

type ReportRepository interface {
	GetDashboardStats(ctx context.Context) (*entity.DashboardReport, error)
}

type ReportService struct {
	repo ReportRepository
}

func NewReportService(repo ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDashboard(ctx context.Context) (*entity.DashboardReport, error) {
	return s.repo.GetDashboardStats(ctx)
}
