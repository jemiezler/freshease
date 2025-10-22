package router

import (
	"freshease/backend/ent"
	"freshease/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RegisterRoutes(api fiber.Router, client *ent.Client) {
	// Register routes here
	log.Debug("[router] registering modules...")

	users.RegisterModuleWithEnt(api, client)

	log.Debug("[router] registered modules")
}
