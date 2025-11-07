package vendors

import (
	"freshease/backend/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/",   ctl.ListVendors)
	r.Get("/:id", ctl.GetVendor)
	r.Post("/",  ctl.CreateVendor)
	r.Patch("/:id", ctl.UpdateVendor)
	r.Delete("/:id", ctl.DeleteVendor)
}

// ListVendors godoc
// @Summary      List vendors
// @Description  Get all vendors
// @Tags         vendors
// @Produce      json
// @Success      200 {array}  GetVendorDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /vendors [get]
func (ctl *Controller) ListVendors(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Vendors Retrieved Successfully"})
}

// GetVendor godoc
// @Summary      Get vendor by ID
// @Tags         vendors
// @Produce      json
// @Param        id   path      string true "Vendor ID (UUID)"
// @Success      200  {object}  GetVendorDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /vendors/{id} [get]
func (ctl *Controller) GetVendor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Vendor Retrieved Successfully"})
}

// CreateVendor godoc
// @Summary      Create vendor
// @Tags         vendors
// @Accept       json
// @Produce      json
// @Param        payload body      CreateVendorDTO true "Vendor payload"
// @Success      201     {object}  GetVendorDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /vendors [post]
func (ctl *Controller) CreateVendor(c *fiber.Ctx) error {
	var dto CreateVendorDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Vendor Created Successfully"})
}

// UpdateVendor godoc
// @Summary      Update vendor
// @Tags         vendors
// @Accept       json
// @Produce      json
// @Param        id      path      string         true "Vendor ID (UUID)"
// @Param        payload body      UpdateVendorDTO true "Partial/Full update"
// @Success      201     {object}  GetVendorDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /vendors/{id} [patch]
func (ctl *Controller) UpdateVendor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateVendorDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Vendor Updated Successfully"})
}

// DeleteVendor godoc
// @Summary      Delete vendor
// @Tags         vendors
// @Produce      json
// @Param        id   path      string true "Vendor ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /vendors/{id} [delete]
func (ctl *Controller) DeleteVendor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Vendor Deleted Successfully"})
}
