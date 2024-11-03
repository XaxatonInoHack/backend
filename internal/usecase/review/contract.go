package review

import (
	"context"

	"xaxaton/internal/model"
)

type reviewEmployee interface {
	CreateReview(ctx context.Context, reviews []model.Review) error
	GetReview(ctx context.Context, userId int) ([]model.Review, error)
}

type llm interface {
	GetFeedbackLLM(ctx context.Context, selfReview, employeeReview map[int64][]string) (string, string, error)
}

type self interface {
	InsertSelfScore(ctx context.Context, selfReviews []model.SelfReview) error
	UpdateSelfResume(ctx context.Context, selfReviews []model.SelfReview) error
}

type feedback interface {
	CreateFeedback(ctx context.Context, feedbacks []model.Feedback) error
}
