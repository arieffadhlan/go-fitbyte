package activity

import (
	"context"

	"github.com/arieffadhlan/go-fitbyte/internal/dto"
)

type UseCase interface {
	PostActivity(context.Context, *dto.ActivityRequest, int) (*dto.ActivityResponse, error)
	UpdateActivity(ctx context.Context, updatedActivity *dto.ActivityRequest, userID int, activityID string) (*dto.ActivityResponse, error)
}
