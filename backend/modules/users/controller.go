package users

import (
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
// @Accept       json
// @Produce      json
// @Param        payload body      CreateUserDTO true "User payload"
// @Success      201     {object}  GetUserDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /users [post]
func (ctl *Controller) CreateUser(c *fiber.Ctx) error {
	var dto CreateUserDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
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
// @Accept       json
// @Produce      json
// @Param        id      path      string        true "User ID (UUID)"
// @Param        payload body      UpdateUserDTO true "Partial/Full update"
// @Success      200     {object}  GetUserDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /users/{id} [put]
func (ctl *Controller) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateUserDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
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
// @Router       /users/{id} [delete]
func (ctl *Controller) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
