package http

import "github.com/gofiber/fiber/v2"

type Ctx = fiber.Ctx

func New() *fiber.App {
	return fiber.New(fiber.Config{})
}
