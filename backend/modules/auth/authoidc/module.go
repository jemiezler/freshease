package authoidc

import (
	"context"

	"freshease/backend/ent"

	"github.com/gofiber/fiber/v2"
)

func RegisterModule(api fiber.Router, db *ent.Client) error {
	svc, err := NewService(context.Background(), db)
	if err != nil {
		return err
	}
	ctl := NewController(svc)

	r := api.Group("/auth")
	r.Get("/:provider/start", ctl.Start)
	r.Get("/:provider/callback", ctl.Callback)
	return nil
}
