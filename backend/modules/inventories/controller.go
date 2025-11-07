package inventories

import (
	"freshease/backend/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/",   ctl.ListInventories)
	r.Get("/:id", ctl.GetInventory)
	r.Post("/",  ctl.CreateInventory)
	r.Patch("/:id", ctl.UpdateInventory)
	r.Delete("/:id", ctl.DeleteInventory)
}

// ListInventories godoc
// @Summary      List inventories
// @Description  Get all inventories
// @Tags         inventories
// @Produce      json
// @Success      200 {array}  GetInventoryDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /inventories [get]
func (ctl *Controller) ListInventories(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Inventories Retrieved Successfully"})
}

// GetInventory godoc
// @Summary      Get inventory by ID
// @Tags         inventories
// @Produce      json
// @Param        id   path      string true "Inventory ID (UUID)"
// @Success      200  {object}  GetInventoryDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /inventories/{id} [get]
func (ctl *Controller) GetInventory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Inventory Retrieved Successfully"})
}

// CreateInventory godoc
// @Summary      Create inventory
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        payload body      CreateInventoryDTO true "Inventory payload"
// @Success      201     {object}  GetInventoryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /inventories [post]
func (ctl *Controller) CreateInventory(c *fiber.Ctx) error {
	var dto CreateInventoryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Inventory Created Successfully"})
}

// UpdateInventory godoc
// @Summary      Update inventory
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id      path      string            true "Inventory ID (UUID)"
// @Param        payload body      UpdateInventoryDTO true "Partial/Full update"
// @Success      201     {object}  GetInventoryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /inventories/{id} [patch]
func (ctl *Controller) UpdateInventory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateInventoryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Inventory Updated Successfully"})
}

// DeleteInventory godoc
// @Summary      Delete inventory
// @Tags         inventories
// @Produce      json
// @Param        id   path      string true "Inventory ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /inventories/{id} [delete]
func (ctl *Controller) DeleteInventory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Inventory Deleted Successfully"})
}
