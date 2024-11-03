package feedback_llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Gateway struct {
}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (g *Gateway) GetFeedbackLLM(ctx context.Context, selfReview, employeeReview map[int64][]string) (string, string, error) {
	promptEmployee := "Here are some reviews about an employee:\n\n"
	promptSelf := "Here are some reviews about myself:\n\n"

	var reviewSelf string

	for employeeID, reviews := range employeeReview {
		promptEmployee += fmt.Sprintf("Review employeeID - %v:\n", employeeID)
		for ind, review := range reviews {
			promptEmployee += fmt.Sprintf("%v review: %v\n\n", ind+1, review)
		}
	}

	promptEmployee += "Based on these reviews, evaluate the employee on a scale from 1 to 5 for the following criteria:\n"
	promptEmployee += "1. Professionalism\n2. Teamwork\n3. Communication\n4. Initiative\n5. Overall Performance\n"
	promptEmployee += "Add short (5 sentences) explanation for each score you assigned."

	reviewEmployee, err := g.getFeedback(ctx, promptEmployee)
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

		promptEmployee += "Based on these self reviews, evaluate the employee on a scale from 1 to 5 for the following criteria:\n"
		promptEmployee += "1. Professionalism\n2. Teamwork\n3. Communication\n4. Initiative\n5. Overall Performance\n"
		promptEmployee += "Add short (5 sentences) explanation for each score you assigned."

		reviewSelf, err = g.getFeedback(ctx, promptEmployee)
		if err != nil {
			return "", "", err
		}
	}

	return reviewEmployee, reviewSelf, nil
}

func (g *Gateway) getFeedback(ctx context.Context, prompt string) (string, error) {
	URL := "https://vk-scoreworker-case-backup.olymp.innopolis.university/generate"

	fmt.Println(prompt, "\n\n\n")
	data := map[string]interface{}{
		"prompt":              prompt,
		"apply_chat_template": true,
		"system_prompt":       "You are a helpful assistant.",
		"max_tokens":          400,
		"n":                   1,
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
