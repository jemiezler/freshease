package vendors

import (
	"freshease/backend/ent"
	"freshease/backend/modules/uploads"

	"github.com/gofiber/fiber/v2"
)

// RegisterModuleWithEnt wires Ent repo -> service -> controller and mounts routes.
func RegisterModuleWithEnt(api fiber.Router, client *ent.Client, uploadsSvc uploads.Service) {
	repo := NewEntRepo(client)
	svc := NewService(repo, uploadsSvc)
	ctl := NewController(svc)
	Routes(api, ctl)
}
