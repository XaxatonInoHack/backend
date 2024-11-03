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
	GetFeedbackLLMFinal(ctx context.Context, employeeReview map[int64][]string, employeeScore map[int64]model.Feedback) (string, error)
}

type self interface {
	CreateSelfReview(ctx context.Context, selfReviews []model.SelfReview) error
	UpdateSelfResume(ctx context.Context, selfReviews []model.SelfReview) error
	GetSelfReviews(ctx context.Context, userId int64) ([]model.SelfReview, error)
	GetSelfReviewsAll(ctx context.Context, userIDs []int64) (map[int64]model.SelfReview, error)
}

type feedback interface {
	CreateFeedback(ctx context.Context, feedbacks []model.Feedback) error
	GetFeedbackAll(ctx context.Context, userIDs []int64) (map[int64]model.Feedback, error)
	UpdateResume(ctx context.Context, feedbacks []model.Feedback) error
}
