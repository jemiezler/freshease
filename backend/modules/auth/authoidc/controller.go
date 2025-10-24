package authoidc

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

type Controller struct{ s *Service }

func NewController(s *Service) *Controller { return &Controller{s: s} }

func randB64(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// GET /api/auth/:provider/start
func (ctl *Controller) Start(c *fiber.Ctx) error {
	p := ProviderName(c.Params("provider"))
	state := randB64(16)
	nonce := randB64(16)

	c.Cookie(&fiber.Cookie{Name: "oidc_state", Value: state, HTTPOnly: true, SameSite: "Lax", Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "oidc_nonce", Value: nonce, HTTPOnly: true, SameSite: "Lax", Path: "/"})

	url, err := ctl.s.AuthCodeURL(p, state, nonce, "")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// GET /api/auth/:provider/callback?code=...&state=...
func (ctl *Controller) Callback(c *fiber.Ctx) error {
	p := ProviderName(c.Params("provider"))
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || c.Cookies("oidc_state") != state {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid state"})
	}

	access, err := ctl.s.ExchangeAndLogin(c.Context(), p, code, "")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	// clear
	c.Cookie(&fiber.Cookie{Name: "oidc_state", Value: "", MaxAge: -1, Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "oidc_nonce", Value: "", MaxAge: -1, Path: "/"})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    fiber.Map{"accessToken": access},
		"message": "Logged in",
	})
}
