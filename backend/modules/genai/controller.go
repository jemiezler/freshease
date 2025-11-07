package genai

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

// GenerateWeeklyMeals godoc
// @Summary      Generate a 7-day meal plan
// @Description  Uses user profile and preferences to generate a weekly plan
// @Tags         genai
// @Accept       json
// @Produce      json
// @Param        payload body      GenerateMealsReq true "Generation request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      502     {object}  map[string]interface{}
// @Router       /genai/weekly [post]
func (ctl *Controller) GenerateWeeklyMeals(c *fiber.Ctx) error {
	var req GenerateMealsReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	out, err := ctl.svc.GenerateWeeklyMeals(context.Background(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(out)
}

// GenerateDailyMeals godoc
// @Summary      Generate a 1-day meal plan
// @Description  Uses user profile and preferences to generate a daily plan
// @Tags         genai
// @Accept       json
// @Produce      json
// @Param        payload body      GenerateMealsReq true "Generation request"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      502     {object}  map[string]interface{}
// @Router       /genai/daily [post]
func (ctl *Controller) GenerateDailyMeals(c *fiber.Ctx) error {
	var req GenerateMealsReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	out, err := ctl.svc.GenerateDailyMeals(context.Background(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(out)
}
