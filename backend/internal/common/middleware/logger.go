package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		c.Set("X-Response-Time", time.Since(start).String())

		return err
	}
}
