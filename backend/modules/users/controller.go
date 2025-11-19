package users

import (
	"encoding/json"
	"strings"

	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/", ctl.ListUsers)
	r.Get("/:id", ctl.GetUser)
	r.Post("/", ctl.CreateUser)
	r.Put("/:id", ctl.UpdateUser)
	r.Delete("/:id", ctl.DeleteUser)
}

// ListUsers godoc
// @Summary      List users
// @Description  Get all users
// @Tags         users
// @Produce      json
// @Success      200 {array}  GetUserDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /users [get]
func (ctl *Controller) ListUsers(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Users Retrieved Successfully"})
}

// GetUser godoc
// @Summary      Get user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string true "User ID (UUID)"
// @Success      200  {object}  GetUserDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /users/{id} [get]
func (ctl *Controller) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.JSON(item)
}

// CreateUser godoc
// @Summary      Create user
// @Tags         users
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        avatar formData file false "User avatar image file"
// @Param        cover formData file false "User cover image file"
// @Param        payload formData string false "User payload (JSON string)"
// @Success      201     {object}  GetUserDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /users [post]
func (ctl *Controller) CreateUser(c *fiber.Ctx) error {
	var dto CreateUserDTO

	// Use binding helper to handle both multipart/form-data and JSON
	// Note: BindMultipartForm only handles one file field, so we'll handle avatar and cover separately
	_, err := middleware.BindMultipartForm(c, &dto, "")
	if err != nil {
		// If binding fails and it's not a multipart request, try standard JSON binding
		if err := middleware.BindAndValidate(c, &dto); err != nil {
			return err
		}
	}

	// Handle avatar upload if provided
	avatarFile, _ := c.FormFile("avatar")
	if avatarFile != nil && avatarFile.Size > 0 {
		objectName, err := ctl.svc.UploadUserAvatar(c.Context(), avatarFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to upload avatar",
				"error":   err.Error(),
			})
		}
		dto.Avatar = &objectName
	}

	// Handle cover upload if provided
	coverFile, _ := c.FormFile("cover")
	if coverFile != nil && coverFile.Size > 0 {
		objectName, err := ctl.svc.UploadUserCover(c.Context(), coverFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to upload cover",
				"error":   err.Error(),
			})
		}
		dto.Cover = &objectName
	}

	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// UpdateUser godoc
// @Summary      Update user
// @Tags         users
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        id      path      string        true "User ID (UUID)"
// @Param        avatar formData file false "User avatar image file"
// @Param        cover formData file false "User cover image file"
// @Param        payload formData string false "User payload (JSON string)"
// @Success      200     {object}  GetUserDTO
// @Failure      400     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /users/{id} [put]
func (ctl *Controller) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	// Check authorization: users can only update their own profile
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
	}
	if userID.(string) != id.String() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden: can only update own profile"})
	}

	var dto UpdateUserDTO
	dto.ID = id

	// Normalize goal in request body before validation
	contentType := string(c.Request().Header.ContentType())
	if strings.Contains(contentType, "application/json") {
		// Read and normalize goal in JSON body before validation
		bodyBytes := c.Body()
		if len(bodyBytes) > 0 {
			var rawBody map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &rawBody); err == nil {
				if goal, ok := rawBody["goal"].(string); ok {
					normalized := normalizeGoal(goal)
					rawBody["goal"] = normalized
					// Re-set the body with normalized value
					if newBodyBytes, err := json.Marshal(rawBody); err == nil {
						c.Request().SetBody(newBodyBytes)
					}
				}
			}
		}
	}

	// Use binding helper to handle both multipart/form-data and JSON
	// allowEmptyPayload=true for updates (allows image-only updates)
	_, err = middleware.BindMultipartForm(c, &dto, "", true)
	if err != nil {
		// If binding fails and it's not a multipart request, try standard JSON binding
		if err := middleware.BindAndValidate(c, &dto); err != nil {
			return err
		}
	}

	// Handle avatar upload if provided
	avatarFile, _ := c.FormFile("avatar")
	if avatarFile != nil && avatarFile.Size > 0 {
		objectName, err := ctl.svc.UploadUserAvatar(c.Context(), avatarFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to upload avatar",
				"error":   err.Error(),
			})
		}
		dto.Avatar = &objectName
	}

	// Handle cover upload if provided
	coverFile, _ := c.FormFile("cover")
	if coverFile != nil && coverFile.Size > 0 {
		objectName, err := ctl.svc.UploadUserCover(c.Context(), coverFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "failed to upload cover",
				"error":   err.Error(),
			})
		}
		dto.Cover = &objectName
	}

	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(item)
}

// DeleteUser godoc
// @Summary      Delete user
// @Tags         users
// @Produce      json
// @Param        id   path      string true "User ID (UUID)"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /users/{id} [delete]
func (ctl *Controller) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	// Check authorization: users can only delete their own profile
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
	}
	if userID.(string) != id.String() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden: can only delete own profile"})
	}

	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// normalizeGoal converts common goal variations to the expected format
func normalizeGoal(goal string) string {
	goalLower := strings.ToLower(strings.TrimSpace(goal))

	// Map common variations to expected values
	switch goalLower {
	case "loss weight", "lose weight", "weight loss", "weightloss":
		return "weight_loss"
	case "gain weight", "weight gain", "weightgain":
		return "weight_gain"
	case "maintain", "maintenance", "maintain weight":
		return "maintenance"
	default:
		// If it already matches one of the expected values, return as-is
		if goalLower == "weight_loss" || goalLower == "weight_gain" || goalLower == "maintenance" {
			return goalLower
		}
		// Default to maintenance if unknown
		return "maintenance"
	}
}
