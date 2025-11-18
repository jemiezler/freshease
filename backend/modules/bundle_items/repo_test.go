package bundle_items

import (
	"context"
	"testing"

	"freshease/backend/ent/enttest"
	"freshease/backend/internal/common/errs"
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

	// Test empty list
	items, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, items)

	// Create required entities
	bundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Test Bundle").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product1, err := client.Product.Create().
		SetName("Product 1").
		SetSku("SKU-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	product2, err := client.Product.Create().
		SetName("Product 2").
		SetSku("SKU-002").
		SetPrice(5.50).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create bundle items
	item1, err := client.Bundle_item.Create().
		SetID(uuid.New()).
		SetQty(2).
		SetBundle(bundle).
		SetProduct(product1).
		Save(ctx)
	require.NoError(t, err)

	item2, err := client.Bundle_item.Create().
		SetID(uuid.New()).
		SetQty(3).
		SetBundle(bundle).
		SetProduct(product2).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	items, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// Verify items are returned
	itemMap := make(map[uuid.UUID]*GetBundle_itemDTO)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	assert.Contains(t, itemMap, item1.ID)
	assert.Contains(t, itemMap, item2.ID)
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	bundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Test Bundle").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetName("Test Product").
		SetSku("SKU-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test bundle item
	createDTO := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		BundleID:  bundle.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing bundle item
	foundItem, err := repo.FindByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, foundItem.ID)
	assert.Equal(t, item.Qty, foundItem.Qty)
	assert.Equal(t, item.BundleID, foundItem.BundleID)
	assert.Equal(t, item.ProductID, foundItem.ProductID)

	// Test bundle item not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestRepository_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	bundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Test Bundle").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetName("Test Product").
		SetSku("SKU-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Test creating new bundle item
	createDTO := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		BundleID:  bundle.ID,
		ProductID: product.ID,
	}
	createdItem, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.NotNil(t, createdItem)
	assert.Equal(t, createDTO.ID, createdItem.ID)
	assert.Equal(t, createDTO.Qty, createdItem.Qty)
	assert.Equal(t, createDTO.BundleID, createdItem.BundleID)
	assert.Equal(t, createDTO.ProductID, createdItem.ProductID)

	// Test Create - error: bundle not found
	nonExistentBundleID := uuid.New()
	createDTO2 := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		BundleID:  nonExistentBundleID,
		ProductID: product.ID,
	}
	_, err = repo.Create(ctx, createDTO2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - error: product not found
	nonExistentProductID := uuid.New()
	createDTO3 := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       3,
		BundleID:  bundle.ID,
		ProductID: nonExistentProductID,
	}
	_, err = repo.Create(ctx, createDTO3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - with zero Qty
	createDTO4 := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       0,
		BundleID:  bundle.ID,
		ProductID: product.ID,
	}
	createdItem4, err := repo.Create(ctx, createDTO4)
	require.NoError(t, err)
	assert.Equal(t, 0, createdItem4.Qty)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	bundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Test Bundle").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetName("Test Product").
		SetSku("SKU-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test bundle item
	createDTO := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		BundleID:  bundle.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test updating bundle item - update Qty
	newQty := 5
	updateDTO := &UpdateBundle_itemDTO{
		ID:  item.ID,
		Qty: &newQty,
	}
	updatedItem, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.NotNil(t, updatedItem)
	assert.Equal(t, *updateDTO.Qty, updatedItem.Qty)

	// Test no fields to update
	noUpdateDTO := &UpdateBundle_itemDTO{ID: item.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())

	// Test Update - error: bundle_item not found
	nonExistentID := uuid.New()
	newQty2 := 5
	updateDTO2 := &UpdateBundle_itemDTO{
		ID:  nonExistentID,
		Qty: &newQty2,
	}
	_, err = repo.Update(ctx, updateDTO2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Update - with zero Qty
	zeroQty := 0
	updateDTO3 := &UpdateBundle_itemDTO{
		ID:  item.ID,
		Qty: &zeroQty,
	}
	updatedItem3, err := repo.Update(ctx, updateDTO3)
	require.NoError(t, err)
	assert.Equal(t, 0, updatedItem3.Qty)
}

func TestRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	bundle, err := client.Bundle.Create().
		SetID(uuid.New()).
		SetName("Test Bundle").
		SetPrice(99.99).
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetName("Test Product").
		SetSku("SKU-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test bundle item
	createDTO := &CreateBundle_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		BundleID:  bundle.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test deleting bundle item
	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)

	// Verify item is deleted
	_, err = repo.FindByID(ctx, item.ID)
	assert.Error(t, err)
}

