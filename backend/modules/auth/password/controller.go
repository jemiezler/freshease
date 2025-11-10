package password

import (
	"freshease/backend/internal/common/middleware"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// Login godoc
// @Summary      Login with email and password
// @Description  Authenticate user with email and password, returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body      LoginRequest true "Login credentials"
// @Success      200     {object}  LoginResponse
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Router       /auth/login [post]
func (ctl *Controller) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	token, user, err := ctl.svc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	roleName := ""
	if user.Edges.Role != nil {
		roleName = user.Edges.Role.Name
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"accessToken": token,
			"user": fiber.Map{
				"id":    user.ID.String(),
				"email": user.Email,
				"name":  user.Name,
				"role":  roleName,
			},
		},
		"message": "Login successful",
	})
}

// InitAdmin godoc
// @Summary      Initialize admin user
// @Description  Create the first admin user with admin role. Can only be called if no admin exists.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body      InitAdminRequest true "Admin user details"
// @Success      201     {object}  InitAdminResponse
// @Failure      400     {object}  map[string]interface{}
// @Failure      409     {object}  map[string]interface{}
// @Router       /auth/init-admin [post]
func (ctl *Controller) InitAdmin(c *fiber.Ctx) error {
	var req InitAdminRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := ctl.svc.InitAdmin(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if err.Error() == "admin user already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	roleName := ""
	if user.Edges.Role != nil {
		roleName = user.Edges.Role.Name
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"user": fiber.Map{
				"id":    user.ID.String(),
				"email": user.Email,
				"name":  user.Name,
				"role":  roleName,
			},
		},
		"message": "Admin user created successfully",
	})
}

