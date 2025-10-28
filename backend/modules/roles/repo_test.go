package roles

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

	// Create test roles
	role1, err := client.Role.Create().
		SetID(uuid.New()).
		SetName("admin").
		SetDescription("Administrator role").
		Save(ctx)
	require.NoError(t, err)

	role2, err := client.Role.Create().
		SetID(uuid.New()).
		SetName("user").
		SetDescription("Regular user role").
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created roles
	foundIDs := make(map[uuid.UUID]bool)
	for _, role := range result {
		foundIDs[role.ID] = true
		assert.NotEmpty(t, role.Name)
		assert.NotEmpty(t, role.Description)
	}

	assert.True(t, foundIDs[role1.ID])
	assert.True(t, foundIDs[role2.ID])
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test role
	createdRole, err := client.Role.Create().
		SetID(uuid.New()).
		SetName("admin").
		SetDescription("Administrator role").
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdRole.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdRole.ID, result.ID)
	assert.Equal(t, "admin", result.Name)
	assert.Equal(t, "Administrator role", result.Description)

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

	dto := &CreateRoleDTO{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID) // ID should be generated
	assert.Equal(t, dto.Name, result.Name)
	assert.Equal(t, dto.Description, result.Description)

	// Verify it was actually created in the database
	dbRole, err := client.Role.Get(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, dto.Name, dbRole.Name)
	assert.Equal(t, dto.Description, dbRole.Description)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test role
	createdRole, err := client.Role.Create().
		SetID(uuid.New()).
		SetName("admin").
		SetDescription("Administrator role").
		Save(ctx)
	require.NoError(t, err)

	// Test Update - full update
	dto := &UpdateRoleDTO{
		ID:          createdRole.ID,
		Name:        stringPtr("updated_admin"),
		Description: stringPtr("Updated description"),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdRole.ID, result.ID)
	assert.Equal(t, "updated_admin", result.Name)
	assert.Equal(t, "Updated description", result.Description)

	// Verify it was actually updated in the database
	dbRole, err := client.Role.Get(ctx, createdRole.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated_admin", dbRole.Name)
	assert.Equal(t, "Updated description", dbRole.Description)

	// Test Update - partial update (only name)
	dto2 := &UpdateRoleDTO{
		ID:   createdRole.ID,
		Name: stringPtr("partial_update"),
	}

	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, createdRole.ID, result2.ID)
	assert.Equal(t, "partial_update", result2.Name)
	assert.Equal(t, "Updated description", result2.Description) // Should remain unchanged

	// Test Update - no fields to update
	dto3 := &UpdateRoleDTO{
		ID: createdRole.ID,
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

	// Create test role
	createdRole, err := client.Role.Create().
		SetID(uuid.New()).
		SetName("admin").
		SetDescription("Administrator role").
		Save(ctx)
	require.NoError(t, err)

	// Test Delete - success
	err = repo.Delete(ctx, createdRole.ID)
	require.NoError(t, err)

	// Verify it was actually deleted from the database
	_, err = client.Role.Get(ctx, createdRole.ID)
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

	// Create multiple roles
	dto1 := &CreateRoleDTO{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
	}

	dto2 := &CreateRoleDTO{
		ID:          uuid.New(),
		Name:        "user",
		Description: "Regular user role",
	}

	role1, err := repo.Create(ctx, dto1)
	require.NoError(t, err)

	role2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)

	// List all roles
	allRoles, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allRoles, 2)

	// Get specific role
	retrievedRole, err := repo.FindByID(ctx, role1.ID)
	require.NoError(t, err)
	assert.Equal(t, role1.ID, retrievedRole.ID)
	assert.Equal(t, dto1.Name, retrievedRole.Name)

	// Update role
	updateDTO := &UpdateRoleDTO{
		ID:   role1.ID,
		Name: stringPtr("updated_admin"),
	}

	updatedRole, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, "updated_admin", updatedRole.Name)
	assert.Equal(t, dto1.Description, updatedRole.Description)

	// Delete one role
	err = repo.Delete(ctx, role1.ID)
	require.NoError(t, err)

	// Verify only one role remains
	remainingRoles, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, remainingRoles, 1)
	assert.Equal(t, role2.ID, remainingRoles[0].ID)
}
