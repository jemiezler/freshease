package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test DTOs for validation testing
type TestCreateDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	Name     string    `json:"name" validate:"required,min=2,max=100"`
	Phone    *string   `json:"phone" validate:"omitempty,min=10,max=20"`
	Bio      *string   `json:"bio" validate:"omitempty,min=10,max=500"`
	Avatar   *string   `json:"avatar" validate:"omitempty,min=10,max=200"`
	Cover    *string   `json:"cover" validate:"omitempty,min=10,max=200"`
	Status   *string   `json:"status" validate:"omitempty"`
}

type TestUpdateDTO struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Email    *string   `json:"email" validate:"omitempty,email"`
	Password *string   `json:"password" validate:"omitempty,min=8,max=100"`
	Name     *string   `json:"name" validate:"omitempty,min=2,max=100"`
}

func TestBindAndValidate(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "success - valid create DTO",
			requestBody: TestCreateDTO{
				ID:       uuid.New(),
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
				Phone:    stringPtr("+1234567890"),
				Bio:      stringPtr("This is a longer test bio that meets the minimum length requirement"),
				Avatar:   stringPtr("https://example.com/avatar.jpg"),
				Cover:    stringPtr("https://example.com/cover.jpg"),
				Status:   stringPtr("active"),
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/test", func(c *fiber.Ctx) error {
				var dto TestCreateDTO
				if err := BindAndValidate(c, &dto); err != nil {
					return err
				}
				return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
			})

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.Contains(t, responseBody["message"].(string), tt.expectedError)
			} else {
				assert.Equal(t, "success", responseBody["message"])
			}
		})
	}
}

func TestRequireAuth(t *testing.T) {
	// Set up test JWT secret
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "success - valid token",
			authHeader:     "Bearer " + generateValidToken(t),
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "error - missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "missing bearer token",
		},
		{
			name:           "error - invalid bearer prefix",
			authHeader:     "Basic " + generateValidToken(t),
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "missing bearer token",
		},
		{
			name:           "error - invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid token",
		},
		{
			name:           "error - expired token",
			authHeader:     "Bearer " + generateExpiredToken(t),
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(RequireAuth())
			app.Get("/protected", func(c *fiber.Ctx) error {
				return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, responseBody["message"])
			} else {
				assert.Equal(t, "success", responseBody["message"])
			}
		})
	}
}

func TestRequestLogger(t *testing.T) {
	app := fiber.New()
	app.Use(RequestLogger())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if X-Response-Time header is set
	responseTime := resp.Header.Get("X-Response-Time")
	assert.NotEmpty(t, responseTime)
	// Response time should be a valid duration string
	assert.Regexp(t, `^\d+(\.\d+)?(ns|µs|ms|s)$`, responseTime)
}

func TestRequestLogger_WithError(t *testing.T) {
	app := fiber.New()
	app.Use(RequestLogger())
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error"})
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Check if X-Response-Time header is still set even with error
	responseTime := resp.Header.Get("X-Response-Time")
	assert.NotEmpty(t, responseTime)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func generateValidToken(t *testing.T) string {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secret)
	require.NoError(t, err)
	return tokenString
}

func generateExpiredToken(t *testing.T) string {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	require.NoError(t, err)
	return tokenString
}
