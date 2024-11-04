package get_info

import (
	"context"
	"xaxaton/internal/lib/parser"
	"xaxaton/internal/model"
)

type UseCase struct {
	feedback feedback
}

func NewUseCase(f feedback) *UseCase {
	return &UseCase{
		feedback: f,
	}
}

func (uc *UseCase) GetInfo(ctx context.Context, userID int64) (model.Feedback, map[string]float64, map[string]string, error) {
	userFeedback, err := uc.feedback.GetFeedback(ctx, userID)
	if err != nil {
		return model.Feedback{}, nil, nil, err
	}

	score := parser.ParseScores(userFeedback.Score)
	feed := parser.ParseCriteriaText(userFeedback.Resume)

	return userFeedback, score, feed, err
}
