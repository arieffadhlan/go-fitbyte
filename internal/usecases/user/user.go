package user

import (
	"context"
	"fmt"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/repositories/user"
)

type UserUsecase struct {
	repository user.Repository
}

func NewUserUsecase(repository user.Repository) UseCase {
	return &UserUsecase{
		repository: repository,
	}
}

func (u *UserUsecase) GetUserById(c context.Context, id int) (*dto.UserResponse, error) {
	user, err := u.repository.GetUserById(c, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &dto.UserResponse{
		Email:    user.Email,
		ImageURI: user.ImageURI,
	}, nil
}

func (u *UserUsecase) UpdateUserImage(c context.Context, id int, imageUri string) error {
	return u.repository.UpdateUserImage(c, id, imageUri)
}
