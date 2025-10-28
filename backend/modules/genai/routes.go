package genai

import "github.com/gofiber/fiber/v2"

func Routes(api fiber.Router, ctl *Controller) {
	r := api.Group("/genai")
	r.Post("/weekly", ctl.GenerateWeeklyMeals)
	r.Post("/daily", ctl.GenerateDailyMeals)

}
