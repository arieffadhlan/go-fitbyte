package user

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type Repository interface {
	GetUserById(c context.Context, id int) (*models.User, error)
	UpdateUserImage(c context.Context, id int, imageUri string) error
}
