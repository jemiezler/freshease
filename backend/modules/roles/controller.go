package roles

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListRoles)
	r.Get("/:id", ctl.GetRole)
	r.Post("/", ctl.CreateRole)
	r.Patch("/:id", ctl.UpdateRole)
	r.Delete("/:id", ctl.DeleteRole)
}

// ListRoles godoc
// @Summary      List roles
// @Description  Get all roles
// @Tags         roles
// @Produce      json
// @Success      200 {array}  GetRoleDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /roles [get]
func (ctl *Controller) ListRoles(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Roles Retrieved Successfully"})
}

// GetRole godoc
// @Summary      Get role by ID
// @Tags         roles
// @Produce      json
// @Param        id   path      string true "Role ID (UUID)"
// @Success      200  {object}  GetRoleDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /roles/{id} [get]
func (ctl *Controller) GetRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Role Retrieved Successfully"})
}

// CreateRole godoc
// @Summary      Create role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        payload body      CreateRoleDTO true "Role payload"
// @Success      201     {object}  GetRoleDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /roles [post]
func (ctl *Controller) CreateRole(c *fiber.Ctx) error {
	var dto CreateRoleDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Role Created Successfully"})
}

// UpdateRole godoc
// @Summary      Update role
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id      path      string       true "Role ID (UUID)"
// @Param        payload body      UpdateRoleDTO true "Partial/Full update"
// @Success      201     {object}  GetRoleDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /roles/{id} [patch]
func (ctl *Controller) UpdateRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateRoleDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Role Updated Successfully"})
}

// DeleteRole godoc
// @Summary      Delete role
// @Tags         roles
// @Produce      json
// @Param        id   path      string true "Role ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /roles/{id} [delete]
func (ctl *Controller) DeleteRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Role Deleted Successfully"})
}
