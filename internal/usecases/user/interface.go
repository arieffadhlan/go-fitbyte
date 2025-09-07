package user

import (
	"context"
	
	"github.com/arieffadhlan/go-fitbyte/internal/dto"
)

type UseCase interface {
	GetUserById(c context.Context, id int) (*dto.UserResponse, error)
	UpdateUserImage(c context.Context, id int, imageUri string) error
}
