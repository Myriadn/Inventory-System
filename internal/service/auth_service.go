package service

import (
	"context"
	"errors"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/repository"

	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req *entity.RegisterRequest) error {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Hash Password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Simpan User
	newUser := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPwd),
		Role:         req.Role,
	}

	return s.userRepo.CreateUser(ctx, newUser)
}

func (s *AuthService) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	// Cari User by Email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Cek Password (Compare Hash)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate Session Token (UUID)
	sessionID := uuid.New()
	token := uuid.New()                         // Kita pakai ini sebagai Bearer Token
	expiredAt := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam

	session := &entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: expiredAt,
	}

	// Simpan Session ke DB
	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &entity.LoginResponse{
		Token:     token.String(),
		ExpiredAt: expiredAt.Format(time.RFC3339),
		Role:      user.Role,
	}, nil
}
