package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/repository"
)

type ReportService struct {
	repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDashboard(ctx context.Context) (*entity.DashboardReport, error) {
	return s.repo.GetDashboardStats(ctx)
}
