package order_items

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
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
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

	// Create order items
	item1, err := client.Order_item.Create().
		SetID(uuid.New()).
		SetQty(2).
		SetUnitPrice(10.99).
		SetLineTotal(21.98).
		SetOrder(order).
		SetProduct(product1).
		Save(ctx)
	require.NoError(t, err)

	item2, err := client.Order_item.Create().
		SetID(uuid.New()).
		SetQty(3).
		SetUnitPrice(5.50).
		SetLineTotal(16.50).
		SetOrder(order).
		SetProduct(product2).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	items, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// Verify items are returned
	itemMap := make(map[uuid.UUID]*GetOrder_itemDTO)
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
	user, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
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

	// Create test order item
	createDTO := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		UnitPrice: 10.99,
		LineTotal: 21.98,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing order item
	foundItem, err := repo.FindByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, foundItem.ID)
	assert.Equal(t, item.Qty, foundItem.Qty)
	assert.Equal(t, item.UnitPrice, foundItem.UnitPrice)
	assert.Equal(t, item.LineTotal, foundItem.LineTotal)
	assert.Equal(t, item.OrderID, foundItem.OrderID)
	assert.Equal(t, item.ProductID, foundItem.ProductID)

	// Test order item not found
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

	order, err := client.Order.Create().
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
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

	// Test creating new order item
	createDTO := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		UnitPrice: 10.99,
		LineTotal: 21.98,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	createdItem, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.NotNil(t, createdItem)
	assert.Equal(t, createDTO.ID, createdItem.ID)
	assert.Equal(t, createDTO.Qty, createdItem.Qty)
	assert.Equal(t, createDTO.UnitPrice, createdItem.UnitPrice)
	assert.Equal(t, createDTO.LineTotal, createdItem.LineTotal)
	assert.Equal(t, createDTO.OrderID, createdItem.OrderID)
	assert.Equal(t, createDTO.ProductID, createdItem.ProductID)

	// Test Create - error: order not found
	nonExistentOrderID := uuid.New()
	createDTO2 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       1,
		UnitPrice: 5.0,
		LineTotal: 5.0,
		OrderID:   nonExistentOrderID,
		ProductID: product.ID,
	}
	_, err = repo.Create(ctx, createDTO2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - error: product not found
	nonExistentProductID := uuid.New()
	createDTO3 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       1,
		UnitPrice: 5.0,
		LineTotal: 5.0,
		OrderID:   order.ID,
		ProductID: nonExistentProductID,
	}
	_, err = repo.Create(ctx, createDTO3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - with zero values
	createDTO4 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       0,
		UnitPrice: 0.0,
		LineTotal: 0.0,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	createdItem4, err := repo.Create(ctx, createDTO4)
	require.NoError(t, err)
	assert.Equal(t, 0, createdItem4.Qty)
	assert.Equal(t, 0.0, createdItem4.UnitPrice)
	assert.Equal(t, 0.0, createdItem4.LineTotal)
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

	order, err := client.Order.Create().
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
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

	// Create test order item
	createDTO := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		UnitPrice: 10.99,
		LineTotal: 21.98,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test updating order item - update Qty
	updateDTO := &UpdateOrder_itemDTO{
		ID:        item.ID,
		Qty:       intPtr(5),
		UnitPrice: float64Ptr(12.50),
		LineTotal: float64Ptr(62.50),
	}
	updatedItem, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.NotNil(t, updatedItem)
	assert.Equal(t, *updateDTO.Qty, updatedItem.Qty)
	assert.Equal(t, *updateDTO.UnitPrice, updatedItem.UnitPrice)
	assert.Equal(t, *updateDTO.LineTotal, updatedItem.LineTotal)

	// Test updating order item - update OrderID (create new item for this test)
	createDTO2 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       1,
		UnitPrice: 5.0,
		LineTotal: 5.0,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item2, err := repo.Create(ctx, createDTO2)
	require.NoError(t, err)

	newOrder, err := client.Order.Create().
		SetOrderNo("ORD-002").
		SetStatus("processing").
		SetSubtotal(200.0).
		SetShippingFee(20.0).
		SetDiscount(0.0).
		SetTotal(220.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	updateDTO2 := &UpdateOrder_itemDTO{
		ID:        item2.ID,
		OrderID:   &newOrder.ID,
		Qty:       intPtr(3), // Also update a field to ensure mutation has fields
	}
	updatedItem2, err := repo.Update(ctx, updateDTO2)
	require.NoError(t, err)
	assert.Equal(t, newOrder.ID, updatedItem2.OrderID)
	assert.Equal(t, 3, updatedItem2.Qty)

	// Test updating order item - update ProductID (create new item for this test)
	createDTO3 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       1,
		UnitPrice: 5.0,
		LineTotal: 5.0,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item3, err := repo.Create(ctx, createDTO3)
	require.NoError(t, err)

	newProduct, err := client.Product.Create().
		SetName("New Product").
		SetSku("SKU-002").
		SetPrice(15.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	updateDTO3 := &UpdateOrder_itemDTO{
		ID:        item3.ID,
		ProductID: &newProduct.ID,
		Qty:       intPtr(4), // Also update a field to ensure mutation has fields
	}
	updatedItem3, err := repo.Update(ctx, updateDTO3)
	require.NoError(t, err)
	assert.Equal(t, newProduct.ID, updatedItem3.ProductID)
	assert.Equal(t, 4, updatedItem3.Qty)

	// Test no fields to update
	noUpdateDTO := &UpdateOrder_itemDTO{ID: item.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())

	// Test order not found (create new item for this test)
	createDTO4 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       1,
		UnitPrice: 5.0,
		LineTotal: 5.0,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item4, err := repo.Create(ctx, createDTO4)
	require.NoError(t, err)

	nonExistentOrderID := uuid.New()
	updateDTO4 := &UpdateOrder_itemDTO{
		ID:        item4.ID,
		OrderID:   &nonExistentOrderID,
		Qty:       intPtr(2), // Also update a field to ensure mutation has fields
	}
	_, err = repo.Update(ctx, updateDTO4)
	assert.Error(t, err)

	// Test product not found (create new item for this test)
	createDTO5 := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       1,
		UnitPrice: 5.0,
		LineTotal: 5.0,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item5, err := repo.Create(ctx, createDTO5)
	require.NoError(t, err)

	nonExistentProductID := uuid.New()
	updateDTO5 := &UpdateOrder_itemDTO{
		ID:        item5.ID,
		ProductID: &nonExistentProductID,
		Qty:       intPtr(2), // Also update a field to ensure mutation has fields
	}
	_, err = repo.Update(ctx, updateDTO5)
	assert.Error(t, err)
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

	order, err := client.Order.Create().
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
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

	// Create test order item
	createDTO := &CreateOrder_itemDTO{
		ID:        uuid.New(),
		Qty:       2,
		UnitPrice: 10.99,
		LineTotal: 21.98,
		OrderID:   order.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test deleting order item
	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)

	// Verify item is deleted
	_, err = repo.FindByID(ctx, item.ID)
	assert.Error(t, err)
}

