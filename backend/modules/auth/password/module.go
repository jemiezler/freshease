package password

import (
	"freshease/backend/ent"

	"github.com/gofiber/fiber/v2"
)

func RegisterModule(api fiber.Router, db *ent.Client) {
	svc := NewService(db)
	ctl := NewController(svc)

	auth := api.Group("/auth")
	auth.Post("/login", ctl.Login)
	auth.Post("/init-admin", ctl.InitAdmin)
}

