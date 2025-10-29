package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("loads with default values", func(t *testing.T) {
		// Clear environment variables to test defaults
		envVars := []string{
			"DATABASE_URL", "HTTP_PORT", "JWT_SECRET", "ENT_DEBUG",
			"OIDC_GOOGLE_ISSUER", "OIDC_GOOGLE_CLIENT_ID", "OIDC_GOOGLE_CLIENT_SECRET",
			"OIDC_GOOGLE_REDIRECT_URI", "GENAI_APIKEY",
		}

		// Store original values
		originalValues := make(map[string]string)
		for _, envVar := range envVars {
			originalValues[envVar] = os.Getenv(envVar)
			os.Unsetenv(envVar)
		}

		// Restore original values after test
		defer func() {
			for envVar, value := range originalValues {
				if value != "" {
					os.Setenv(envVar, value)
				}
			}
		}()

		cfg := Load()

		// Test default values
		assert.Equal(t, "postgres://postgres:user1234@localhost:5432/trail-teller_db?sslmode=disable", cfg.DatabaseURL)
		assert.Equal(t, ":8080", cfg.HTTPPort)
		assert.Equal(t, "secret", cfg.JWTSecret)
		assert.False(t, cfg.Ent.Debug)
		assert.Empty(t, cfg.OIDC_GOOGLE_ISSUER)
		assert.Empty(t, cfg.OIDC_GOOGLE_CLIENT_ID)
		assert.Empty(t, cfg.OIDC_GOOGLE_CLIENT_SECRET)
		assert.Empty(t, cfg.OIDC_GOOGLE_REDIRECT_URI)
		assert.Empty(t, cfg.GENAI_APIKEY)
	})

	t.Run("loads with environment variables", func(t *testing.T) {
		// Set test environment variables
		testEnv := map[string]string{
			"DATABASE_URL":              "postgres://test:test@localhost:5432/test_db",
			"HTTP_PORT":                 ":9090",
			"JWT_SECRET":                "test-secret",
			"ENT_DEBUG":                 "true",
			"OIDC_GOOGLE_ISSUER":        "https://accounts.google.com",
			"OIDC_GOOGLE_CLIENT_ID":     "test-client-id",
			"OIDC_GOOGLE_CLIENT_SECRET": "test-client-secret",
			"OIDC_GOOGLE_REDIRECT_URI":  "http://localhost:8080/auth/google/callback",
			"GENAI_APIKEY":              "test-genai-key",
		}

		// Store original values
		originalValues := make(map[string]string)
		for envVar, value := range testEnv {
			originalValues[envVar] = os.Getenv(envVar)
			os.Setenv(envVar, value)
		}

		// Restore original values after test
		defer func() {
			for envVar, value := range originalValues {
				if value != "" {
					os.Setenv(envVar, value)
				} else {
					os.Unsetenv(envVar)
				}
			}
		}()

		cfg := Load()

		// Test environment variable values
		assert.Equal(t, "postgres://test:test@localhost:5432/test_db", cfg.DatabaseURL)
		assert.Equal(t, ":9090", cfg.HTTPPort)
		assert.Equal(t, "test-secret", cfg.JWTSecret)
		assert.True(t, cfg.Ent.Debug)
		assert.Equal(t, "https://accounts.google.com", cfg.OIDC_GOOGLE_ISSUER)
		assert.Equal(t, "test-client-id", cfg.OIDC_GOOGLE_CLIENT_ID)
		assert.Equal(t, "test-client-secret", cfg.OIDC_GOOGLE_CLIENT_SECRET)
		assert.Equal(t, "http://localhost:8080/auth/google/callback", cfg.OIDC_GOOGLE_REDIRECT_URI)
		assert.Equal(t, "test-genai-key", cfg.GENAI_APIKEY)
	})

	t.Run("handles mixed environment variables and defaults", func(t *testing.T) {
		// Set only some environment variables
		testEnv := map[string]string{
			"DATABASE_URL": "postgres://mixed:test@localhost:5432/mixed_db",
			"JWT_SECRET":   "mixed-secret",
			"ENT_DEBUG":    "true",
		}

		// Store original values
		originalValues := make(map[string]string)
		for envVar, value := range testEnv {
			originalValues[envVar] = os.Getenv(envVar)
			os.Setenv(envVar, value)
		}

		// Clear other variables
		otherVars := []string{"HTTP_PORT", "OIDC_GOOGLE_ISSUER", "OIDC_GOOGLE_CLIENT_ID", "OIDC_GOOGLE_CLIENT_SECRET", "OIDC_GOOGLE_REDIRECT_URI", "GENAI_APIKEY"}
		for _, envVar := range otherVars {
			originalValues[envVar] = os.Getenv(envVar)
			os.Unsetenv(envVar)
		}

		// Restore original values after test
		defer func() {
			for envVar, value := range originalValues {
				if value != "" {
					os.Setenv(envVar, value)
				} else {
					os.Unsetenv(envVar)
				}
			}
		}()

		cfg := Load()

		// Test mixed values
		assert.Equal(t, "postgres://mixed:test@localhost:5432/mixed_db", cfg.DatabaseURL)
		assert.Equal(t, ":8080", cfg.HTTPPort) // default
		assert.Equal(t, "mixed-secret", cfg.JWTSecret)
		assert.True(t, cfg.Ent.Debug)
		assert.Empty(t, cfg.OIDC_GOOGLE_ISSUER)        // default
		assert.Empty(t, cfg.OIDC_GOOGLE_CLIENT_ID)     // default
		assert.Empty(t, cfg.OIDC_GOOGLE_CLIENT_SECRET) // default
		assert.Empty(t, cfg.OIDC_GOOGLE_REDIRECT_URI)  // default
		assert.Empty(t, cfg.GENAI_APIKEY)              // default
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("returns environment variable when set", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test-value")
		defer os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default-value")
		assert.Equal(t, "test-value", result)
	})

	t.Run("returns default when environment variable not set", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default-value")
		assert.Equal(t, "default-value", result)
	})

	t.Run("returns default when environment variable is empty", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default-value")
		assert.Equal(t, "default-value", result)
	})

	t.Run("handles special characters in environment variables", func(t *testing.T) {
		specialValue := "test@#$%^&*()_+-=[]{}|;':\",./<>?"
		os.Setenv("TEST_SPECIAL", specialValue)
		defer os.Unsetenv("TEST_SPECIAL")

		result := getEnv("TEST_SPECIAL", "default")
		assert.Equal(t, specialValue, result)
	})
}

