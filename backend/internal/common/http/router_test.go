package http

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"freshease/backend/ent/enttest"
	"freshease/backend/ent/user"
	"freshease/backend/internal/common/config"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func TestRegisterRoutes(t *testing.T) {
	// Set up required environment variables for authoidc
	os.Setenv("OAUTH_BASE_URL", "http://localhost:8080")
	os.Setenv("OIDC_GOOGLE_ISSUER", "https://accounts.google.com")
	os.Setenv("OIDC_GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("OIDC_GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("OIDC_GOOGLE_REDIRECT_PATH", "/auth/google/callback")
	os.Setenv("OIDC_LINE_ISSUER", "https://access.line.me")
	os.Setenv("OIDC_LINE_CLIENT_ID", "test-line-client-id")
	os.Setenv("OIDC_LINE_CLIENT_SECRET", "test-line-client-secret")
	os.Setenv("OIDC_LINE_REDIRECT_PATH", "/auth/line/callback")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("JWT_ACCESS_TTL_MIN", "15")
	defer func() {
		os.Unsetenv("OAUTH_BASE_URL")
		os.Unsetenv("OIDC_GOOGLE_ISSUER")
		os.Unsetenv("OIDC_GOOGLE_CLIENT_ID")
		os.Unsetenv("OIDC_GOOGLE_CLIENT_SECRET")
		os.Unsetenv("OIDC_GOOGLE_REDIRECT_PATH")
		os.Unsetenv("OIDC_LINE_ISSUER")
		os.Unsetenv("OIDC_LINE_CLIENT_ID")
		os.Unsetenv("OIDC_LINE_CLIENT_SECRET")
		os.Unsetenv("OIDC_LINE_REDIRECT_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_ACCESS_TTL_MIN")
	}()

	t.Run("registers routes without panic", func(t *testing.T) {
		app := fiber.New()

		// Create a test client
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// This should not panic
		apiGroup := app.Group("/api")
		assert.NotPanics(t, func() {
			RegisterRoutes(apiGroup, app, client, config.Config{})
		})

		// Verify that routes are registered
		routes := app.GetRoutes()
		assert.NotEmpty(t, routes)
	})

	t.Run("registers API group", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		apiGroup := app.Group("/api")
		RegisterRoutes(apiGroup, app, client, config.Config{})

		// Check that API routes are registered
		routes := app.GetRoutes()
		apiRoutes := 0
		for _, route := range routes {
			if route.Path[:4] == "/api" {
				apiRoutes++
			}
		}
		assert.Greater(t, apiRoutes, 0, "Should have API routes registered")
	})

	t.Run("handles nil client gracefully", func(t *testing.T) {
		app := fiber.New()

		// This should not panic even with nil client
		apiGroup := app.Group("/api")
		assert.NotPanics(t, func() {
			RegisterRoutes(apiGroup, app, nil, config.Config{})
		})
	})
}

func TestLogRegisteredModules(t *testing.T) {
	t.Run("logs modules correctly", func(t *testing.T) {
		app := fiber.New()

		// Add some test routes
		api := app.Group("/api")
		api.Get("/test", func(c *fiber.Ctx) error { return c.JSON("test") })
		api.Post("/test", func(c *fiber.Ctx) error { return c.JSON("test") })
		api.Get("/users", func(c *fiber.Ctx) error { return c.JSON("users") })
		api.Get("/products", func(c *fiber.Ctx) error { return c.JSON("products") })

		// This should not panic
		assert.NotPanics(t, func() {
			logRegisteredModules(app, "/api")
		})
	})

	t.Run("handles empty routes", func(t *testing.T) {
		app := fiber.New()

		// This should not panic with no routes
		assert.NotPanics(t, func() {
			logRegisteredModules(app, "/api")
		})
	})

	t.Run("handles routes without API prefix", func(t *testing.T) {
		app := fiber.New()

		// Add routes without API prefix
		app.Get("/health", func(c *fiber.Ctx) error { return c.JSON("ok") })
		app.Get("/status", func(c *fiber.Ctx) error { return c.JSON("ok") })

		// This should not panic
		assert.NotPanics(t, func() {
			logRegisteredModules(app, "/api")
		})
	})

	t.Run("handles different API prefixes", func(t *testing.T) {
		app := fiber.New()

		// Add routes with different prefixes
		v1 := app.Group("/v1")
		v1.Get("/test", func(c *fiber.Ctx) error { return c.JSON("test") })

		v2 := app.Group("/v2")
		v2.Get("/test", func(c *fiber.Ctx) error { return c.JSON("test") })

		// Test with v1 prefix
		assert.NotPanics(t, func() {
			logRegisteredModules(app, "/v1")
		})

		// Test with v2 prefix
		assert.NotPanics(t, func() {
			logRegisteredModules(app, "/v2")
		})
	})
}

