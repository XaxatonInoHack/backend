package main

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	configure "xaxaton/internal/configure"
	"xaxaton/internal/gateway/feedback_llm"
	"xaxaton/internal/repo/feedback"
	"xaxaton/internal/repo/review"
	"xaxaton/internal/repo/self_review"
	reviewUC "xaxaton/internal/usecase/review"
)

func main() {
	log.Info("worker prepare is started")

	ctx := context.Background()

	cfg := configure.MustConfig(nil)

	dbpool := configure.NewPostgres(ctx, cfg.Postgres)
	defer dbpool.Close()

	// Repo layer
	feedbackDB := feedback.New(dbpool)
	reviewDB := review.New(dbpool)
	selfDB := self_review.New(dbpool)

	// Gateway layer
	llmGW := feedback_llm.NewGateway()

	// UseCase layer
	reviewData := reviewUC.NewUseCase(reviewDB, llmGW, selfDB, feedbackDB)

	if err := cfg.Postgres.MigrationsUp(); err != nil && err.Error() != "no change" {
		panic(err)
	}

	if err := reviewData.ParseJSON(ctx); err != nil {
		panic(err)
	}
}
