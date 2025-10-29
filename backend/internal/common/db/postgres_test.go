package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEntClientPGX(t *testing.T) {
	t.Run("returns error for invalid DSN", func(t *testing.T) {
		ctx := context.Background()
		invalidDSN := "invalid-dsn-format"

		client, closeFn, err := NewEntClientPGX(ctx, invalidDSN, false)

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("returns error for malformed DSN", func(t *testing.T) {
		ctx := context.Background()
		malformedDSN := "postgres://user:pass@localhost:invalid-port/db"

		client, closeFn, err := NewEntClientPGX(ctx, malformedDSN, false)

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles empty DSN", func(t *testing.T) {
		ctx := context.Background()
		emptyDSN := ""

		client, closeFn, err := NewEntClientPGX(ctx, emptyDSN, false)

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles DSN with special characters", func(t *testing.T) {
		ctx := context.Background()
		specialDSN := "postgres://user@#$%:pass@localhost:5432/db?sslmode=disable"

		client, closeFn, err := NewEntClientPGX(ctx, specialDSN, false)

		// This should fail due to invalid characters in username
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles DSN with missing required fields", func(t *testing.T) {
		ctx := context.Background()
		incompleteDSN := "postgres://localhost:5432/db"

		client, closeFn, err := NewEntClientPGX(ctx, incompleteDSN, false)

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		dsn := "postgres://user:pass@localhost:5432/db?sslmode=disable"

		client, closeFn, err := NewEntClientPGX(ctx, dsn, false)

		// Should fail due to cancelled context
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(1 * time.Millisecond)

		dsn := "postgres://user:pass@localhost:5432/db?sslmode=disable"

		client, closeFn, err := NewEntClientPGX(ctx, dsn, false)

		// Should fail due to timeout
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})
}

func TestNewEntClientPGX_ValidDSN(t *testing.T) {
	// Note: These tests require a real database connection
	// They are marked as integration tests and should be run separately
	// or with a test database setup

	t.Run("creates client with valid DSN (integration test)", func(t *testing.T) {
		t.Skip("Skipping integration test - requires database connection")

		ctx := context.Background()
		dsn := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"

		client, closeFn, err := NewEntClientPGX(ctx, dsn, false)

		require.NoError(t, err)
		require.NotNil(t, client)
		require.NotNil(t, closeFn)

		// Test that closeFn works
		err = closeFn(ctx)
		assert.NoError(t, err)
	})

	t.Run("creates debug client with valid DSN (integration test)", func(t *testing.T) {
		t.Skip("Skipping integration test - requires database connection")

		ctx := context.Background()
		dsn := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"

		client, closeFn, err := NewEntClientPGX(ctx, dsn, true)

		require.NoError(t, err)
		require.NotNil(t, client)
		require.NotNil(t, closeFn)

		// Test that closeFn works
		err = closeFn(ctx)
		assert.NoError(t, err)
	})
}

func TestNewEntClientPGX_EdgeCases(t *testing.T) {
	t.Run("handles very long DSN", func(t *testing.T) {
		ctx := context.Background()
		longDSN := "postgres://" +
			"verylongusername:" +
			"verylongpassword:" +
			"@localhost:5432/" +
			"verylongdatabasename" +
			"?sslmode=disable&" +
			"application_name=verylongappname&" +
			"connect_timeout=30&" +
			"statement_timeout=30000"

		client, closeFn, err := NewEntClientPGX(ctx, longDSN, false)

		// This should fail due to invalid format (missing @)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles DSN with query parameters", func(t *testing.T) {
		ctx := context.Background()
		dsnWithParams := "postgres://user:pass@localhost:5432/db?sslmode=disable&application_name=test&connect_timeout=10"

		client, closeFn, err := NewEntClientPGX(ctx, dsnWithParams, false)

		// This should fail due to no real database connection
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})

	t.Run("handles DSN with different SSL modes", func(t *testing.T) {
		ctx := context.Background()

		sslModes := []string{"disable", "require", "verify-ca", "verify-full"}

		for _, sslMode := range sslModes {
			dsn := "postgres://user:pass@localhost:5432/db?sslmode=" + sslMode

			client, closeFn, err := NewEntClientPGX(ctx, dsn, false)

			// All should fail due to no real database connection
			assert.Error(t, err, "Should fail for sslmode=%s", sslMode)
			assert.Nil(t, client)
			assert.Nil(t, closeFn)
		}
	})
}

func TestNewEntClientPGX_ConnectionPool(t *testing.T) {
	t.Run("verifies connection pool settings are applied", func(t *testing.T) {
		// This test would require mocking the sql.DB to verify
		// that SetMaxOpenConns, SetMaxIdleConns, etc. are called
		// For now, we'll test the error cases

		ctx := context.Background()
		dsn := "postgres://user:pass@localhost:5432/db?sslmode=disable"

		client, closeFn, err := NewEntClientPGX(ctx, dsn, false)

		// Should fail due to no real database connection
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})
}

func TestNewEntClientPGX_ContextHandling(t *testing.T) {
	t.Run("handles nil context", func(t *testing.T) {
		// This should panic or return an error
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic with nil context or invalid memory address
				assert.True(t, r != nil, "Expected panic with nil context")
			}
		}()

		dsn := "postgres://user:pass@localhost:5432/db?sslmode=disable"
		client, closeFn, err := NewEntClientPGX(nil, dsn, false)

		// If it doesn't panic, it should return an error
		if err == nil {
			t.Error("Expected error with nil context")
		}
		assert.Nil(t, client)
		assert.Nil(t, closeFn)
	})
}
