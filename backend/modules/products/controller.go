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
// @Param        payload formData string false "Product payload (JSON string)" example({"name":"Product","price":10.99,"description":"Description","unit_label":"kg","is_active":"active","quantity":100,"restock_amount":50})
// @Success      201     {object}  GetProductDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /products [post]
func (ctl *Controller) CreateProduct(c *fiber.Ctx) error {
	var dto CreateProductDTO

	// Check if this is multipart/form-data (file upload)
	contentType := string(c.Request().Header.ContentType())
	if len(contentType) > 0 && contentType[:19] == "multipart/form-data" {
		// Handle file upload if present
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			// Upload image to MinIO
			_, err = ctl.svc.UploadProductImage(c.Context(), file)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to upload image", "error": err.Error()})
			}
			// Image is uploaded to MinIO, URL is handled by uploads service
		}

		// Parse other form fields - try JSON payload first, then individual fields
		if jsonStr := c.FormValue("payload"); jsonStr != "" {
			if err := c.App().Config().JSONDecoder([]byte(jsonStr), &dto); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload", "error": err.Error()})
			}
		} else if err := c.BodyParser(&dto); err != nil {
			// If BodyParser fails, try individual form fields
			if name := c.FormValue("name"); name != "" {
				dto.Name = name
			}
			// Image URL is handled by uploads service, not stored in product
		}
	} else {
		// Standard JSON request
		if err := middleware.BindAndValidate(c, &dto); err != nil {
			return err
		}
	}

	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
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
// @Param        payload body      UpdateProductDTO true "Partial/Full update"
// @Success      201     {object}  GetProductDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /products/{id} [patch]
func (ctl *Controller) UpdateProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateProductDTO

	// Check if this is multipart/form-data (file upload)
	contentType := string(c.Request().Header.ContentType())
	if len(contentType) > 0 && contentType[:19] == "multipart/form-data" {
		// Handle file upload if present
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			// Upload image to MinIO
			_, err = ctl.svc.UploadProductImage(c.Context(), file)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to upload image", "error": err.Error()})
			}
			// Image is uploaded to MinIO, URL is handled by uploads service
		}

		// Parse other form fields - try JSON payload first, then body parser
		if jsonStr := c.FormValue("payload"); jsonStr != "" {
			if err := c.App().Config().JSONDecoder([]byte(jsonStr), &dto); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid payload", "error": err.Error()})
			}
		} else if err := c.BodyParser(&dto); err != nil {
			// Image URL is handled by uploads service, not stored in product
		}
	} else {
		// Standard JSON request
		if err := middleware.BindAndValidate(c, &dto); err != nil {
			return err
		}
	}

	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Product Updated Successfully"})
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
