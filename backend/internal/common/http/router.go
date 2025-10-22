package http

import (
	"freshease/backend/ent"
	"freshease/backend/modules/roles"
	"freshease/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RegisterRoutes(app *fiber.App, client *ent.Client) {
	api := app.Group("/api")

	log.Debug("[router] registering modules...")

	// --- Core domain modules ---
	users.RegisterModuleWithEnt(api, client)
	roles.RegisterModuleWithEnt(api, client)

	// --- Future modules (just add below) ---
	// products.RegisterModuleWithEnt(api, client)
	// orders.RegisterModuleWithEnt(api, client)
	// auth.RegisterModuleWithEnt(api, client)

	log.Debug("[router] all modules registered ✅")
}
