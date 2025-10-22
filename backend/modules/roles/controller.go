package roles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"freshease/backend/internal/common/middleware"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListRoles)
	r.Get("/:id", ctl.GetRole)
	r.Post("/", ctl.CreateRole)
	r.Put("/:id", ctl.UpdateRole)
	r.Delete("/:id", ctl.DeleteRole)
}

func (ctl *Controller) ListRoles(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(items)
}

func (ctl *Controller) GetRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"}) }
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"}) }
	return c.JSON(item)
}

func (ctl *Controller) CreateRole(c *fiber.Ctx) error {
	var dto CreateRoleDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil { return err }
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (ctl *Controller) UpdateRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"}) }
	var dto UpdateRoleDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil { return err }
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(item)
}

func (ctl *Controller) DeleteRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"}) }
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
