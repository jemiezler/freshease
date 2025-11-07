package cart_items

import (
	"freshease/backend/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/",   ctl.ListCart_items)
	r.Get("/:id", ctl.GetCart_item)
	r.Post("/",  ctl.CreateCart_item)
	r.Patch("/:id", ctl.UpdateCart_item)
	r.Delete("/:id", ctl.DeleteCart_item)
}

// ListCart_items godoc
// @Summary      List cart items
// @Description  Get all cart items
// @Tags         cart_items
// @Produce      json
// @Success      200 {array}  GetCart_itemDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /cart_items [get]
func (ctl *Controller) ListCart_items(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Cart_items Retrieved Successfully"})
}

// GetCart_item godoc
// @Summary      Get cart item by ID
// @Tags         cart_items
// @Produce      json
// @Param        id   path      string true "CartItem ID (UUID)"
// @Success      200  {object}  GetCart_itemDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /cart_items/{id} [get]
func (ctl *Controller) GetCart_item(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Cart_item Retrieved Successfully"})
}

// CreateCart_item godoc
// @Summary      Create cart item
// @Tags         cart_items
// @Accept       json
// @Produce      json
// @Param        payload body      CreateCart_itemDTO true "CartItem payload"
// @Success      201     {object}  GetCart_itemDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /cart_items [post]
func (ctl *Controller) CreateCart_item(c *fiber.Ctx) error {
	var dto CreateCart_itemDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Cart_item Created Successfully"})
}

// UpdateCart_item godoc
// @Summary      Update cart item
// @Tags         cart_items
// @Accept       json
// @Produce      json
// @Param        id      path      string              true "CartItem ID (UUID)"
// @Param        payload body      UpdateCart_itemDTO  true "Partial/Full update"
// @Success      201     {object}  GetCart_itemDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /cart_items/{id} [patch]
func (ctl *Controller) UpdateCart_item(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateCart_itemDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Cart_item Updated Successfully"})
}

// DeleteCart_item godoc
// @Summary      Delete cart item
// @Tags         cart_items
// @Produce      json
// @Param        id   path      string true "CartItem ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /cart_items/{id} [delete]
func (ctl *Controller) DeleteCart_item(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Cart_item Deleted Successfully"})
}
