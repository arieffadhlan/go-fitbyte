package auth

import (
	"context"
	"errors"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/jwt"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/security"
	"github.com/arieffadhlan/go-fitbyte/internal/repositories/auth"
)

type AuthUsecase struct {
	authRepository auth.Repository
	cfg            *config.Config
}

func NewAuthUsecase(authRepository auth.Repository, cfg *config.Config) *AuthUsecase {
	return &AuthUsecase{
		authRepository: authRepository,
		cfg:            cfg,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, req *dto.AuthRequest) (int, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	hashedPassword, err := security.HashingPassword(req.Password)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	id, err := uc.authRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, req *dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.authRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := security.ComparePassword(req.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, uc.cfg.JwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		Email: user.Email,
		Token: token,
	}, nil
}
