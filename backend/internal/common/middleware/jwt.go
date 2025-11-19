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

		// Store user info in context for use in handlers
		if claims, ok := tok.Claims.(jwt.MapClaims); ok {
			// Extract user_id from "sub" claim
			if sub, exists := claims["sub"]; exists {
				if userID, ok := sub.(string); ok && userID != "" {
					c.Locals("user_id", userID)
				}
			}
			// Extract user_email from "email" claim
			if email, exists := claims["email"]; exists {
				if userEmail, ok := email.(string); ok && userEmail != "" {
					c.Locals("user_email", userEmail)
				}
			}
		} else {
			// Try to extract from RegisteredClaims if MapClaims fails
			if regClaims, ok := tok.Claims.(*jwt.RegisteredClaims); ok && regClaims.Subject != "" {
				c.Locals("user_id", regClaims.Subject)
			}
		}

		return c.Next()
	}
}
