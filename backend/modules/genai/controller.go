package genai

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

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
