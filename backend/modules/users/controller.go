package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"freshease/backend/internal/common/middleware"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListUsers)
	r.Get("/:id", ctl.GetUser)
	r.Post("/", ctl.CreateUser)
	r.Put("/:id", ctl.UpdateUser)
	r.Delete("/:id", ctl.DeleteUser)
}

func (ctl *Controller) ListUsers(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(items)
}

func (ctl *Controller) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"}) }
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"}) }
	return c.JSON(item)
}

func (ctl *Controller) CreateUser(c *fiber.Ctx) error {
	var dto CreateUserDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil { return err }
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (ctl *Controller) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"}) }
	var dto UpdateUserDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil { return err }
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(item)
}

func (ctl *Controller) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"}) }
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
