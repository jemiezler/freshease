package product_categories

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListProduct_categories)
	r.Get("/:id", ctl.GetProduct_category)
	r.Post("/", ctl.CreateProduct_category)
	r.Patch("/:id", ctl.UpdateProduct_category)
	r.Delete("/:id", ctl.DeleteProduct_category)
}

// ListProduct_categories godoc
// @Summary      List product categories
// @Description  Get all product categories
// @Tags         product_categories
// @Produce      json
// @Success      200 {array}  GetProductCategoryDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /product_categories [get]
func (ctl *Controller) ListProduct_categories(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Product_categories Retrieved Successfully"})
}

// GetProduct_category godoc
// @Summary      Get product category by ID
// @Tags         product_categories
// @Produce      json
// @Param        id   path      string true "ProductCategory ID (UUID)"
// @Success      200  {object}  GetProductCategoryDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /product_categories/{id} [get]
func (ctl *Controller) GetProduct_category(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Product_category Retrieved Successfully"})
}

// CreateProduct_category godoc
// @Summary      Create product category
// @Tags         product_categories
// @Accept       json
// @Produce      json
// @Param        payload body      CreateProductCategoryDTO true "ProductCategory payload"
// @Success      201     {object}  GetProductCategoryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /product_categories [post]
func (ctl *Controller) CreateProduct_category(c *fiber.Ctx) error {
	var dto CreateProductCategoryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Product_category Created Successfully"})
}

// UpdateProduct_category godoc
// @Summary      Update product category
// @Tags         product_categories
// @Accept       json
// @Produce      json
// @Param        id      path      string                  true "ProductCategory ID (UUID)"
// @Param        payload body      UpdateProductCategoryDTO true "Partial/Full update"
// @Success      201     {object}  GetProductCategoryDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /product_categories/{id} [patch]
func (ctl *Controller) UpdateProduct_category(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateProductCategoryDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Product_category Updated Successfully"})
}

// DeleteProduct_category godoc
// @Summary      Delete product category
// @Tags         product_categories
// @Produce      json
// @Param        id   path      string true "ProductCategory ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /product_categories/{id} [delete]
func (ctl *Controller) DeleteProduct_category(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Product_category Deleted Successfully"})
}
