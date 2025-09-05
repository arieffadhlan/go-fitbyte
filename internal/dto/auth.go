package dto

import (
	"github.com/arieffadhlan/go-fitbyte/internal/pkg/exceptions"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (a *AuthRequest) Validate() error {
	if a.Email == "" || a.Password == "" {
		return exceptions.NewBadRequest("email and password are required")
	}

	passLen := len(a.Password)
	if passLen < 8 || passLen > 32 {
		return exceptions.NewBadRequest("password must be between 8 and 32 characters")
	}
	return nil
}
