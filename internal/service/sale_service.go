package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
)

type SaleRepository interface {
	CreateTransaction(ctx context.Context, sale *entity.Sale, items []entity.SaleItemRequest) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Sale, error)
	FindByID(ctx context.Context, id int64) (*entity.Sale, error)
}

type SaleService struct {
	repo SaleRepository
}

func NewSaleService(repo SaleRepository) *SaleService {
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

func (s *SaleService) GetAll(ctx context.Context, page, limit int) ([]entity.Sale, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *SaleService) GetByID(ctx context.Context, id int64) (*entity.Sale, error) {
	sale, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, errors.New("transaction not found")
	}
	return sale, nil
}
