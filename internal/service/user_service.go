package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
)

type UserManagementRepo interface {
	FindAll(ctx context.Context, limit, offset int) ([]entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateRole(ctx context.Context, id int64, newRole string) error
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	repo UserManagementRepo
}

func NewUserService(repo UserManagementRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll(ctx context.Context, page, limit int) ([]entity.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *UserService) UpdateRole(ctx context.Context, id int64, newRole string) error {
	// Validasi Role yang diizinkan
	if newRole != "super_admin" && newRole != "admin" && newRole != "staff" {
		return errors.New("invalid role: must be super_admin, admin, or staff")
	}

	// Cek apakah user target ada
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Update Role
	return s.repo.UpdateRole(ctx, id, newRole)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	// Cek user ada atau tidak
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Hapus User
	return s.repo.Delete(ctx, id)
}
