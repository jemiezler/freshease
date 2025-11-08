package notifications

import "github.com/gofiber/fiber/v2"

// Routes keeps routes isolated from wiring; controller methods attach here.
func Routes(app fiber.Router, ctl *Controller) {
	grp := app.Group("/notifications")
	ctl.Register(grp)
}
