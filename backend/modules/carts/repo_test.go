package carts

import (
	"context"
	"testing"

	"freshease/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user (required for cart)
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Test empty list
	carts, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, carts)

	// Create test carts with user
	cart1 := &CreateCartDTO{
		Status: stringPtr("pending"),
		Total:  float64Ptr(100.50),
		UserID: &user.ID,
	}

	cart2 := &CreateCartDTO{
		Status: stringPtr("completed"),
		Total:  float64Ptr(250.75),
		UserID: &user.ID,
	}

	createdCart1, err := repo.Create(ctx, cart1)
	require.NoError(t, err)

	createdCart2, err := repo.Create(ctx, cart2)
	require.NoError(t, err)

	// Test populated list
	carts, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, carts, 2)

	// Verify carts are returned
	cartMap := make(map[uuid.UUID]*GetCartDTO)
	for _, cart := range carts {
		cartMap[cart.ID] = cart
	}

	assert.Contains(t, cartMap, createdCart1.ID)
	assert.Contains(t, cartMap, createdCart2.ID)
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Test cart not found
	nonExistentID := uuid.New()
	_, err := repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)

	// Create test user (required for cart)
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test cart
	createDTO := &CreateCartDTO{
		Status: stringPtr("pending"),
		Total:  float64Ptr(150.25),
		UserID: &user.ID,
	}

	createdCart, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing cart
	foundCart, err := repo.FindByID(ctx, createdCart.ID)
	require.NoError(t, err)
	assert.Equal(t, createdCart.ID, foundCart.ID)
	assert.Equal(t, *createDTO.Status, foundCart.Status)
	assert.Equal(t, *createDTO.Total, foundCart.Total)
}

func TestRepository_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user (required for cart)
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	tests := []struct {
		name      string
		createDTO *CreateCartDTO
		wantError bool
	}{
		{
			name: "success - creates cart with all fields",
			createDTO: &CreateCartDTO{
				Status: stringPtr("pending"),
				Total:  float64Ptr(99.99),
				UserID: &user.ID,
			},
			wantError: false,
		},
		{
			name: "success - creates cart with minimal fields",
			createDTO: &CreateCartDTO{
				Status: stringPtr("completed"),
				Total:  float64Ptr(0.0),
				UserID: &user.ID,
			},
			wantError: false,
		},
		{
			name: "success - creates cart with nil total (uses default)",
			createDTO: &CreateCartDTO{
				Status: stringPtr("pending"), // Status is required
				Total:  nil,                  // Total can be nil (defaults to 0.0)
				UserID: &user.ID,             // UserID is required
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(ctx, tt.createDTO)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.ID)

				// Check fields
				if tt.createDTO.Status != nil {
					assert.Equal(t, *tt.createDTO.Status, result.Status)
				}
				if tt.createDTO.Total != nil {
					assert.Equal(t, *tt.createDTO.Total, result.Total)
				} else {
					assert.Equal(t, 0.0, result.Total) // default value
				}
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user (required for cart)
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create initial cart
	createDTO := &CreateCartDTO{
		Status: stringPtr("pending"),
		Total:  float64Ptr(100.00),
		UserID: &user.ID,
	}

	createdCart, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	tests := []struct {
		name      string
		updateDTO *UpdateCartDTO
		wantError bool
	}{
		{
			name: "success - updates all fields",
			updateDTO: &UpdateCartDTO{
				ID:     createdCart.ID,
				Status: stringPtr("completed"),
				Total:  float64Ptr(200.00),
			},
			wantError: false,
		},
		{
			name: "success - updates partial fields",
			updateDTO: &UpdateCartDTO{
				ID:     createdCart.ID,
				Status: stringPtr("cancelled"),
			},
			wantError: false,
		},
		{
			name: "error - updates non-existent cart",
			updateDTO: &UpdateCartDTO{
				ID:     uuid.New(),
				Status: stringPtr("completed"),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Update(ctx, tt.updateDTO)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.updateDTO.ID, result.ID)

				// Check updated fields
				if tt.updateDTO.Status != nil {
					assert.Equal(t, *tt.updateDTO.Status, result.Status)
				}
				if tt.updateDTO.Total != nil {
					assert.Equal(t, *tt.updateDTO.Total, result.Total)
				}
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user (required for cart)
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Test deleting non-existent cart
	nonExistentID := uuid.New()
	err = repo.Delete(ctx, nonExistentID)
	assert.Error(t, err)

	// Create test cart
	createDTO := &CreateCartDTO{
		Status: stringPtr("pending"),
		Total:  float64Ptr(75.50),
		UserID: &user.ID,
	}

	createdCart, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Verify cart exists
	_, err = repo.FindByID(ctx, createdCart.ID)
	require.NoError(t, err)

	// Test deleting existing cart
	err = repo.Delete(ctx, createdCart.ID)
	require.NoError(t, err)

	// Verify cart is deleted
	_, err = repo.FindByID(ctx, createdCart.ID)
	assert.Error(t, err)
}
