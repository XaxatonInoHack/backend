package get_info

import (
	"context"

	"xaxaton/internal/model"
)

type feedback interface {
	GetFeedback(ctx context.Context, userId int64) (model.Feedback, error)
}
