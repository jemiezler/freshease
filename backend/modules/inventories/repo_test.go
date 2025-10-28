package inventories

import (
	"context"
	"testing"
	"time"

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

	// Create test inventories
	inventory1, err := client.Inventory.Create().
		SetQuantity(100).
		SetRestockAmount(50).
		Save(ctx)
	require.NoError(t, err)

	inventory2, err := client.Inventory.Create().
		SetQuantity(200).
		SetRestockAmount(75).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created inventories
	foundIDs := make(map[uuid.UUID]bool)
	for _, inv := range result {
		foundIDs[inv.ID] = true
		assert.Greater(t, inv.Quantity, 0)
		assert.Greater(t, inv.RestockAmount, 0)
		assert.False(t, inv.UpdatedAt.IsZero())
	}

	assert.True(t, foundIDs[inventory1.ID])
	assert.True(t, foundIDs[inventory2.ID])
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test inventory
	createdInventory, err := client.Inventory.Create().
		SetQuantity(150).
		SetRestockAmount(60).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdInventory.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdInventory.ID, result.ID)
	assert.Equal(t, 150, result.Quantity)
	assert.Equal(t, 60, result.RestockAmount)
	assert.False(t, result.UpdatedAt.IsZero())

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

	dto := &CreateInventoryDTO{
		Quantity:      300,
		RestockAmount: 100,
		UpdatedAt:     time.Now(),
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, dto.Quantity, result.Quantity)
	assert.Equal(t, dto.RestockAmount, result.RestockAmount)
	assert.False(t, result.UpdatedAt.IsZero())

	// Verify it was actually created in the database
	dbInventory, err := client.Inventory.Get(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, dto.Quantity, dbInventory.Quantity)
	assert.Equal(t, dto.RestockAmount, dbInventory.RestockAmount)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test inventory
	createdInventory, err := client.Inventory.Create().
		SetQuantity(100).
		SetRestockAmount(50).
		Save(ctx)
	require.NoError(t, err)

	// Test Update - full update
	dto := &UpdateInventoryDTO{
		ID:            createdInventory.ID,
		Quantity:      intPtr(400),
		RestockAmount: intPtr(150),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdInventory.ID, result.ID)
	assert.Equal(t, 400, result.Quantity)
	assert.Equal(t, 150, result.RestockAmount)
	assert.False(t, result.UpdatedAt.IsZero())

	// Verify it was actually updated in the database
	dbInventory, err := client.Inventory.Get(ctx, createdInventory.ID)
	require.NoError(t, err)
	assert.Equal(t, 400, dbInventory.Quantity)
	assert.Equal(t, 150, dbInventory.RestockAmount)

	// Test Update - partial update (only quantity)
	dto2 := &UpdateInventoryDTO{
		ID:       createdInventory.ID,
		Quantity: intPtr(500),
	}

	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, createdInventory.ID, result2.ID)
	assert.Equal(t, 500, result2.Quantity)
	assert.Equal(t, 150, result2.RestockAmount) // Should remain unchanged

	// Test Update - no fields to update
	dto3 := &UpdateInventoryDTO{
		ID: createdInventory.ID,
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

	// Create test inventory
	createdInventory, err := client.Inventory.Create().
		SetQuantity(200).
		SetRestockAmount(80).
		Save(ctx)
	require.NoError(t, err)

	// Test Delete - success
	err = repo.Delete(ctx, createdInventory.ID)
	require.NoError(t, err)

	// Verify it was actually deleted from the database
	_, err = client.Inventory.Get(ctx, createdInventory.ID)
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

	// Create multiple inventories
	dto1 := &CreateInventoryDTO{
		Quantity:      100,
		RestockAmount: 50,
		UpdatedAt:     time.Now(),
	}

	dto2 := &CreateInventoryDTO{
		Quantity:      200,
		RestockAmount: 75,
		UpdatedAt:     time.Now(),
	}

	inventory1, err := repo.Create(ctx, dto1)
	require.NoError(t, err)

	inventory2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)

	// List all inventories
	allInventories, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allInventories, 2)

	// Get specific inventory
	retrievedInventory, err := repo.FindByID(ctx, inventory1.ID)
	require.NoError(t, err)
	assert.Equal(t, inventory1.ID, retrievedInventory.ID)
	assert.Equal(t, dto1.Quantity, retrievedInventory.Quantity)

	// Update inventory
	updateDTO := &UpdateInventoryDTO{
		ID:       inventory1.ID,
		Quantity: intPtr(300),
	}

	updatedInventory, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, 300, updatedInventory.Quantity)
	assert.Equal(t, dto1.RestockAmount, updatedInventory.RestockAmount)

	// Delete one inventory
	err = repo.Delete(ctx, inventory1.ID)
	require.NoError(t, err)

	// Verify only one inventory remains
	remainingInventories, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, remainingInventories, 1)
	assert.Equal(t, inventory2.ID, remainingInventories[0].ID)
}
