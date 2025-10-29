package shop

import (
	"freshease/backend/ent"

	"github.com/gofiber/fiber/v2"
)

// RegisterModuleWithEnt wires Ent repo -> service -> controller and mounts routes.
func RegisterModuleWithEnt(api fiber.Router, client *ent.Client) {
	repo := NewEntRepo(client)
	svc := NewService(repo)
	ctl := NewController(svc)
	Routes(api, ctl)
}
