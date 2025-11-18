package authoidc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestStateStore_Delete(t *testing.T) {
	store := NewStateStore()
	state := "test-state-123"
	
	// Store a state
	store.Store(state)
	
	// Verify it exists
	assert.True(t, store.Validate(state))
	
	// Delete it
	store.Delete(state)
	
	// Verify it's gone
	assert.False(t, store.Validate(state))
}

func TestStateStore_Validate_Expired(t *testing.T) {
	store := &StateStore{
		states: make(map[string]time.Time),
	}
	
	// Store a state with expired time
	expiredState := "expired-state"
	store.mutex.Lock()
	store.states[expiredState] = time.Now().Add(-1 * time.Minute) // Expired 1 minute ago
	store.mutex.Unlock()
	
	// Validate should return false for expired state
	assert.False(t, store.Validate(expiredState))
}

func TestController_Exchange(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		requestBody    map[string]string
		mockSetup      func(*MockService)
		expectedStatus int
	}{
		{
			name:     "error - invalid request body",
			provider: "google",
			requestBody: map[string]string{
				"invalid": "json",
			},
			mockSetup:      func(*MockService) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "error - missing code",
			provider: "google",
			requestBody: map[string]string{
				"state": "test-state",
			},
			mockSetup:      func(*MockService) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "error - missing state",
			provider: "google",
			requestBody: map[string]string{
				"code": "test-code",
			},
			mockSetup:      func(*MockService) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "error - invalid state",
			provider: "google",
			requestBody: map[string]string{
				"code":  "test-code",
				"state": "invalid-state",
			},
			mockSetup:      func(*MockService) {},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:     "error - unknown provider",
			provider: "unknown",
			requestBody: map[string]string{
				"code":  "test-code",
				"state": "test-state",
			},
			mockSetup: func(ms *MockService) {
				// Mock service doesn't have unknown provider
			},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:     "success - valid exchange",
			provider: "google",
			requestBody: map[string]string{
				"code":  "test-code",
				"state": "test-state",
			},
			mockSetup: func(ms *MockService) {
				ms.clients[ProviderGoogle] = &providerClient{}
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &MockService{
				clients: map[ProviderName]*providerClient{},
			}
			tt.mockSetup(service)
			controller := NewController(service)
			
			// Store state if needed
			if tt.requestBody["state"] != "" && tt.name != "error - invalid state" {
				controller.stateStore.Store(tt.requestBody["state"])
			}
			
			app := fiber.New()
			app.Post("/auth/:provider/exchange", controller.Exchange)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/"+tt.provider+"/exchange", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Set cookie for state validation fallback
			if tt.requestBody["state"] != "" {
				req.Header.Set("Cookie", "oidc_state="+tt.requestBody["state"])
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestController_Start_Success(t *testing.T) {
	service := &MockService{
		clients: map[ProviderName]*providerClient{
			ProviderGoogle: {},
		},
	}
	controller := NewController(service)
	app := fiber.New()
	app.Get("/auth/:provider/start", controller.Start)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/start", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTemporaryRedirect, resp.StatusCode)
	
	// Verify cookies are set
	setCookies := resp.Header.Values("Set-Cookie")
	hasState := false
	hasNonce := false
	for _, cookie := range setCookies {
		if len(cookie) > 0 && (cookie[:10] == "oidc_state" || 
			(len(cookie) > 10 && cookie[1:11] == "oidc_state")) {
			hasState = true
		}
		if len(cookie) > 0 && (cookie[:10] == "oidc_nonce" || 
			(len(cookie) > 10 && cookie[1:11] == "oidc_nonce")) {
			hasNonce = true
		}
	}
	assert.True(t, hasState, "oidc_state cookie should be set")
	assert.True(t, hasNonce, "oidc_nonce cookie should be set")
}


func TestController_Callback_Success(t *testing.T) {
	service := &MockService{
		clients: map[ProviderName]*providerClient{},
	}
	controller := NewController(service)
	state := "test-state-123"
	controller.stateStore.Store(state)
	
	app := fiber.New()
	app.Get("/auth/:provider/callback", controller.Callback)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state="+state+"&code=test-code", nil)
	req.Header.Set("Cookie", "oidc_state="+state)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTemporaryRedirect, resp.StatusCode)
	
	// Verify redirect URL contains code and state
	location := resp.Header.Get("Location")
	assert.Contains(t, location, "freshease://callback")
	assert.Contains(t, location, "code=test-code")
	assert.Contains(t, location, "state="+state)
}

func TestController_Start_Web(t *testing.T) {
	service := &MockService{
		clients: map[ProviderName]*providerClient{
			ProviderGoogle: {},
		},
	}
	controller := NewController(service)
	app := fiber.New()
	app.Get("/auth/:provider/start", controller.Start)

	callbackURL := "http://localhost:8080/auth/callback"
	req := httptest.NewRequest(http.MethodGet, "/auth/google/start?platform=web&callback_url="+callbackURL, nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTemporaryRedirect, resp.StatusCode)
	
	// Verify cookies are set
	setCookies := resp.Header.Values("Set-Cookie")
	hasState := false
	hasNonce := false
	for _, cookie := range setCookies {
		if len(cookie) > 0 && (cookie[:10] == "oidc_state" || 
			(len(cookie) > 10 && cookie[1:11] == "oidc_state")) {
			hasState = true
		}
		if len(cookie) > 0 && (cookie[:10] == "oidc_nonce" || 
			(len(cookie) > 10 && cookie[1:11] == "oidc_nonce")) {
			hasNonce = true
		}
	}
	assert.True(t, hasState, "oidc_state cookie should be set")
	assert.True(t, hasNonce, "oidc_nonce cookie should be set")
}

func TestController_Callback_Web(t *testing.T) {
	service := &MockService{
		clients: map[ProviderName]*providerClient{},
	}
	controller := NewController(service)
	state := "test-state-456"
	webCallbackURL := "http://localhost:8080/auth/callback"
	controller.stateStore.StoreWithCallback(state, webCallbackURL)
	
	app := fiber.New()
	app.Get("/auth/:provider/callback", controller.Callback)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state="+state+"&code=test-code-web", nil)
	req.Header.Set("Cookie", "oidc_state="+state)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusTemporaryRedirect, resp.StatusCode)
	
	// Verify redirect URL is the web callback URL with code and state
	location := resp.Header.Get("Location")
	assert.Contains(t, location, webCallbackURL)
	assert.Contains(t, location, "code=test-code-web")
	assert.Contains(t, location, "state="+state)
	assert.NotContains(t, location, "freshease://")
}
