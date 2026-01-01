package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/repository"
)

type WarehouseService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseService(repo *repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) Create(ctx context.Context, req *entity.WarehouseRequest) (*entity.Warehouse, error) {
	warehouse := &entity.Warehouse{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}
	err := s.repo.Create(ctx, warehouse)
	return warehouse, err
}

func (s *WarehouseService) GetAll(ctx context.Context, page, limit int) ([]entity.Warehouse, error) {
	// Simple Pagination Logic
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	return s.repo.FindAll(ctx, limit, offset)
}

func (s *WarehouseService) GetByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("warehouses not found")
	}
	return cat, nil
}

func (s *WarehouseService) Update(ctx context.Context, id int64, req *entity.WarehouseRequest) error {
	// Cek dulu datanya ada gak
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("warehouse not found")
	}

	existing.Name = req.Name
	existing.Location = req.Location
	existing.Description = req.Description

	return s.repo.Update(ctx, existing)
}

func (s *WarehouseService) Delete(ctx context.Context, id int64) error {
	// Cek dulu datanya ada gak
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("warehouse not found")
	}

	// Repository akan melempar error jika kategori sedang dipakai oleh produk (FK Constraint)
	return s.repo.Delete(ctx, id)
}
