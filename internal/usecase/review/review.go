package review

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/sync/errgroup"

	"xaxaton/internal/model"
)

type UseCase struct {
	feedback feedback
	llm      llm
}

func NewUseCase(f feedback, llm llm) *UseCase {
	return &UseCase{
		feedback: f,
		llm:      llm,
	}
}

type User map[int64][]string

func (u *UseCase) ParseJSON(ctx context.Context) error {
	plan, err := os.ReadFile("internal/usecase/review/review_dataset.json")
	if err != nil {
		return err
	}

	var data []Review

	err = json.Unmarshal(plan, &data)
	if err != nil {
		return err
	}

	g, ctxErr := errgroup.WithContext(ctx)
	g.Go(func() error {
		return u.saveToDB(ctxErr, &data)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) saveToDB(ctx context.Context, reviews *[]Review) error {
	data := make([]model.Review, 0, len(*reviews))

	for _, review := range *reviews {
		data = append(data, model.Review{
			UserID:   review.UserID,
			ReviewID: review.ReviewID,
			Feedback: review.Feedback,
		})
	}

	if err := u.feedback.CreateReview(ctx, data); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) createFeedbackOne(ctx context.Context, reviews *[]Review) error {
	employeeReviews := make(map[int64]User, len(*reviews))
	selfReviews := make(map[int64]User, len(*reviews))

	for _, review := range *reviews {
		if _, ok := employeeReviews[review.UserID]; !ok {
			employeeReviews[review.UserID] = make(User, 100)
		}

		if review.UserID == review.ReviewID {
			selfReviews[review.UserID][review.UserID] = append(selfReviews[review.UserID][review.UserID], review.Feedback)

			continue
		}

		employeeReviews[review.UserID][review.ReviewID] = append(employeeReviews[review.UserID][review.ReviewID], review.Feedback)
	}

	g, errCtx := errgroup.WithContext(ctx)
	for index := range employeeReviews {
		g.Go(func() error {
			employeeReview := employeeReviews[index]

			selfReview, ok := selfReviews[index]
			if !ok {
				selfReview = nil
			}

			return u.llm.GetFeedbackLLM(errCtx, selfReview, employeeReview)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
