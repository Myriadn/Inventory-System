package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Product, error)
	FindByID(ctx context.Context, id int64) (*entity.Product, error)
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id int64) error
	FindLowStock(ctx context.Context, threshold int) ([]entity.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req *entity.ProductRequest) (*entity.Product, error) {
	product := &entity.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		RackID:      req.RackID,
		WarehouseID: req.WarehouseID,
	}

	// Repository akan return error jika SKU duplikat atau FK tidak valid
	err := s.repo.Create(ctx, product)
	return product, err
}

func (s *ProductService) GetAll(ctx context.Context, page, limit int) ([]entity.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, id int64, req *entity.ProductRequest) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}

	// Update field
	existing.SKU = req.SKU
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Stock = req.Stock
	existing.Price = req.Price
	existing.CategoryID = req.CategoryID
	existing.RackID = req.RackID
	existing.WarehouseID = req.WarehouseID

	return s.repo.Update(ctx, existing)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	// Cek existance
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("product not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetLowStock(ctx context.Context) ([]entity.Product, error) {
	const minStockThreshold = 5
	return s.repo.FindLowStock(ctx, minStockThreshold)
}
