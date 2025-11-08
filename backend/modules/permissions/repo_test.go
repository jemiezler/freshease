package permissions

import (
	"context"
	"testing"

	"freshease/backend/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestRepository_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test permissions
	permission1, err := client.Permission.Create().
		SetID(uuid.New()).
		SetCode("read_users").
		SetDescription("Permission to read user data").
		Save(ctx)
	require.NoError(t, err)

	permission2, err := client.Permission.Create().
		SetID(uuid.New()).
		SetCode("write_users").
		SetDescription("Permission to write user data").
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created permissions
	foundIDs := make(map[uuid.UUID]bool)
	for _, perm := range result {
		foundIDs[perm.ID] = true
		assert.NotEmpty(t, perm.Code)
	}

	assert.True(t, foundIDs[permission1.ID])
	assert.True(t, foundIDs[permission2.ID])
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test permission
	createdPermission, err := client.Permission.Create().
		SetID(uuid.New()).
		SetCode("read_users").
		SetDescription("Permission to read user data").
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdPermission.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdPermission.ID, result.ID)
	assert.Equal(t, "read_users", result.Code)
	assert.Equal(t, "Permission to read user data", result.Description)

	// Test FindByID - not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestRepository_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	dto := &CreatePermissionDTO{
		ID:          uuid.New(),
		Code:        "read_users",
		Description: stringPtr("Permission to read user data"),
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID) // ID should be generated
	assert.Equal(t, dto.Code, result.Code)
	assert.Equal(t, dto.Description, result.Description)

	// Verify it was actually created in the database
	dbPermission, err := client.Permission.Get(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, dto.Code, dbPermission.Code)
	assert.Equal(t, dto.Description, dbPermission.Description)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test permission
	createdPermission, err := client.Permission.Create().
		SetID(uuid.New()).
		SetCode("read_users").
		SetDescription("Permission to read user data").
		Save(ctx)
	require.NoError(t, err)

	// Test Update - full update
	dto := &UpdatePermissionDTO{
		ID:          createdPermission.ID,
		Code:        stringPtr("updated_code"),
		Description: stringPtr("Updated description"),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdPermission.ID, result.ID)
	assert.Equal(t, "updated_code", result.Code)
	assert.Equal(t, "Updated description", result.Description)

	// Verify it was actually updated in the database
	dbPermission, err := client.Permission.Get(ctx, createdPermission.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated_code", dbPermission.Code)
	assert.Equal(t, "Updated description", dbPermission.Description)

	// Test Update - partial update (only name)
	dto2 := &UpdatePermissionDTO{
		ID:   createdPermission.ID,
		Code: stringPtr("partial_update"),
	}

	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, createdPermission.ID, result2.ID)
	assert.Equal(t, "partial_update", result2.Code)
	assert.Equal(t, "Updated description", result2.Description) // Should remain unchanged

	// Test Update - no fields to update
	dto3 := &UpdatePermissionDTO{
		ID: createdPermission.ID,
		// No fields to update
	}

	_, err = repo.Update(ctx, dto3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
}

func TestRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test permission
	createdPermission, err := client.Permission.Create().
		SetID(uuid.New()).
		SetCode("read_users").
		SetDescription("Permission to read user data").
		Save(ctx)
	require.NoError(t, err)

	// Test Delete - success
	err = repo.Delete(ctx, createdPermission.ID)
	require.NoError(t, err)

	// Verify it was actually deleted from the database
	_, err = client.Permission.Get(ctx, createdPermission.ID)
	assert.Error(t, err)

	// Test Delete - non-existent ID
	nonExistentID := uuid.New()
	err = repo.Delete(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestRepository_Integration(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create multiple permissions
	dto1 := &CreatePermissionDTO{
		ID:          uuid.New(),
		Code:        "read_users",
		Description: stringPtr("Permission to read user data"),
	}

	dto2 := &CreatePermissionDTO{
		ID:          uuid.New(),
		Code:        "write_users",
		Description: stringPtr("Permission to write user data"),
	}

	permission1, err := repo.Create(ctx, dto1)
	require.NoError(t, err)

	permission2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)

	// List all permissions
	allPermissions, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allPermissions, 2)

	// Get specific permission
	retrievedPermission, err := repo.FindByID(ctx, permission1.ID)
	require.NoError(t, err)
	assert.Equal(t, permission1.ID, retrievedPermission.ID)
	assert.Equal(t, dto1.Code, retrievedPermission.Code)

	// Update permission
	updateDTO := &UpdatePermissionDTO{
		ID:   permission1.ID,
		Code: stringPtr("updated_read_users"),
	}

	updatedPermission, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, "updated_read_users", updatedPermission.Code)
	assert.Equal(t, dto1.Description, updatedPermission.Description)

	// Delete one permission
	err = repo.Delete(ctx, permission1.ID)
	require.NoError(t, err)

	// Verify only one permission remains
	remainingPermissions, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, remainingPermissions, 1)
	assert.Equal(t, permission2.ID, remainingPermissions[0].ID)
}
