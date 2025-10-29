package authoidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ServiceInterface interface {
	AuthCodeURL(p ProviderName, state, nonce, codeChallenge string) (string, error)
	ExchangeAndLogin(ctx context.Context, p ProviderName, code, codeVerifier string) (string, error)
}

type Controller struct {
	s          ServiceInterface
	stateStore *StateStore
}

type StateStore struct {
	states map[string]time.Time
	mutex  sync.RWMutex
}

func NewStateStore() *StateStore {
	store := &StateStore{
		states: make(map[string]time.Time),
	}
	// Start cleanup goroutine
	go store.cleanup()
	return store
}

func (s *StateStore) Store(state string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.states[state] = time.Now().Add(10 * time.Minute) // 10 minute expiration
}

func (s *StateStore) Validate(state string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	expiry, exists := s.states[state]
	if !exists {
		return false
	}
	return time.Now().Before(expiry)
}

func (s *StateStore) Delete(state string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.states, state)
}

func (s *StateStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for state, expiry := range s.states {
			if now.After(expiry) {
				delete(s.states, state)
			}
		}
		s.mutex.Unlock()
	}
}

func NewController(s ServiceInterface) *Controller {
	return &Controller{
		s:          s,
		stateStore: NewStateStore(),
	}
}

func randB64(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// GET /api/auth/:provider/start
func (ctl *Controller) Start(c *fiber.Ctx) error {
	p := ProviderName(c.Params("provider"))
	state := randB64(16)
	nonce := randB64(16)
	ctl.stateStore.Store(state)
	c.Cookie(&fiber.Cookie{Name: "oidc_state", Value: state, HTTPOnly: true, SameSite: "Lax", Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "oidc_nonce", Value: nonce, HTTPOnly: true, SameSite: "Lax", Path: "/"})

	url, err := ctl.s.AuthCodeURL(p, state, nonce, "")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// GET /api/auth/:provider/callback?code=...&state=...
func (ctl *Controller) Callback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "missing code or state"})
	}

	// Validate state against our memory store (primary) or cookie (fallback)
	stateValid := ctl.stateStore.Validate(state) || c.Cookies("oidc_state") == state
	if !stateValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid state"})
	}

	// For mobile apps, redirect to custom scheme with the authorization code
	// The mobile app will then call the exchange endpoint
	redirectURL := fmt.Sprintf("freshease://callback?code=%s&state=%s", code, state)
	return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
}

// POST /api/auth/:provider/exchange
func (ctl *Controller) Exchange(c *fiber.Ctx) error {
	p := ProviderName(c.Params("provider"))

	var req struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if req.Code == "" || req.State == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Missing code or state"})
	}

	// Verify state against our memory store (primary) or cookie (fallback)
	stateValid := ctl.stateStore.Validate(req.State) || req.State == c.Cookies("oidc_state")
	if !stateValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid state"})
	}

	fmt.Printf("🔄 [OAuth Exchange] Calling ExchangeAndLogin...\n")
	access, err := ctl.s.ExchangeAndLogin(c.Context(), p, req.Code, "")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	// Clean up state from memory store
	ctl.stateStore.Delete(req.State)

	// Clear cookies
	c.Cookie(&fiber.Cookie{Name: "oidc_state", Value: "", MaxAge: -1, Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "oidc_nonce", Value: "", MaxAge: -1, Path: "/"})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    fiber.Map{"accessToken": access},
		"message": "Authentication successful",
	})
}
