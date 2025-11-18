package cart_items

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

	// Test empty list
	items, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, items)

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	cart, err := client.Cart.Create().
		SetStatus("pending").
		AddUser(user).
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

	// Create cart items
	item1, err := client.Cart_item.Create().
		SetID(uuid.New()).
		SetQty(2).
		SetUnitPrice(10.99).
		SetLineTotal(21.98).
		SetCart(cart).
		SetProduct(product1).
		Save(ctx)
	require.NoError(t, err)

	item2, err := client.Cart_item.Create().
		SetID(uuid.New()).
		SetQty(3).
		SetUnitPrice(5.50).
		SetLineTotal(16.50).
		SetCart(cart).
		SetProduct(product2).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	items, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// Verify items are returned
	itemMap := make(map[uuid.UUID]*GetCart_itemDTO)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	assert.Contains(t, itemMap, item1.ID)
	assert.Contains(t, itemMap, item2.ID)

	// Verify first item details
	foundItem1 := itemMap[item1.ID]
	assert.Equal(t, item1.Qty, foundItem1.Qty)
	assert.Equal(t, item1.UnitPrice, foundItem1.UnitPrice)
	assert.Equal(t, item1.LineTotal, foundItem1.LineTotal)
	assert.Equal(t, cart.ID, foundItem1.CartID)
	assert.Equal(t, product1.ID, foundItem1.ProductID)
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	cart, err := client.Cart.Create().
		SetStatus("pending").
		AddUser(user).
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

	// Create cart item
	item, err := client.Cart_item.Create().
		SetID(uuid.New()).
		SetQty(2).
		SetUnitPrice(10.99).
		SetLineTotal(21.98).
		SetCart(cart).
		SetProduct(product).
		Save(ctx)
	require.NoError(t, err)

	// Test finding existing item
	foundItem, err := repo.FindByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, foundItem.ID)
	assert.Equal(t, item.Qty, foundItem.Qty)
	assert.Equal(t, item.UnitPrice, foundItem.UnitPrice)
	assert.Equal(t, item.LineTotal, foundItem.LineTotal)
	assert.Equal(t, cart.ID, foundItem.CartID)
	assert.Equal(t, product.ID, foundItem.ProductID)

	// Test finding non-existing item
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
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	cart, err := client.Cart.Create().
		SetStatus("pending").
		AddUser(user).
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

	// Create cart item DTO
	createDTO := &CreateCart_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		UnitPrice: 10.99,
		LineTotal: 21.98,
		CartID:    cart.ID,
		ProductID: product.ID,
	}

	// Test creating cart item
	createdItem, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.Equal(t, createDTO.ID, createdItem.ID)
	assert.Equal(t, createDTO.Qty, createdItem.Qty)
	assert.Equal(t, createDTO.UnitPrice, createdItem.UnitPrice)
	assert.Equal(t, createDTO.LineTotal, createdItem.LineTotal)
	assert.Equal(t, createDTO.CartID, createdItem.CartID)
	assert.Equal(t, createDTO.ProductID, createdItem.ProductID)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	cart, err := client.Cart.Create().
		SetStatus("pending").
		AddUser(user).
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

	// Create cart item
	item, err := client.Cart_item.Create().
		SetID(uuid.New()).
		SetQty(2).
		SetUnitPrice(10.99).
		SetLineTotal(21.98).
		SetCart(cart).
		SetProduct(product1).
		Save(ctx)
	require.NoError(t, err)

	// Update cart item
	newQty := 5
	newUnitPrice := 12.50
	newLineTotal := 62.50
	updateDTO := &UpdateCart_itemDTO{
		ID:        item.ID,
		Qty:       &newQty,
		UnitPrice: &newUnitPrice,
		LineTotal: &newLineTotal,
		ProductID: &product2.ID,
	}

	updatedItem, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, item.ID, updatedItem.ID)
	assert.Equal(t, newQty, updatedItem.Qty)
	assert.Equal(t, newUnitPrice, updatedItem.UnitPrice)
	assert.Equal(t, newLineTotal, updatedItem.LineTotal)
	assert.Equal(t, product2.ID, updatedItem.ProductID)
}

func TestRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	cart, err := client.Cart.Create().
		SetStatus("pending").
		AddUser(user).
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

	// Create cart item
	item, err := client.Cart_item.Create().
		SetID(uuid.New()).
		SetQty(2).
		SetUnitPrice(10.99).
		SetLineTotal(21.98).
		SetCart(cart).
		SetProduct(product).
		Save(ctx)
	require.NoError(t, err)

	// Test deleting cart item
	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)

	// Verify item is deleted
	_, err = repo.FindByID(ctx, item.ID)
	assert.Error(t, err)
}

// Helper functions to create pointers (intPtr and float64Ptr are in controller_test.go)
func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
