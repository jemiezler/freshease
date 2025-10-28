package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth() fiber.Handler {
	sec := []byte(os.Getenv("JWT_SECRET"))
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "missing bearer token"})
		}
		raw := strings.TrimPrefix(h, "Bearer ")
		tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) { return sec, nil })
		if err != nil || !tok.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
		}
		return c.Next()
	}
}
