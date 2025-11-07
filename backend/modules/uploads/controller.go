package uploads

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	svc Service
}

func NewController(s Service) *Controller {
	return &Controller{svc: s}
}

func (ctl *Controller) Register(r fiber.Router) {
	r.Post("/images", ctl.UploadImage)
	r.Post("/images/:folder", ctl.UploadImageToFolder)
	r.Delete("/images/:path", ctl.DeleteImage)
}

// UploadImage godoc
// @Summary      Upload an image
// @Description  Upload an image file (supports: jpg, jpeg, png, gif, webp). Max size: 10MB
// @Tags         uploads
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Image file to upload"
// @Param        folder formData string false "Folder to store the image (default: 'images')"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /uploads/images [post]
func (ctl *Controller) UploadImage(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "file is required",
			"error":   err.Error(),
		})
	}

	// Get folder from form or use default
	folder := c.FormValue("folder")
	if folder == "" {
		folder = "images"
	}

	// Sanitize folder name
	folder = strings.Trim(folder, "/")
	if folder == "" {
		folder = "images"
	}

	// Upload file
	objectName, err := ctl.svc.UploadImage(c.Context(), file, folder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to upload image",
			"error":   err.Error(),
		})
	}

	// Get URL
	url, err := ctl.svc.GetImageURL(c.Context(), objectName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate image URL",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Image uploaded successfully",
		"object_name": objectName,
		"url":         url,
	})
}

// UploadImageToFolder godoc
// @Summary      Upload an image to a specific folder
// @Description  Upload an image file to a specific folder path
// @Tags         uploads
// @Accept       multipart/form-data
// @Produce      json
// @Param        folder path string true "Folder path (e.g., 'products', 'users/avatars')"
// @Param        file formData file true "Image file to upload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /uploads/images/{folder} [post]
func (ctl *Controller) UploadImageToFolder(c *fiber.Ctx) error {
	folder := c.Params("folder")
	if folder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "folder parameter is required",
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "file is required",
			"error":   err.Error(),
		})
	}

	// Sanitize folder name
	folder = strings.Trim(folder, "/")

	// Upload file
	objectName, err := ctl.svc.UploadImage(c.Context(), file, folder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to upload image",
			"error":   err.Error(),
		})
	}

	// Get URL
	url, err := ctl.svc.GetImageURL(c.Context(), objectName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate image URL",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Image uploaded successfully",
		"object_name": objectName,
		"url":         url,
	})
}

// DeleteImage godoc
// @Summary      Delete an image
// @Description  Delete an image file from storage
// @Tags         uploads
// @Produce      json
// @Param        path path string true "Object path (e.g., 'images/uuid.jpg' or 'products/uuid.png')"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /uploads/images/{path} [delete]
func (ctl *Controller) DeleteImage(c *fiber.Ctx) error {
	path := c.Params("path")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "path parameter is required",
		})
	}

	// Decode URL-encoded path
	path = strings.ReplaceAll(path, "%2F", "/")

	err := ctl.svc.DeleteImage(c.Context(), path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete image",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}
