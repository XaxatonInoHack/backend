package info_user

import (
	"context"

	"xaxaton/internal/model"
)

type info interface {
	GetInfo(ctx context.Context, userID int64) (model.Feedback, map[string]float64, map[string]string, error)
}
