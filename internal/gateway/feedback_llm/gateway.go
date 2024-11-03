package feedback_llm

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type Gateway struct {
}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (g *Gateway) GetFeedbackLLM(ctx context.Context, selfReview, employeeReview map[int64][]string) error {
	promptEmployee := "Here are some reviews about an employee:\n\n"
	promptSelf := "Here are some reviews about myself:\n\n"

	gg, errCtx := errgroup.WithContext(ctx)

	gg.Go(func() error {
		for employeeID, reviews := range employeeReview {
			promptEmployee += fmt.Sprintf("Review employeeID - %v:\n", employeeID)
			for ind, review := range reviews {
				promptEmployee += fmt.Sprintf("%v review: %v\n\n", ind+1, review)
			}
		}

		promptEmployee += "Based on these reviews, evaluate the employee on a scale from 1 to 5 for the following criteria:\n"
		promptEmployee += "1. Professionalism\n2. Teamwork\n3. Communication\n4. Initiative\n5. Overall Performance\n"
		promptEmployee += "Add short (5 sentences) explanation for each score you assigned."

		return g.getEmployeeFeedback(errCtx, promptEmployee)
	})

	if selfReview != nil {
		gg.Go(func() error {
			for selfID, reviews := range selfReview {
				promptSelf += fmt.Sprintf("Review selfID - %v:\n\n", selfID)
				for ind, review := range reviews {
					promptSelf += fmt.Sprintf("%v review: %v\n\n", ind+1, review)
				}
			}

			promptEmployee += "Based on these self reviews, evaluate the employee on a scale from 1 to 5 for the following criteria:\n"
			promptEmployee += "1. Professionalism\n2. Teamwork\n3. Communication\n4. Initiative\n5. Overall Performance\n"
			promptEmployee += "Add short (5 sentences) explanation for each score you assigned."

			return g.getSelfFeedback(errCtx, promptEmployee)
		})
	}

	if err := gg.Wait(); err != nil {
		return err
	}

	return nil
}

func (g *Gateway) getEmployeeFeedback(ctx context.Context, prompt string) error {
	data := map[string]interface{}{
		"prompt":              prompt,
		"apply_chat_template": true,
		"system_prompt":       "You are a helpful assistant.",
		"max_tokens":          400,
		"n":                   1,
		"temperature":         0.7,
	}

}

func (g *Gateway) getSelfFeedback(ctx context.Context, prompt string) error {
	data := map[string]interface{}{
		"prompt":              prompt,
		"apply_chat_template": true,
		"system_prompt":       "You are a helpful assistant.",
		"max_tokens":          400,
		"n":                   1,
		"temperature":         0.7,
	}
}
