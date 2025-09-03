package activity

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type Repository interface {
	Post(c context.Context, activity *models.Activity) (*models.Activity, error)
	Update(c context.Context, updatedActivity *models.Activity) (*models.Activity, error)
}
