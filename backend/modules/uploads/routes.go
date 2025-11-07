package uploads

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(api fiber.Router, ctl *Controller) {
	uploads := api.Group("/uploads")
	ctl.Register(uploads)
}