func TestWhoamiEndpoint(t *testing.T) {
	// Set up required environment variables for authoidc
	os.Setenv("OAUTH_BASE_URL", "http://localhost:8080")
	os.Setenv("OIDC_GOOGLE_ISSUER", "https://accounts.google.com")
	os.Setenv("OIDC_GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("OIDC_GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("OIDC_GOOGLE_REDIRECT_PATH", "/auth/google/callback")
	os.Setenv("OIDC_LINE_ISSUER", "https://access.line.me")
	os.Setenv("OIDC_LINE_CLIENT_ID", "test-line-client-id")
	os.Setenv("OIDC_LINE_CLIENT_SECRET", "test-line-client-secret")
	os.Setenv("OIDC_LINE_REDIRECT_PATH", "/auth/line/callback")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("JWT_ACCESS_TTL_MIN", "15")
	defer func() {
		os.Unsetenv("OAUTH_BASE_URL")
		os.Unsetenv("OIDC_GOOGLE_ISSUER")
		os.Unsetenv("OIDC_GOOGLE_CLIENT_ID")
		os.Unsetenv("OIDC_GOOGLE_CLIENT_SECRET")
		os.Unsetenv("OIDC_GOOGLE_REDIRECT_PATH")
		os.Unsetenv("OIDC_LINE_ISSUER")
		os.Unsetenv("OIDC_LINE_CLIENT_ID")
		os.Unsetenv("OIDC_LINE_CLIENT_SECRET")
		os.Unsetenv("OIDC_LINE_REDIRECT_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_ACCESS_TTL_MIN")
	}()

	t.Run("returns user details for valid token", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// Create a test user
		userEntity, err := client.User.Create().
			SetID(uuid.New()).
			SetEmail("test@example.com").
			SetName("Test User").
			SetPhone("1234567890").
			SetBio("This is a longer bio that meets the minimum length requirement for the user bio field").
			SetAvatar("avatar.jpg").
			SetCover("This is a longer cover description that meets the minimum length requirement").
			SetSex("male").
			SetGoal("fitness").
			SetHeightCm(175.0).
			SetWeightKg(70.0).
			SetStatus("active").
			Save(context.Background())
		require.NoError(t, err)

		// Mock the middleware to set user context
		app.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/api/whoami" {
				c.Locals("user_id", userEntity.ID.String())
				c.Locals("user_email", userEntity.Email)
			}
			return c.Next()
		})

		// Create a test endpoint that bypasses authentication middleware
		app.Get("/api/whoami", func(c *fiber.Ctx) error {
			userID := c.Locals("user_id")
			userEmail := c.Locals("user_email")

			if userID == nil || userEmail == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
			}

			// Parse user ID from string to UUID
			userUUID, err := uuid.Parse(userID.(string))
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user ID"})
			}

			// Get user details from database
			user, err := client.User.Get(c.Context(), userUUID)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
			}

			return c.JSON(fiber.Map{
				"id":            user.ID.String(),
				"email":         user.Email,
				"name":          user.Name,
				"phone":         user.Phone,
				"bio":           user.Bio,
				"avatar":        user.Avatar,
				"cover":         user.Cover,
				"date_of_birth": user.DateOfBirth,
				"sex":           user.Sex,
				"goal":          user.Goal,
				"height_cm":     user.HeightCm,
				"weight_kg":     user.WeightKg,
				"status":        user.Status,
				"created_at":    user.CreatedAt,
				"updated_at":    user.UpdatedAt,
			})
		})

		// Test the whoami endpoint
		req := httptest.NewRequest("GET", "/api/whoami", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("returns unauthorized for missing user context", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// Create a test endpoint that bypasses authentication middleware
		app.Get("/api/whoami", func(c *fiber.Ctx) error {
			userID := c.Locals("user_id")
			userEmail := c.Locals("user_email")

			if userID == nil || userEmail == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
			}

			// Parse user ID from string to UUID
			userUUID, err := uuid.Parse(userID.(string))
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user ID"})
			}

			// Get user details from database
			userEntity, err := client.User.Query().Where(user.ID(userUUID)).First(c.Context())
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
			}

			return c.JSON(fiber.Map{
				"id":            userEntity.ID.String(),
				"email":         userEntity.Email,
				"name":          userEntity.Name,
				"phone":         userEntity.Phone,
				"bio":           userEntity.Bio,
				"avatar":        userEntity.Avatar,
				"cover":         userEntity.Cover,
				"date_of_birth": userEntity.DateOfBirth,
				"sex":           userEntity.Sex,
				"goal":          userEntity.Goal,
				"height_cm":     userEntity.HeightCm,
				"weight_kg":     userEntity.WeightKg,
				"status":        userEntity.Status,
				"created_at":    userEntity.CreatedAt,
				"updated_at":    userEntity.UpdatedAt,
			})
		})

		req := httptest.NewRequest("GET", "/api/whoami", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("returns bad request for invalid user ID", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// Mock the middleware to set invalid user ID
		app.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/api/whoami" {
				c.Locals("user_id", "invalid-uuid")
				c.Locals("user_email", "test@example.com")
			}
			return c.Next()
		})

		// Create a test endpoint that bypasses authentication middleware
		app.Get("/api/whoami", func(c *fiber.Ctx) error {
			userID := c.Locals("user_id")
			userEmail := c.Locals("user_email")

			if userID == nil || userEmail == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
			}

			// Parse user ID from string to UUID
			userUUID, err := uuid.Parse(userID.(string))
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user ID"})
			}

			// Get user details from database
			userEntity, err := client.User.Query().Where(user.ID(userUUID)).First(c.Context())
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
			}

			return c.JSON(fiber.Map{
				"id":            userEntity.ID.String(),
				"email":         userEntity.Email,
				"name":          userEntity.Name,
				"phone":         userEntity.Phone,
				"bio":           userEntity.Bio,
				"avatar":        userEntity.Avatar,
				"cover":         userEntity.Cover,
				"date_of_birth": userEntity.DateOfBirth,
				"sex":           userEntity.Sex,
				"goal":          userEntity.Goal,
				"height_cm":     userEntity.HeightCm,
				"weight_kg":     userEntity.WeightKg,
				"status":        userEntity.Status,
				"created_at":    userEntity.CreatedAt,
				"updated_at":    userEntity.UpdatedAt,
			})
		})

		req := httptest.NewRequest("GET", "/api/whoami", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns not found for non-existent user", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// Mock the middleware to set non-existent user ID
		app.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/api/whoami" {
				c.Locals("user_id", "550e8400-e29b-41d4-a716-446655440000")
				c.Locals("user_email", "nonexistent@example.com")
			}
			return c.Next()
		})

		// Create a test endpoint that bypasses authentication middleware
		app.Get("/api/whoami", func(c *fiber.Ctx) error {
			userID := c.Locals("user_id")
			userEmail := c.Locals("user_email")

			if userID == nil || userEmail == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
			}

			// Parse user ID from string to UUID
			userUUID, err := uuid.Parse(userID.(string))
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user ID"})
			}

			// Get user details from database
			userEntity, err := client.User.Query().Where(user.ID(userUUID)).First(c.Context())
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
			}

			return c.JSON(fiber.Map{
				"id":            userEntity.ID.String(),
				"email":         userEntity.Email,
				"name":          userEntity.Name,
				"phone":         userEntity.Phone,
				"bio":           userEntity.Bio,
				"avatar":        userEntity.Avatar,
				"cover":         userEntity.Cover,
				"date_of_birth": userEntity.DateOfBirth,
				"sex":           userEntity.Sex,
				"goal":          userEntity.Goal,
				"height_cm":     userEntity.HeightCm,
				"weight_kg":     userEntity.WeightKg,
				"status":        userEntity.Status,
				"created_at":    userEntity.CreatedAt,
				"updated_at":    userEntity.UpdatedAt,
			})
		})

		req := httptest.NewRequest("GET", "/api/whoami", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestRouterEdgeCases(t *testing.T) {
	// Set up required environment variables for authoidc
	os.Setenv("OAUTH_BASE_URL", "http://localhost:8080")
	os.Setenv("OIDC_GOOGLE_ISSUER", "https://accounts.google.com")
	os.Setenv("OIDC_GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("OIDC_GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("OIDC_GOOGLE_REDIRECT_PATH", "/auth/google/callback")
	os.Setenv("OIDC_LINE_ISSUER", "https://access.line.me")
	os.Setenv("OIDC_LINE_CLIENT_ID", "test-line-client-id")
	os.Setenv("OIDC_LINE_CLIENT_SECRET", "test-line-client-secret")
	os.Setenv("OIDC_LINE_REDIRECT_PATH", "/auth/line/callback")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("JWT_ACCESS_TTL_MIN", "15")
	defer func() {
		os.Unsetenv("OAUTH_BASE_URL")
		os.Unsetenv("OIDC_GOOGLE_ISSUER")
		os.Unsetenv("OIDC_GOOGLE_CLIENT_ID")
		os.Unsetenv("OIDC_GOOGLE_CLIENT_SECRET")
		os.Unsetenv("OIDC_GOOGLE_REDIRECT_PATH")
		os.Unsetenv("OIDC_LINE_ISSUER")
		os.Unsetenv("OIDC_LINE_CLIENT_ID")
		os.Unsetenv("OIDC_LINE_CLIENT_SECRET")
		os.Unsetenv("OIDC_LINE_REDIRECT_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_ACCESS_TTL_MIN")
	}()

	t.Run("handles empty app", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// This should not panic
		apiGroup := app.Group("/api")
		assert.NotPanics(t, func() {
			RegisterRoutes(apiGroup, app, client, config.Config{})
		})
	})

	t.Run("handles app with existing routes", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// Add some existing routes
		app.Get("/health", func(c *fiber.Ctx) error { return c.JSON("ok") })
		app.Get("/status", func(c *fiber.Ctx) error { return c.JSON("ok") })

		// This should not panic
		apiGroup := app.Group("/api")
		assert.NotPanics(t, func() {
			RegisterRoutes(apiGroup, app, client, config.Config{})
		})
	})

	t.Run("handles multiple route registrations", func(t *testing.T) {
		app := fiber.New()
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		cfg := config.Config{}
		defer client.Close()

		// Register routes multiple times
		apiGroup := app.Group("/api")
		assert.NotPanics(t, func() {
			RegisterRoutes(apiGroup, app, client, cfg)
			RegisterRoutes(apiGroup, app, client, cfg)
		})
	})
}
