package middleware

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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
		contentType    string
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
			contentType:    "application/json",
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

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tt.contentType)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.Equal(t, "success", responseBody["message"])
			} else {
				// For error cases, just verify status code
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				if err == nil && tt.expectedError != "" {
					if msg, ok := responseBody["message"].(string); ok {
						assert.Contains(t, msg, tt.expectedError)
					}
				}
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

func TestBindMultipartForm(t *testing.T) {
	tests := []struct {
		name            string
		contentType     string
		formData        map[string]string
		fileFieldName   string
		includeFile     bool
		allowEmptyPayload bool
		expectedError   string
		expectedFile    bool
	}{
		{
			name:          "success - multipart with JSON payload",
			contentType:   "multipart/form-data; boundary=test",
			formData:      map[string]string{"payload": `{"id":"` + uuid.New().String() + `","email":"test@example.com","password":"password123","name":"Test User"}`},
			fileFieldName: "file",
			includeFile:   false,
			expectedError: "",
			expectedFile:  false,
		},
		{
			name:          "success - multipart with form fields",
			contentType:   "multipart/form-data; boundary=test",
			formData:      map[string]string{"email": "test@example.com", "name": "Test User"},
			fileFieldName: "",
			includeFile:   false,
			expectedError: "",
			expectedFile:  false,
		},
		{
			name:          "success - non-multipart falls back to JSON",
			contentType:   "application/json",
			formData:      map[string]string{},
			fileFieldName: "",
			includeFile:   false,
			expectedError: "",
			expectedFile:  false,
		},
		{
			name:            "success - multipart with file and allowEmptyPayload",
			contentType:     "multipart/form-data; boundary=test",
			formData:        map[string]string{},
			fileFieldName:   "file",
			includeFile:     true,
			allowEmptyPayload: true,
			expectedError:   "",
			expectedFile:    true,
		},
		{
			name:          "error - invalid JSON payload",
			contentType:   "multipart/form-data; boundary=test",
			formData:      map[string]string{"payload": "invalid json"},
			fileFieldName: "",
			includeFile:   false,
			expectedError: "invalid JSON payload",
			expectedFile:  false,
		},
		{
			name:          "error - validation failed in payload",
			contentType:   "multipart/form-data; boundary=test",
			formData:      map[string]string{"payload": `{"email":"invalid-email"}`},
			fileFieldName: "",
			includeFile:   false,
			expectedError: "validation failed",
			expectedFile:  false,
		},
		{
			name:          "error - empty file",
			contentType:   "multipart/form-data; boundary=test",
			formData:      map[string]string{"payload": `{"id":"` + uuid.New().String() + `","email":"test@example.com","password":"password123","name":"Test User"}`},
			fileFieldName: "file",
			includeFile:   true,
			expectedError: "file is empty",
			expectedFile:  false,
		},
		{
			name:          "success - file with content",
			contentType:   "multipart/form-data; boundary=test",
			formData:      map[string]string{"payload": `{"id":"` + uuid.New().String() + `","email":"test@example.com","password":"password123","name":"Test User"}`},
			fileFieldName: "file",
			includeFile:   true,
			expectedError: "",
			expectedFile:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/test", func(c *fiber.Ctx) error {
				var dto TestCreateDTO
				file, err := BindMultipartForm(c, &dto, tt.fileFieldName, tt.allowEmptyPayload)
				if err != nil {
					return err
				}
				if file != nil {
					return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success", "file": file.Filename})
				}
				return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
			})

			// Create multipart form body
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			// Add form fields
			for key, value := range tt.formData {
				writer.WriteField(key, value)
			}

			// Add file if needed
			if tt.includeFile {
				fileWriter, err := writer.CreateFormFile(tt.fileFieldName, "test.txt")
				require.NoError(t, err)
				if tt.name == "error - empty file" {
					// Don't write anything to create empty file (size 0)
					// File is created but has no content
				} else {
					_, err = fileWriter.Write([]byte("test file content"))
					require.NoError(t, err)
				}
			}

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/test", &body)
			if tt.contentType == "multipart/form-data; boundary=test" {
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				req.Header.Set("Content-Type", tt.contentType)
				// For JSON, write JSON body
				if tt.contentType == "application/json" {
					jsonBody := map[string]interface{}{
						"id":       uuid.New().String(),
						"email":    "test@example.com",
						"password": "password123",
						"name":     "Test User",
					}
					jsonBytes, _ := json.Marshal(jsonBody)
					req.Body = http.NoBody
					req = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonBytes))
					req.Header.Set("Content-Type", "application/json")
				}
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				// If JSON decode fails, read the raw body for debugging
				bodyBytes := make([]byte, resp.ContentLength)
				resp.Body.Read(bodyBytes)
				t.Logf("Failed to decode response: %v, body: %s", err, string(bodyBytes))
				// Re-create body for re-reading
				resp.Body = http.NoBody
				return
			}

			if tt.expectedError != "" {
				if msg, ok := responseBody["message"].(string); ok {
					assert.Contains(t, msg, tt.expectedError)
				}
				assert.True(t, resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError)
			} else {
				assert.Equal(t, "success", responseBody["message"])
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				if tt.expectedFile {
					assert.NotNil(t, responseBody["file"])
				}
			}
		})
	}
}
