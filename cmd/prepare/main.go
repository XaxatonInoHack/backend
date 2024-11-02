package main

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	configure "xaxaton/internal/configure"
	"xaxaton/internal/repo/review"
	reviewUC "xaxaton/internal/usecase/review"
)

func main() {
	log.Info("worker prepare is started")

	ctx := context.Background()

	cfg := configure.MustConfig(nil)

	dbpool := configure.NewPostgres(ctx, cfg.Postgres)
	defer dbpool.Close()

	// Repo layer
	//feedbackDB := feedback.New(dbpool)
	reviewDB := review.New(dbpool)

	// UseCase layer
	reviewData := reviewUC.NewUseCase(reviewDB)

	if err := cfg.Postgres.MigrationsUp(); err != nil && err.Error() != "no change" {
		panic(err)
	}

	if err := reviewData.ParseJSON(ctx); err != nil {
		panic(err)
	}
}
