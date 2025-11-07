package permissions

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListPermissions)
	r.Get("/:id", ctl.GetPermission)
	r.Post("/", ctl.CreatePermission)
	r.Patch("/:id", ctl.UpdatePermission)
	r.Delete("/:id", ctl.DeletePermission)
}

// ListPermissions godoc
// @Summary      List permissions
// @Description  Get all permissions
// @Tags         permissions
// @Produce      json
// @Success      200 {array}  GetPermissionDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /permissions [get]
func (ctl *Controller) ListPermissions(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Permissions Retrieved Successfully"})
}

// GetPermission godoc
// @Summary      Get permission by ID
// @Tags         permissions
// @Produce      json
// @Param        id   path      string true "Permission ID (UUID)"
// @Success      200  {object}  GetPermissionDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /permissions/{id} [get]
func (ctl *Controller) GetPermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Permission Retrieved Successfully"})
}

// CreatePermission godoc
// @Summary      Create permission
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        payload body      CreatePermissionDTO true "Permission payload"
// @Success      201     {object}  GetPermissionDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /permissions [post]
func (ctl *Controller) CreatePermission(c *fiber.Ctx) error {
	var dto CreatePermissionDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Permission Created Successfully"})
}

// UpdatePermission godoc
// @Summary      Update permission
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        id      path      string             true "Permission ID (UUID)"
// @Param        payload body      UpdatePermissionDTO true "Partial/Full update"
// @Success      201     {object}  GetPermissionDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /permissions/{id} [patch]
func (ctl *Controller) UpdatePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdatePermissionDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Permission Updated Successfully"})
}

// DeletePermission godoc
// @Summary      Delete permission
// @Tags         permissions
// @Produce      json
// @Param        id   path      string true "Permission ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /permissions/{id} [delete]
func (ctl *Controller) DeletePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Permission Deleted Successfully"})
}
