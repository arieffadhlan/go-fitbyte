package dto

import (
	"errors"
	"regexp"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (a *AuthRequest) Validate() error {
	if a.Email == "" || a.Password == "" {
		return errors.New("email and password are required")
	}

	if !emailRegex.MatchString(a.Email) {
		return errors.New("invalid email format")
	}

	passLen := len(a.Password)
	if passLen < 8 || passLen > 32 {
		return errors.New("password must be between 8 and 32 characters")
	}
	return nil
}
