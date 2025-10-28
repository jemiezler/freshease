package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRepository is a generic mock repository for testing
type MockRepository[T any] struct {
	ListFunc     func(ctx context.Context) ([]*T, error)
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*T, error)
	CreateFunc   func(ctx context.Context, dto interface{}) (*T, error)
	UpdateFunc   func(ctx context.Context, dto interface{}) (*T, error)
	DeleteFunc   func(ctx context.Context, id uuid.UUID) error
}

// TestContext provides common testing utilities
type TestContext struct {
	t   *testing.T
	app *fiber.App
}

// NewTestContext creates a new test context
func NewTestContext(t *testing.T, app *fiber.App) *TestContext {
	return &TestContext{
		t:   t,
		app: app,
	}
}

// MakeRequest makes an HTTP request and returns the response using Fiber's testing method
func (tc *TestContext) MakeRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(tc.t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(tc.t, err)

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return tc.app.Test(req)
}

// AssertJSONResponse asserts that the response contains expected JSON
func (tc *TestContext) AssertJSONResponse(resp *http.Response, expectedStatus int, expectedBody interface{}) {
	assert.Equal(tc.t, expectedStatus, resp.StatusCode)

	if expectedBody != nil {
		var actualBody interface{}
		err := json.NewDecoder(resp.Body).Decode(&actualBody)
		require.NoError(tc.t, err)

		expectedJSON, err := json.Marshal(expectedBody)
		require.NoError(tc.t, err)

		var expectedBodyParsed interface{}
		err = json.Unmarshal(expectedJSON, &expectedBodyParsed)
		require.NoError(tc.t, err)

		assert.Equal(tc.t, expectedBodyParsed, actualBody)
	}
}

// CreateTestUserDTO creates a test user DTO with valid data
func CreateTestUserDTO() map[string]interface{} {
	return map[string]interface{}{
		"id":       uuid.New().String(),
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
		"phone":    "+1234567890",
		"bio":      "Test bio",
		"avatar":   "https://example.com/avatar.jpg",
		"cover":    "https://example.com/cover.jpg",
		"status":   "active",
	}
}

// CreateTestProductDTO creates a test product DTO with valid data
func CreateTestProductDTO() map[string]interface{} {
	return map[string]interface{}{
		"id":          uuid.New().String(),
		"name":        "Test Product",
		"description": "Test product description",
		"price":       99.99,
		"sku":         "TEST-SKU-001",
		"status":      "active",
	}
}

// AssertErrorResponse asserts that the response contains an error message
func (tc *TestContext) AssertErrorResponse(resp *http.Response, expectedStatus int, expectedMessage string) {
	assert.Equal(tc.t, expectedStatus, resp.StatusCode)

	var errorResponse map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&errorResponse)
	require.NoError(tc.t, err)

	if expectedMessage != "" {
		assert.Equal(tc.t, expectedMessage, errorResponse["message"])
	}
}

// GenerateTestUUID generates a new UUID for testing
func GenerateTestUUID() uuid.UUID {
	return uuid.New()
}

// CreateFiberTestApp creates a Fiber app for testing
func CreateFiberTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	})
	return app
}

// MockHTTPHandler creates a mock HTTP handler for testing
func MockHTTPHandler(statusCode int, body interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(statusCode).JSON(body)
	}
}

// AssertContains asserts that a slice contains a specific element
func AssertContains[T comparable](t *testing.T, slice []T, element T) {
	for _, item := range slice {
		if item == element {
			return
		}
	}
	t.Errorf("Expected slice to contain %v, but it didn't", element)
}

// AssertNotContains asserts that a slice does not contain a specific element
func AssertNotContains[T comparable](t *testing.T, slice []T, element T) {
	for _, item := range slice {
		if item == element {
			t.Errorf("Expected slice to not contain %v, but it did", element)
			return
		}
	}
}

// TestRequest performs a test request using Fiber's built-in testing method
func TestRequest(app *fiber.App, req *http.Request) (*http.Response, error) {
	return app.Test(req)
}

// AssertValidUUID asserts that a string is a valid UUID
func AssertValidUUID(t *testing.T, uuidStr string) {
	_, err := uuid.Parse(uuidStr)
	assert.NoError(t, err, "Expected valid UUID, got: %s", uuidStr)
}

// CreateTestContext creates a test context with timeout
func CreateTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// AssertStatusCode asserts the HTTP status code
func AssertStatusCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected status code %d, got %d", expected, actual)
}

// AssertContentType asserts the content type header
func AssertContentType(t *testing.T, resp *http.Response, expected string) {
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, expected, "Expected content type to contain %s, got %s", expected, contentType)
}