func TestEntConfig(t *testing.T) {
	t.Run("ent config debug flag", func(t *testing.T) {
		cfg := Config{
			Ent: EntConfig{
				Debug: true,
			},
		}

		assert.True(t, cfg.Ent.Debug)

		cfg.Ent.Debug = false
		assert.False(t, cfg.Ent.Debug)
	})
}

func TestConfigStruct(t *testing.T) {
	t.Run("config struct initialization", func(t *testing.T) {
		cfg := Config{
			DatabaseURL:               "test-db-url",
			HTTPPort:                  ":8080",
			JWTSecret:                 "test-secret",
			Ent:                       EntConfig{Debug: true},
			OIDC_GOOGLE_ISSUER:        "test-issuer",
			OIDC_GOOGLE_CLIENT_ID:     "test-client-id",
			OIDC_GOOGLE_CLIENT_SECRET: "test-client-secret",
			OIDC_GOOGLE_REDIRECT_URI:  "test-redirect-uri",
			GENAI_APIKEY:              "test-genai-key",
		}

		assert.Equal(t, "test-db-url", cfg.DatabaseURL)
		assert.Equal(t, ":8080", cfg.HTTPPort)
		assert.Equal(t, "test-secret", cfg.JWTSecret)
		assert.True(t, cfg.Ent.Debug)
		assert.Equal(t, "test-issuer", cfg.OIDC_GOOGLE_ISSUER)
		assert.Equal(t, "test-client-id", cfg.OIDC_GOOGLE_CLIENT_ID)
		assert.Equal(t, "test-client-secret", cfg.OIDC_GOOGLE_CLIENT_SECRET)
		assert.Equal(t, "test-redirect-uri", cfg.OIDC_GOOGLE_REDIRECT_URI)
		assert.Equal(t, "test-genai-key", cfg.GENAI_APIKEY)
	})
}
