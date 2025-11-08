package product_categories

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

	// Create test products and categories first
	product1, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetSku("APPLE-001").
		SetPrice(2.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	product2, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Banana").
		SetSku("BANANA-001").
		SetPrice(1.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	category1, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	category2, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Vegetables").
		SetSlug("vegetables").
		Save(ctx)
	require.NoError(t, err)

	// Create product_category joins
	pc1, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetProduct(product1).
		SetCategory(category1).
		Save(ctx)
	require.NoError(t, err)

	pc2, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetProduct(product2).
		SetCategory(category2).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created product categories
	foundIDs := make(map[uuid.UUID]bool)
	for _, pc := range result {
		foundIDs[pc.ID] = true
		assert.NotEqual(t, uuid.Nil, pc.ProductID)
		assert.NotEqual(t, uuid.Nil, pc.CategoryID)
	}

	assert.True(t, foundIDs[pc1.ID])
	assert.True(t, foundIDs[pc2.ID])
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test product and category first
	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetSku("APPLE-001").
		SetPrice(2.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	category, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	// Create product_category join
	createdPC, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetProduct(product).
		SetCategory(category).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdPC.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdPC.ID, result.ID)
	assert.Equal(t, product.ID, result.ProductID)
	assert.Equal(t, category.ID, result.CategoryID)

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

	// Create test product and category first
	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetSku("APPLE-001").
		SetPrice(2.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	category, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	dto := &CreateProductCategoryDTO{
		ID:         uuid.New(),
		ProductID:  product.ID,
		CategoryID: category.ID,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.ProductID, result.ProductID)
	assert.Equal(t, dto.CategoryID, result.CategoryID)

	// Verify it was actually created in the database
	dbPC, err := client.Product_category.Get(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, dto.ID, dbPC.ID)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test products and categories first
	product1, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetSku("APPLE-001").
		SetPrice(2.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	product2, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Banana").
		SetSku("BANANA-001").
		SetPrice(1.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	category1, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	category2, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Vegetables").
		SetSlug("vegetables").
		Save(ctx)
	require.NoError(t, err)

	// Create product_category join
	createdPC, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetProduct(product1).
		SetCategory(category1).
		Save(ctx)
	require.NoError(t, err)

	// Test Update - update category
	dto := &UpdateProductCategoryDTO{
		ID:         createdPC.ID,
		CategoryID: func() *uuid.UUID { id := category2.ID; return &id }(),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdPC.ID, result.ID)
	assert.Equal(t, product1.ID, result.ProductID)   // Product should remain unchanged
	assert.Equal(t, category2.ID, result.CategoryID) // Category should be updated

	// Test Update - update product
	dto2 := &UpdateProductCategoryDTO{
		ID:        createdPC.ID,
		ProductID: func() *uuid.UUID { id := product2.ID; return &id }(),
	}

	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, createdPC.ID, result2.ID)
	assert.Equal(t, product2.ID, result2.ProductID)   // Product should be updated
	assert.Equal(t, category2.ID, result2.CategoryID) // Category should remain unchanged

	// Test Update - no fields to update
	dto3 := &UpdateProductCategoryDTO{
		ID: createdPC.ID,
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

	// Create test product and category first
	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetSku("APPLE-001").
		SetPrice(2.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	category, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	// Create product_category join
	createdPC, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetProduct(product).
		SetCategory(category).
		Save(ctx)
	require.NoError(t, err)

	// Test Delete - success
	err = repo.Delete(ctx, createdPC.ID)
	require.NoError(t, err)

	// Verify it was actually deleted from the database
	_, err = client.Product_category.Get(ctx, createdPC.ID)
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

	// Create test products and categories first
	product1, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetSku("APPLE-001").
		SetPrice(2.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	product2, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Banana").
		SetSku("BANANA-001").
		SetPrice(1.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		Save(ctx)
	require.NoError(t, err)

	category1, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	category2, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Vegetables").
		SetSlug("vegetables").
		Save(ctx)
	require.NoError(t, err)

	// Create multiple product categories
	dto1 := &CreateProductCategoryDTO{
		ID:         uuid.New(),
		ProductID:  product1.ID,
		CategoryID: category1.ID,
	}

	dto2 := &CreateProductCategoryDTO{
		ID:         uuid.New(),
		ProductID:  product2.ID,
		CategoryID: category2.ID,
	}

	pc1, err := repo.Create(ctx, dto1)
	require.NoError(t, err)

	pc2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)

	// List all product categories
	allPCs, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allPCs, 2)

	// Get specific product category
	retrievedPC, err := repo.FindByID(ctx, pc1.ID)
	require.NoError(t, err)
	assert.Equal(t, pc1.ID, retrievedPC.ID)
	assert.Equal(t, dto1.ProductID, retrievedPC.ProductID)
	assert.Equal(t, dto1.CategoryID, retrievedPC.CategoryID)

	// Update product category
	updateDTO := &UpdateProductCategoryDTO{
		ID:        pc1.ID,
		ProductID: func() *uuid.UUID { id := product2.ID; return &id }(),
	}

	updatedPC, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, product2.ID, updatedPC.ProductID)
	assert.Equal(t, dto1.CategoryID, updatedPC.CategoryID) // Category should remain unchanged

	// Delete one product category
	err = repo.Delete(ctx, pc1.ID)
	require.NoError(t, err)

	// Verify only one product category remains
	remainingPCs, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, remainingPCs, 1)
	assert.Equal(t, pc2.ID, remainingPCs[0].ID)
}
