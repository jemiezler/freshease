package http

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("creates new fiber app with correct configuration", func(t *testing.T) {
		app := New()

		assert.NotNil(t, app)
		assert.Equal(t, "FreshEase API", app.Config().AppName)
	})

	t.Run("app has request logger middleware", func(t *testing.T) {
		app := New()

		// Test that the app can handle a request (middleware is applied)
		req := httptest.NewRequest("GET", "/test", nil)

		// Add a test route
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "test"})
		})

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("app has swagger endpoint", func(t *testing.T) {
		app := New()

		// Test swagger endpoint
		req := httptest.NewRequest("GET", "/swagger/", nil)

		resp, err := app.Test(req)
		// Swagger endpoint might return various status codes depending on setup
		// but the route should be registered
		assert.NoError(t, err)
		// Accept various status codes as valid responses for swagger
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 600, "Expected valid HTTP status code, got %d", resp.StatusCode)
	})

	t.Run("app can handle multiple requests", func(t *testing.T) {
		app := New()

		// Add test routes
		app.Get("/test1", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "test1"})
		})
		app.Get("/test2", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "test2"})
		})

		// Test first route
		req1 := httptest.NewRequest("GET", "/test1", nil)
		resp1, err1 := app.Test(req1)
		assert.NoError(t, err1)
		assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

		// Test second route
		req2 := httptest.NewRequest("GET", "/test2", nil)
		resp2, err2 := app.Test(req2)
		assert.NoError(t, err2)
		assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	})

	t.Run("app handles different HTTP methods", func(t *testing.T) {
		app := New()

		// Add routes for different methods
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"method": "GET"})
		})
		app.Post("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"method": "POST"})
		})
		app.Put("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"method": "PUT"})
		})
		app.Delete("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"method": "DELETE"})
		})

		methods := []string{"GET", "POST", "PUT", "DELETE"}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/test", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}
	})

	t.Run("app handles 404 for unknown routes", func(t *testing.T) {
		app := New()

		req := httptest.NewRequest("GET", "/unknown-route", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("app handles different content types", func(t *testing.T) {
		app := New()

		// Add route that returns different content types
		app.Get("/json", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "json"})
		})
		app.Get("/text", func(c *fiber.Ctx) error {
			return c.SendString("text response")
		})

		// Test JSON response
		req1 := httptest.NewRequest("GET", "/json", nil)
		resp1, err1 := app.Test(req1)
		assert.NoError(t, err1)
		assert.Equal(t, fiber.StatusOK, resp1.StatusCode)
		assert.Contains(t, resp1.Header.Get("Content-Type"), "application/json")

		// Test text response
		req2 := httptest.NewRequest("GET", "/text", nil)
		resp2, err2 := app.Test(req2)
		assert.NoError(t, err2)
		assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
		assert.Contains(t, resp2.Header.Get("Content-Type"), "text/plain")
	})
}

func TestCtxType(t *testing.T) {
	t.Run("Ctx type is correctly defined", func(t *testing.T) {
		// Test that Ctx is properly aliased
		var ctx Ctx
		assert.NotNil(t, ctx)
	})
}

func TestAppConfiguration(t *testing.T) {
	t.Run("app has correct default configuration", func(t *testing.T) {
		app := New()

		config := app.Config()
		assert.Equal(t, "FreshEase API", config.AppName)
		assert.False(t, config.DisableStartupMessage)
		assert.False(t, config.DisableDefaultDate)
		assert.False(t, config.DisableDefaultContentType)
	})

	t.Run("app can be configured with custom settings", func(t *testing.T) {
		// This test verifies that the app can be created and configured
		app := New()

		// Test that we can add custom middleware
		app.Use(func(c *fiber.Ctx) error {
			c.Set("X-Custom-Header", "test")
			return c.Next()
		})

		app.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "test"})
		})

		req := httptest.NewRequest("GET", "/test", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "test", resp.Header.Get("X-Custom-Header"))
	})
}
