package feedback_llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"xaxaton/internal/lib/converter"
	"xaxaton/internal/lib/parser"
	"xaxaton/internal/model"
	"xaxaton/internal/utils"
)

type Gateway struct {
}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (g *Gateway) GetFeedbackLLM(ctx context.Context, selfReview, employeeReview map[int64][]string) (string, string, error) {
	promptEmployee := t1
	promptSelf := t2

	var reviewSelf string

	for employeeID, reviews := range employeeReview {
		promptEmployee += fmt.Sprintf("Review employeeID - %v:\n", employeeID)
		for ind, review := range reviews {
			promptEmployee += fmt.Sprintf("%v review: %v\n\n", ind+1, review)
		}
	}

	promptEmployee += t3
	promptEmployee = template

	reviewEmployee, err := utils.NewRetry(g.getFeedback).Retry(ctx, promptEmployee)
	if err != nil {
		return "", "", err
	}

	if selfReview != nil {
		for selfID, reviews := range selfReview {
			promptSelf += fmt.Sprintf("Review selfID - %v:\n\n", selfID)
			for ind, review := range reviews {
				promptSelf += fmt.Sprintf("%v review: %v\n\n", ind+1, review)
			}
		}

		promptSelf += t4
		promptSelf += template

		reviewSelf, err = utils.NewRetry(g.getFeedback).Retry(ctx, promptSelf)
		if err != nil {
			return "", "", err
		}
	}

	return reviewEmployee, reviewSelf, nil
}

func (g *Gateway) GetFeedbackLLMFinal(ctx context.Context, employeeReview map[int64][]string, employeeScores map[int64]model.Feedback) (string, error) {
	promptEmployee := t1
	for employeeID, reviews := range employeeReview {
		promptEmployee += fmt.Sprintf("Айди ревьюера - %v:\n", employeeID)
		for ind, review := range reviews {
			promptEmployee += fmt.Sprintf("%v отзыв: %v\n\n", ind+1, review)
		}

		employeeScore, ok := employeeScores[employeeID]
		if ok {
			weight := converter.ScoreToResult(parser.ParseScores(employeeScore.Score))
			promptEmployee += fmt.Sprintf("Вес оценивающего сотрудника: %s\n\n", weight)
		}
	}

	promptEmployee += t5
	promptEmployee += template

	reviewEmployee, err := utils.NewRetry(g.getFeedback).Retry(ctx, promptEmployee)
	if err != nil {
		return "", err
	}

	return reviewEmployee, nil
}

func (g *Gateway) getFeedback(ctx context.Context, prompt string) (string, error) {
	URL := "https://vk-scoreworker-case-backup.olymp.innopolis.university/generate"

	data := map[string]interface{}{
		"prompt":              prompt,
		"apply_chat_template": true,
		"system_prompt":       "You are a helpful assistant.",
		"max_tokens":          400,
		"n":                   5,
		"best_of":             10,
		"temperature":         0.7,
	}

	marshalData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(marshalData)

	request, err := http.NewRequestWithContext(ctx, "POST", URL, r)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status: %s", response.Status)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}
