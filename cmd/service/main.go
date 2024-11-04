package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"strings"
	"xaxaton/internal/handlers/info_user"
	"xaxaton/internal/repo/feedback"
	"xaxaton/internal/usecase/get_info"

	configure "xaxaton/internal/configure"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(
		cors.New(
			cors.Config{
				Next:             nil,
				AllowOriginsFunc: nil,
				AllowOrigins:     "*",
				AllowMethods: strings.Join([]string{
					fiber.MethodGet,
					fiber.MethodPost,
					fiber.MethodHead,
					fiber.MethodPut,
					fiber.MethodDelete,
					fiber.MethodPatch,
				}, ","),
				AllowCredentials: false,
				MaxAge:           0,
			},
		),
	)

	ctx := context.Background()

	cfg := configure.MustConfig(nil)

	dbpool := configure.NewPostgres(ctx, cfg.Postgres)
	defer dbpool.Close()

	// Repo layer
	feedbackDB := feedback.New(dbpool)

	// UseCase layer
	infoUserUC := get_info.NewUseCase(feedbackDB)

	// handlers layer
	infoUserHandler := info_user.NewHandler(infoUserUC)

	if err := cfg.Postgres.MigrationsUp(); err != nil && err.Error() != "no change" {
		panic(err)
	}

	app.Post("/api/1/get/user", infoUserHandler.Handle)
	//app.Post("/api/1/get/users")

	if err := app.Listen(cfg.Fiber.String()); err != nil {
		panic("app not start")
	}
}
