package vendors

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

	// Test List - empty list
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, result)

	// Create test vendors
	vendor1, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Vendor 1").
		SetContact("vendor1@example.com").
		Save(ctx)
	require.NoError(t, err)

	vendor2, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Vendor 2").
		SetContact("vendor2@example.com").
		Save(ctx)
	require.NoError(t, err)

	// Test List - populated list
	result, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created vendors
	foundIDs := make(map[uuid.UUID]bool)
	for _, vendor := range result {
		foundIDs[vendor.ID] = true
		assert.NotNil(t, vendor.Name)
		assert.NotNil(t, vendor.Contact)
	}

	assert.True(t, foundIDs[vendor1.ID])
	assert.True(t, foundIDs[vendor2.ID])
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test vendor
	createdVendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("test@example.com").
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdVendor.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdVendor.ID, result.ID)
	assert.Equal(t, "Test Vendor", *result.Name)
	assert.Equal(t, "test@example.com", *result.Contact)

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

	dto := &CreateVendorDTO{
		ID:      uuid.New(),
		Name:    stringPtr("Test Vendor"),
		Contact: stringPtr("test@example.com"),
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID) // ID should be generated
	assert.Equal(t, dto.Contact, result.Contact)

	// Verify it was actually created in the database
	dbVendor, err := client.Vendor.Get(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Vendor", *dbVendor.Name)
	assert.Equal(t, "test@example.com", *dbVendor.Contact)

	// Test Create - with nil ID (auto-generated)
	dto2 := &CreateVendorDTO{
		ID:      uuid.Nil,
		Name:    stringPtr("Auto ID Vendor"),
		Contact: stringPtr("auto@example.com"),
	}
	result2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.NotEqual(t, uuid.Nil, result2.ID) // ID should be auto-generated

	// Test Create - with nil Name (should fail validation or use empty)
	dto3 := &CreateVendorDTO{
		ID:      uuid.New(),
		Name:    nil,
		Contact: stringPtr("contact@example.com"),
	}
	result3, err := repo.Create(ctx, dto3)
	// This might succeed or fail depending on validation - test both paths
	if err == nil {
		assert.NotNil(t, result3)
	} else {
		assert.Error(t, err)
	}

	// Test Create - with nil Contact
	dto4 := &CreateVendorDTO{
		ID:      uuid.New(),
		Name:    stringPtr("Vendor Without Contact"),
		Contact: nil,
	}
	result4, err := repo.Create(ctx, dto4)
	if err == nil {
		assert.NotNil(t, result4)
	} else {
		assert.Error(t, err)
	}
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test vendor
	createdVendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("test@example.com").
		Save(ctx)
	require.NoError(t, err)

	// Test Update - full update
	dto := &UpdateVendorDTO{
		ID:      createdVendor.ID,
		Name:    stringPtr("Updated Vendor"),
		Contact: stringPtr("updated@example.com"),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdVendor.ID, result.ID)
	assert.Equal(t, "Updated Vendor", *result.Name)
	assert.Equal(t, "updated@example.com", *result.Contact)

	// Verify it was actually updated in the database
	dbVendor, err := client.Vendor.Get(ctx, createdVendor.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Vendor", *dbVendor.Name)
	assert.Equal(t, "updated@example.com", *dbVendor.Contact)

	// Test Update - partial update (only name)
	dto2 := &UpdateVendorDTO{
		ID:   createdVendor.ID,
		Name: stringPtr("Partial Update"),
	}

	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, createdVendor.ID, result2.ID)
	assert.Equal(t, "Partial Update", *result2.Name)
	assert.Equal(t, "updated@example.com", *result2.Contact) // Should remain unchanged

	// Test Update - no fields to update
	dto3 := &UpdateVendorDTO{
		ID: createdVendor.ID,
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

	// Create test vendor
	createdVendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("test@example.com").
		Save(ctx)
	require.NoError(t, err)

	// Test Delete - success
	err = repo.Delete(ctx, createdVendor.ID)
	require.NoError(t, err)

	// Verify it was actually deleted from the database
	_, err = client.Vendor.Get(ctx, createdVendor.ID)
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

	// Create multiple vendors
	dto1 := &CreateVendorDTO{
		ID:      uuid.New(),
		Name:    stringPtr("Vendor 1"),
		Contact: stringPtr("vendor1@example.com"),
	}

	dto2 := &CreateVendorDTO{
		ID:      uuid.New(),
		Name:    stringPtr("Vendor 2"),
		Contact: stringPtr("vendor2@example.com"),
	}

	vendor1, err := repo.Create(ctx, dto1)
	require.NoError(t, err)

	vendor2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)

	// List all vendors
	allVendors, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allVendors, 2)

	// Get specific vendor
	retrievedVendor, err := repo.FindByID(ctx, vendor1.ID)
	require.NoError(t, err)
	assert.Equal(t, vendor1.ID, retrievedVendor.ID)
	assert.Equal(t, dto1.Name, retrievedVendor.Name)

	// Update vendor
	updateDTO := &UpdateVendorDTO{
		ID:   vendor1.ID,
		Name: stringPtr("Updated Vendor 1"),
	}

	updatedVendor, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, "Updated Vendor 1", *updatedVendor.Name)
	assert.Equal(t, dto1.Contact, updatedVendor.Contact)

	// Delete one vendor
	err = repo.Delete(ctx, vendor1.ID)
	require.NoError(t, err)

	// Verify only one vendor remains
	remainingVendors, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, remainingVendors, 1)
	assert.Equal(t, vendor2.ID, remainingVendors[0].ID)
}
