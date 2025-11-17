package bundles

import (
	"context"
	"testing"

	"freshease/backend/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestEntRepo_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test bundles
	bundle1, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Bundle One").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	bundle2, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Bundle Two").
		SetPrice(149.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created bundles
	foundIDs := make(map[uuid.UUID]bool)
	for _, bundle := range result {
		foundIDs[bundle.ID] = true
		assert.NotEmpty(t, bundle.Name)
		assert.Greater(t, bundle.Price, 0.0)
	}

	assert.True(t, foundIDs[bundle1.ID])
	assert.True(t, foundIDs[bundle2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test bundle
	createdBundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Test Bundle").
		SetPrice(199.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdBundle.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdBundle.ID, result.ID)
	assert.Equal(t, "Test Bundle", result.Name)
	assert.Equal(t, 199.99, result.Price)
	assert.True(t, result.IsActive)

	// Test FindByID - not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestEntRepo_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	desc := "Test bundle description"
	dto := &CreateBundleDTO{
		ID:          uuid.New(),
		Name:        "New Bundle",
		Description: &desc,
		Price:       299.99,
		IsActive:    true,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Name, result.Name)
	assert.Equal(t, *dto.Description, *result.Description)
	assert.Equal(t, dto.Price, result.Price)
	assert.Equal(t, dto.IsActive, result.IsActive)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test bundle
	createdBundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Original Bundle").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	// Update bundle
	newName := "Updated Bundle"
	newPrice := 199.99
	desc := "Updated description"
	dto := &UpdateBundleDTO{
		ID:          createdBundle.ID,
		Name:        &newName,
		Price:       &newPrice,
		Description: &desc,
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdBundle.ID, result.ID)
	assert.Equal(t, "Updated Bundle", result.Name)
	assert.Equal(t, 199.99, result.Price)
	assert.Equal(t, "Updated description", *result.Description)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test bundle
	createdBundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("To Delete").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	// Delete bundle
	err = repo.Delete(ctx, createdBundle.ID)
	require.NoError(t, err)

	// Verify bundle is deleted
	_, err = repo.FindByID(ctx, createdBundle.ID)
	assert.Error(t, err)
}

