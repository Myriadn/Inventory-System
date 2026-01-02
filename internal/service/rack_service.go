package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
)

type RackRepository interface {
	Create(ctx context.Context, r *entity.Rack) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Rack, error)
	FindByID(ctx context.Context, id int64) (*entity.Rack, error)
	Update(ctx context.Context, r *entity.Rack) error
	Delete(ctx context.Context, id int64) error
}

type RackService struct {
	repo RackRepository
}

func NewRackService(repo RackRepository) *RackService {
	return &RackService{repo: repo}
}

func (s *RackService) Create(ctx context.Context, req *entity.RackRequest) (*entity.Rack, error) {
	rack := &entity.Rack{
		Name:     req.Name,
		Category: req.Category,
	}
	err := s.repo.Create(ctx, rack)
	return rack, err
}

func (s *RackService) GetAll(ctx context.Context, page, limit int) ([]entity.Rack, error) {
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

func (s *RackService) GetByID(ctx context.Context, id int64) (*entity.Rack, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("rack not found")
	}
	return cat, nil
}

func (s *RackService) Update(ctx context.Context, id int64, req *entity.RackRequest) error {
	// Cek dulu datanya ada gak
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("rack not found")
	}

	existing.Name = req.Name
	existing.Category = req.Category

	return s.repo.Update(ctx, existing)
}

func (s *RackService) Delete(ctx context.Context, id int64) error {
	// Cek dulu datanya ada gak
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("rack not found")
	}

	// Repository akan melempar error jika kategori sedang dipakai oleh produk (FK Constraint)
	return s.repo.Delete(ctx, id)
}
