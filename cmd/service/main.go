package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"strings"

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

	if err := cfg.Postgres.MigrationsUp(); err != nil && err.Error() != "no change" {
		panic(err)
	}

	if err := app.Listen(cfg.Fiber.String()); err != nil {
		panic("app not start")
	}
}
