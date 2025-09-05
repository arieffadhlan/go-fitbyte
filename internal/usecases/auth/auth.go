package auth

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/jwt"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/security"
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
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
		return 0, exceptions.NewInternal("failed to hash password")
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	id,err := uc.authRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, req *dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	usr, err := uc.authRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		 return nil, exceptions.NewNotFound("user not found")
	}

	if usr == nil {
		 return nil, exceptions.NewNotFound("invalid email or password")
	}

	if err := security.ComparePassword(req.Password, usr.Password); err != nil {
		 return nil, exceptions.NewNotFound("invalid email or password")
	}

	token, err := jwt.GenerateToken(usr.ID, usr.Email, uc.cfg.JwtSecret)
	if err != nil {
		 return nil, exceptions.NewInternal("failed to generate tokens")
	}

	return &dto.AuthResponse{
		Email: usr.Email,
		Token: token,
	}, nil
}
