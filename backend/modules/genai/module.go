package genai

import (
	"freshease/backend/ent"

	"github.com/gofiber/fiber/v2"
)

func RegisterModuleWithEnt(api fiber.Router, client *ent.Client) {
	repo := NewEntRepo(client) // or nil if you don't need DB
	svc := NewService(repo)
	ctl := NewController(svc)
	Routes(api, ctl)
}
