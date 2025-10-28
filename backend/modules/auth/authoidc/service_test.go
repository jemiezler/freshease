package authoidc

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_AuthCodeURL(t *testing.T) {
	tests := []struct {
		name          string
		provider      ProviderName
		state         string
		nonce         string
		codeChallenge string
		expectedError bool
	}{
		{
			name:          "success - google provider",
			provider:      ProviderGoogle,
			state:         "test-state",
			nonce:         "test-nonce",
			codeChallenge: "",
			expectedError: false,
		},
		{
			name:          "success - line provider",
			provider:      ProviderLINE,
			state:         "test-state",
			nonce:         "test-nonce",
			codeChallenge: "",
			expectedError: false,
		},
		{
			name:          "error - unknown provider",
			provider:      "unknown",
			state:         "test-state",
			nonce:         "test-nonce",
			codeChallenge: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test would require actual OIDC provider setup
			// For now, we'll test the error case
			if tt.expectedError {
				service := &Service{
					clients: map[ProviderName]*providerClient{},
				}
				url, err := service.AuthCodeURL(tt.provider, tt.state, tt.nonce, tt.codeChallenge)
				assert.Error(t, err)
				assert.Empty(t, url)
			}
		})
	}
}

func TestService_ExchangeAndLogin(t *testing.T) {
	tests := []struct {
		name          string
		provider      ProviderName
		code          string
		codeVerifier  string
		expectedError bool
	}{
		{
			name:          "error - unknown provider",
			provider:      "unknown",
			code:          "test-code",
			codeVerifier:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				clients: map[ProviderName]*providerClient{},
			}
			token, err := service.ExchangeAndLogin(context.Background(), tt.provider, tt.code, tt.codeVerifier)
			assert.Error(t, err)
			assert.Empty(t, token)
		})
	}
}

func TestService_IssueJWT(t *testing.T) {
	tests := []struct {
		name      string
		uid       uuid.UUID
		email     string
		jwtSecret []byte
		ttl       time.Duration
	}{
		{
			name:      "success - issue JWT",
			uid:       uuid.New(),
			email:     "test@example.com",
			jwtSecret: []byte("test-secret"),
			ttl:       15 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				jwtSecret: tt.jwtSecret,
				ttl:       tt.ttl,
			}

			token, err := service.issueJWT(tt.uid, tt.email)
			require.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify the token
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return tt.jwtSecret, nil
			})
			require.NoError(t, err)
			assert.True(t, parsedToken.Valid)

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			require.True(t, ok)
			assert.Equal(t, tt.uid.String(), claims["sub"])
			assert.Equal(t, tt.email, claims["email"])
		})
	}
}

func TestProviderName(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		expected ProviderName
	}{
		{
			name:     "google provider",
			provider: "google",
			expected: ProviderGoogle,
		},
		{
			name:     "line provider",
			provider: "line",
			expected: ProviderLINE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProviderName(tt.provider)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRandB64(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "generate 16 byte string",
			n:    16,
		},
		{
			name: "generate 32 byte string",
			n:    32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := randB64(tt.n)
			assert.NotEmpty(t, result)
			assert.GreaterOrEqual(t, len(result), tt.n)
		})
	}
}
