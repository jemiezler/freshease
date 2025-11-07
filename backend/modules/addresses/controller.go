package addresses

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListAddresses)
	r.Get("/:id", ctl.GetAddress)
	r.Post("/", ctl.CreateAddress)
	r.Patch("/:id", ctl.UpdateAddress)
	r.Delete("/:id", ctl.DeleteAddress)
}

// ListAddresses godoc
// @Summary      List addresses
// @Description  Get all addresses for current user
// @Tags         addresses
// @Produce      json
// @Success      200 {array}  GetAddressDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /addresses [get]
func (ctl *Controller) ListAddresses(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Addresses Retrieved Successfully"})
}

// GetAddress godoc
// @Summary      Get address by ID
// @Tags         addresses
// @Produce      json
// @Param        id   path      string true "Address ID (UUID)"
// @Success      200  {object}  GetAddressDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /addresses/{id} [get]
func (ctl *Controller) GetAddress(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Address Retrieved Successfully"})
}

// CreateAddress godoc
// @Summary      Create address
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        payload body      CreateAddressDTO true "Address payload"
// @Success      201     {object}  GetAddressDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /addresses [post]
func (ctl *Controller) CreateAddress(c *fiber.Ctx) error {
	var dto CreateAddressDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Address Created Successfully"})
}

// UpdateAddress godoc
// @Summary      Update address
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        id      path      string           true "Address ID (UUID)"
// @Param        payload body      UpdateAddressDTO true "Partial/Full update"
// @Success      201     {object}  GetAddressDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /addresses/{id} [patch]
func (ctl *Controller) UpdateAddress(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateAddressDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Address Updated Successfully"})
}

// DeleteAddress godoc
// @Summary      Delete address
// @Tags         addresses
// @Produce      json
// @Param        id   path      string true "Address ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /addresses/{id} [delete]
func (ctl *Controller) DeleteAddress(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Address Deleted Successfully"})
}
