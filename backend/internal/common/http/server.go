package http

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Ctx = fiber.Ctx

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "FreshEase API",
	})
	app.Use(middleware.RequestLogger())
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register global middlewares if needed:
	// app.Use(logger.New(), recover.New(), cors.New())

	return app
}
