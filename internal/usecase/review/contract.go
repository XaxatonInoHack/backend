package review

import (
	"context"
	"xaxaton/internal/model"
)

type feedback interface {
	CreateReview(ctx context.Context, reviews []model.Review) error
	GetReview(ctx context.Context, userId int) ([]model.Review, error)
}
