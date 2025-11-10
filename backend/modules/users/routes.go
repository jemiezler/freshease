package users

import "github.com/gofiber/fiber/v2"

// Routes keeps routes isolated from wiring; controller methods attach here.
func Routes(app fiber.Router, ctl *Controller) {
	grp := app.Group("/users")
	ctl.Register(grp)
}

// RegisterPublicRoutes registers only public routes (GET endpoints for viewing profiles)
func RegisterPublicRoutes(app fiber.Router, ctl *Controller) {
	grp := app.Group("/users")
	grp.Get("/:id", ctl.GetUser)
}

// RegisterSecuredRoutes registers only secured routes (PUT, DELETE for modifying user data)
func RegisterSecuredRoutes(app fiber.Router, ctl *Controller) {
	grp := app.Group("/users")
	grp.Put("/:id", ctl.UpdateUser)
	grp.Delete("/:id", ctl.DeleteUser)
}
