package service

import (
	"context"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/repository"
)

type SaleService struct {
	repo *repository.SaleRepository
}

func NewSaleService(repo *repository.SaleRepository) *SaleService {
	return &SaleService{repo: repo}
}

func (s *SaleService) CreateSale(ctx context.Context, userID int64, req *entity.SaleRequest) (*entity.Sale, error) {
	sale := &entity.Sale{
		UserID: userID,
	}

	// Logic transaksi ada di repository karena melibatkan update multi-tabel atomik
	err := s.repo.CreateTransaction(ctx, sale, req.Items)
	if err != nil {
		return nil, err
	}

	return sale, nil
}
