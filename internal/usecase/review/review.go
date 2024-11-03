package review

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"xaxaton/internal/lib/parser"

	"golang.org/x/sync/errgroup"

	"xaxaton/internal/model"
)

type UseCase struct {
	review   reviewEmployee
	llm      llm
	self     self
	feedback feedback
}

func NewUseCase(r reviewEmployee, llm llm, s self, f feedback) *UseCase {
	return &UseCase{
		review:   r,
		llm:      llm,
		self:     s,
		feedback: f,
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

	g.Go(func() error {
		return u.createFeedbackOne(ctx, &data)
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

	if err := u.review.CreateReview(ctx, data); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) createFeedbackOne(ctx context.Context, reviews *[]Review) error {
	employeeReviews := make(map[int64]User, len(*reviews))
	selfReviews := make(map[int64]User, len(*reviews))

	for _, review := range *reviews {
		if review.UserID == review.ReviewID {
			if _, ok := selfReviews[review.UserID]; !ok {
				selfReviews[review.UserID] = make(User, 100)
			}

			selfReviews[review.UserID][review.UserID] = append(selfReviews[review.UserID][review.UserID], review.Feedback)

			continue
		}

		if _, ok := employeeReviews[review.UserID]; !ok {
			employeeReviews[review.UserID] = make(User, 100)
		}

		employeeReviews[review.UserID][review.ReviewID] = append(employeeReviews[review.UserID][review.ReviewID], review.Feedback)
	}

	g, errCtx := errgroup.WithContext(ctx)
	for userID := range employeeReviews {
		g.Go(func() error {
			employeeReview := employeeReviews[userID]

			selfReview, ok := selfReviews[userID]
			if !ok {
				selfReview = nil
			}

			employeeFeedback, selfFeedback, err := u.llm.GetFeedbackLLM(errCtx, selfReview, employeeReview)
			if err != nil {
				return err
			}

			if selfFeedback != "" {
				selfScore, _ := parser.ParseReview(selfFeedback)
				err = u.self.InsertSelfScore(ctx, []model.SelfReview{
					{
						UserID: userID,
						Score:  employeeScoreToDB(selfScore),
					},
				})
				if err != nil {
					return fmt.Errorf("insert employee feed score: %w", err)
				}
			}

			employeeScore, _ := parser.ParseReview(employeeFeedback)
			fmt.Println(employeeScoreToDB(employeeScore), employeeFeedback)
			err = u.feedback.CreateFeedback(ctx, []model.Feedback{
				{
					UserID: userID,
					Score:  employeeScoreToDB(employeeScore),
				},
			})

			if err != nil {
				return fmt.Errorf("insert employee feed score: %w", err)
			}

			return nil
		})
		break
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
