package authoidc

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	clients map[ProviderName]*providerClient
}

func (m *MockService) AuthCodeURL(p ProviderName, state, nonce, codeChallenge string) (string, error) {
	if _, ok := m.clients[p]; !ok {
		return "", errors.New("unknown provider")
	}
	return "https://example.com/auth?state=" + state, nil
}

func (m *MockService) ExchangeAndLogin(ctx context.Context, p ProviderName, code, codeVerifier string) (string, error) {
	if _, ok := m.clients[p]; !ok {
		return "", errors.New("unknown provider")
	}
	return "mock-jwt-token", nil
}

func TestController_Start(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		expectedStatus int
	}{
		{
			name:           "error - unknown provider",
			provider:       "unknown",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "error - empty provider",
			provider:       "",
			expectedStatus: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service for testing
			service := &MockService{
				clients: map[ProviderName]*providerClient{},
			}
			controller := NewController(service)
			app := fiber.New()
			app.Get("/auth/:provider/start", controller.Start)

			req := httptest.NewRequest(http.MethodGet, "/auth/"+tt.provider+"/start", nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestController_Callback(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		state          string
		code           string
		cookieState    string
		expectedStatus int
	}{
		{
			name:           "error - missing state",
			provider:       "google",
			state:          "",
			code:           "test-code",
			cookieState:    "test-state",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "error - state mismatch",
			provider:       "google",
			state:          "different-state",
			code:           "test-code",
			cookieState:    "test-state",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "error - unknown provider",
			provider:       "unknown",
			state:          "test-state",
			code:           "test-code",
			cookieState:    "test-state",
			expectedStatus: fiber.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service for testing
			service := &MockService{
				clients: map[ProviderName]*providerClient{},
			}
			controller := NewController(service)
			app := fiber.New()
			app.Get("/auth/:provider/callback", controller.Callback)

			req := httptest.NewRequest(http.MethodGet, "/auth/"+tt.provider+"/callback?state="+tt.state+"&code="+tt.code, nil)
			if tt.cookieState != "" {
				req.Header.Set("Cookie", "oidc_state="+tt.cookieState)
			}

			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestController_Integration(t *testing.T) {
	t.Run("test error cases", func(t *testing.T) {
		// Create a mock service for testing
		service := &MockService{
			clients: map[ProviderName]*providerClient{},
		}
		controller := NewController(service)
		app := fiber.New()
		app.Get("/auth/:provider/start", controller.Start)
		app.Get("/auth/:provider/callback", controller.Callback)

		// Test start flow with unknown provider
		startReq := httptest.NewRequest(http.MethodGet, "/auth/unknown/start", nil)
		startResp, err := app.Test(startReq)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, startResp.StatusCode)

		// Test callback flow with unknown provider
		callbackReq := httptest.NewRequest(http.MethodGet, "/auth/unknown/callback?state=test-state&code=test-code", nil)
		callbackReq.Header.Set("Cookie", "oidc_state=test-state")
		callbackResp, err := app.Test(callbackReq)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusTemporaryRedirect, callbackResp.StatusCode)
	})
}
