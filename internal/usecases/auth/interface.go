package auth

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
)

type Repository interface {
	Login(context.Context, *dto.AuthRequest) (*dto.AuthResponse, error)
	Register(context.Context, *dto.AuthRequest) (*dto.AuthResponse, error)
}
