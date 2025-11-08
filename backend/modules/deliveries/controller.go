package deliveries

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListDeliveries)
	r.Get("/:id", ctl.GetDelivery)
	r.Post("/", ctl.CreateDelivery)
	r.Patch("/:id", ctl.UpdateDelivery)
	r.Delete("/:id", ctl.DeleteDelivery)
}

// ListDeliveries godoc
// @Summary      List deliveries
// @Description  Get all deliveries
// @Tags         deliveries
// @Produce      json
// @Success      200 {array}  GetDeliveryDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /deliveries [get]
func (ctl *Controller) ListDeliveries(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Deliveries Retrieved Successfully"})
}

// GetDelivery godoc
// @Summary      Get delivery by ID
// @Tags         deliveries
// @Produce      json
// @Param        id   path      string true "Delivery ID (UUID)"
// @Success      200  {object}  GetDeliveryDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /deliveries/{id} [get]
func (ctl *Controller) GetDelivery(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Delivery Retrieved Successfully"})
}

// CreateDelivery godoc
// @Summary      Create delivery
// @Tags         deliveries
// @Accept       json
// @Produce      json
// @Param        payload body      CreateDeliveryDTO true "Delivery data"
// @Success      201     {object}  GetDeliveryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /deliveries [post]
func (ctl *Controller) CreateDelivery(c *fiber.Ctx) error {
	var dto CreateDeliveryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}

	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Delivery Created Successfully"})
}

// UpdateDelivery godoc
// @Summary      Update delivery
// @Tags         deliveries
// @Accept       json
// @Produce      json
// @Param        id      path      string           true "Delivery ID (UUID)"
// @Param        payload body      UpdateDeliveryDTO true "Partial/Full update"
// @Success      201     {object}  GetDeliveryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /deliveries/{id} [patch]
func (ctl *Controller) UpdateDelivery(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateDeliveryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}

	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Delivery Updated Successfully"})
}

// DeleteDelivery godoc
// @Summary      Delete delivery
// @Tags         deliveries
// @Produce      json
// @Param        id   path      string true "Delivery ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /deliveries/{id} [delete]
func (ctl *Controller) DeleteDelivery(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Delivery Deleted Successfully"})
}

