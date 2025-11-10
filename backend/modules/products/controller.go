package products

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListProducts)
	r.Get("/:id", ctl.GetProduct)
	r.Post("/", ctl.CreateProduct)
	r.Patch("/:id", ctl.UpdateProduct)
	r.Delete("/:id", ctl.DeleteProduct)
}

// ListProducts godoc
// @Summary      List products
// @Description  Get all products
// @Tags         products
// @Produce      json
// @Success      200 {array}  GetProductDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /products [get]
func (ctl *Controller) ListProducts(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Products Retrieved Successfully"})
}

// GetProduct godoc
// @Summary      Get product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      string true "Product ID (UUID)"
// @Success      200  {object}  GetProductDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func (ctl *Controller) GetProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Product Retrieved Successfully"})
}

// CreateProduct godoc
// @Summary      Create product
// @Tags         products
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        image formData file false "Product image file"
// @Param        payload formData string false "Product payload (JSON string)" example({"name":"Product","price":10.99,"description":"Description","unit_label":"kg","is_active":true,"quantity":100,"reorder_level":50})
// @Success      201     {object}  GetProductDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /products [post]
func (ctl *Controller) CreateProduct(c *fiber.Ctx) error {
	var dto CreateProductDTO

	// Use binding helper to handle both multipart/form-data and JSON
	file, err := middleware.BindMultipartForm(c, &dto, "image")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"message": fiberErr.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Handle image upload if file is provided
	if file != nil {
		// Upload image to MinIO and get object name (path, not URL)
		objectName, err := ctl.svc.UploadProductImage(c.Context(), file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to upload image",
				"error":   err.Error(),
			})
		}
		// Set the object name in DTO (this will be stored in database)
		dto.ImageURL = &objectName
	}

	// Create product (service will handle URL generation when retrieving)
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Product Created Successfully"})
}

// UpdateProduct godoc
// @Summary      Update product
// @Tags         products
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        id      path      string           true "Product ID (UUID)"
// @Param        image formData file false "Product image file"
// @Param        payload formData string false "Product payload (JSON string)" example({"name":"Updated Product","price":15.99})
// @Success      200     {object}  GetProductDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /products/{id} [patch]
func (ctl *Controller) UpdateProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	var dto UpdateProductDTO
	dto.ID = id

	// Use binding helper to handle both multipart/form-data and JSON
	// allowEmptyPayload=true for updates (allows image-only updates)
	file, err := middleware.BindMultipartForm(c, &dto, "image", true)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"message": fiberErr.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Handle image upload if file is provided
	if file != nil {
		// Upload image to MinIO and get object name (path, not URL)
		objectName, err := ctl.svc.UploadProductImage(c.Context(), file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to upload image",
				"error":   err.Error(),
			})
		}
		// Set the object name in DTO (this will be stored in database)
		dto.ImageURL = &objectName
	}

	// Update product (service will handle URL generation when retrieving)
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Product Updated Successfully"})
}

// DeleteProduct godoc
// @Summary      Delete product
// @Tags         products
// @Produce      json
// @Param        id   path      string true "Product ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /products/{id} [delete]
func (ctl *Controller) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Product Deleted Successfully"})
}
