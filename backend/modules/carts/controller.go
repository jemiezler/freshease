package carts

import (
	"freshease/backend/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/",   ctl.ListCarts)
	r.Get("/:id", ctl.GetCart)
	r.Post("/",  ctl.CreateCart)
	r.Patch("/:id", ctl.UpdateCart)
	r.Delete("/:id", ctl.DeleteCart)
}

// ListCarts godoc
// @Summary      List carts
// @Description  Get all carts
// @Tags         carts
// @Produce      json
// @Success      200 {array}  GetCartDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /carts [get]
func (ctl *Controller) ListCarts(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Carts Retrieved Successfully"})
}

// GetCart godoc
// @Summary      Get cart by ID
// @Tags         carts
// @Produce      json
// @Param        id   path      string true "Cart ID (UUID)"
// @Success      200  {object}  GetCartDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /carts/{id} [get]
func (ctl *Controller) GetCart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Cart Retrieved Successfully"})
}

// CreateCart godoc
// @Summary      Create cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        payload body      CreateCartDTO true "Cart payload"
// @Success      201     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts [post]
func (ctl *Controller) CreateCart(c *fiber.Ctx) error {
	var dto CreateCartDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Cart Created Successfully"})
}

// UpdateCart godoc
// @Summary      Update cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        id      path      string        true "Cart ID (UUID)"
// @Param        payload body      UpdateCartDTO true "Partial/Full update"
// @Success      201     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts/{id} [patch]
func (ctl *Controller) UpdateCart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateCartDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Cart Updated Successfully"})
}

// DeleteCart godoc
// @Summary      Delete cart
// @Tags         carts
// @Produce      json
// @Param        id   path      string true "Cart ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /carts/{id} [delete]
func (ctl *Controller) DeleteCart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Cart Deleted Successfully"})
}
