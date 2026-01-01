package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, req *entity.CategoryRequest) (*entity.Category, error) {
	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.repo.Create(ctx, category)
	return category, err
}

func (s *CategoryService) GetAll(ctx context.Context, page, limit int) ([]entity.Category, error) {
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

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (s *CategoryService) Update(ctx context.Context, id int64, req *entity.CategoryRequest) error {
	// Cek dulu datanya ada gak
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}

	existing.Name = req.Name
	existing.Description = req.Description

	return s.repo.Update(ctx, existing)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	// Cek dulu datanya ada gak
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}

	// Repository akan melempar error jika kategori sedang dipakai oleh produk (FK Constraint)
	return s.repo.Delete(ctx, id)
}
