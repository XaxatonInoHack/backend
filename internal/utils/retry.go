package utils

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

const (
	colorReset = "\033[0m"

	colorRed = "\033[31m"
)

type Retry struct {
	attempts int64
	sleeps   []time.Duration
	fn       func(context.Context, string) (string, error)
	random   time.Duration
}

func NewRetry(fn func(context.Context, string) (string, error)) *Retry {
	sleeps := []time.Duration{3 * time.Second, 5 * time.Second, 10 * time.Second, 5 * time.Second, 5 * time.Second}
	return &Retry{
		attempts: int64(len(sleeps)),
		sleeps:   sleeps,
		fn:       fn,
		random:   time.Duration(randRange(1, 30)) * time.Second,
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func (r *Retry) Retry(ctx context.Context, prompt string) (string, error) {
	prompt, err := r.fn(ctx, prompt)
	if err != nil {
		if r.attempts == 0 {
			fmt.Println(colorRed, "Retrying after error:", colorReset, err, colorRed, "\nFatal error after retrying")
			return "", err
		}
		fmt.Println(colorRed, "Retrying after error:", colorReset, err, colorRed, "\nAttempt: ", colorReset, int64(len(r.sleeps))-r.attempts+1)
		time.Sleep(r.sleeps[int64(len(r.sleeps))-r.attempts] + r.random)
		r.attempts--
		prompt, err = r.Retry(ctx, prompt)
		if err != nil {
			return "", err
		}
	}

	return prompt, nil
}
