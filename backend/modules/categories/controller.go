package categories

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListCategories)
	r.Get("/:id", ctl.GetCategory)
	r.Post("/", ctl.CreateCategory)
	r.Patch("/:id", ctl.UpdateCategory)
	r.Delete("/:id", ctl.DeleteCategory)
}

// ListCategories godoc
// @Summary      List categories
// @Description  Get all categories
// @Tags         categories
// @Produce      json
// @Success      200 {array}  GetCategoryDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /categories [get]
func (ctl *Controller) ListCategories(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Categories Retrieved Successfully"})
}

// GetCategory godoc
// @Summary      Get category by ID
// @Tags         categories
// @Produce      json
// @Param        id   path      string true "Category ID (UUID)"
// @Success      200  {object}  GetCategoryDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /categories/{id} [get]
func (ctl *Controller) GetCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Category Retrieved Successfully"})
}

// CreateCategory godoc
// @Summary      Create category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        payload body      CreateCategoryDTO true "Category data"
// @Success      201     {object}  GetCategoryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /categories [post]
func (ctl *Controller) CreateCategory(c *fiber.Ctx) error {
	var dto CreateCategoryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}

	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Category Created Successfully"})
}

// UpdateCategory godoc
// @Summary      Update category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id      path      string           true "Category ID (UUID)"
// @Param        payload body      UpdateCategoryDTO true "Partial/Full update"
// @Success      201     {object}  GetCategoryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /categories/{id} [patch]
func (ctl *Controller) UpdateCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateCategoryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}

	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Category Updated Successfully"})
}

// DeleteCategory godoc
// @Summary      Delete category
// @Tags         categories
// @Produce      json
// @Param        id   path      string true "Category ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /categories/{id} [delete]
func (ctl *Controller) DeleteCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Category Deleted Successfully"})
}

