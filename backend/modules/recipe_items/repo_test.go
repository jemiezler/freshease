package recipe_items

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
	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
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

	// Create recipe items
	item1, err := client.Recipe_item.Create().
		SetID(uuid.New()).
		SetAmount(2.5).
		SetUnit("kg").
		SetRecipe(recipe).
		SetProduct(product1).
		Save(ctx)
	require.NoError(t, err)

	item2, err := client.Recipe_item.Create().
		SetID(uuid.New()).
		SetAmount(1.0).
		SetUnit("cup").
		SetRecipe(recipe).
		SetProduct(product2).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	items, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// Verify items are returned
	itemMap := make(map[uuid.UUID]*GetRecipe_itemDTO)
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
	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
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

	// Create test recipe item
	createDTO := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    2.5,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing recipe item
	foundItem, err := repo.FindByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, foundItem.ID)
	assert.Equal(t, item.Amount, foundItem.Amount)
	assert.Equal(t, item.Unit, foundItem.Unit)
	assert.Equal(t, item.RecipeID, foundItem.RecipeID)
	assert.Equal(t, item.ProductID, foundItem.ProductID)

	// Test recipe item not found
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
	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
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

	// Test creating new recipe item
	createDTO := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    2.5,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	createdItem, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.NotNil(t, createdItem)
	assert.Equal(t, createDTO.ID, createdItem.ID)
	assert.Equal(t, createDTO.Amount, createdItem.Amount)
	assert.Equal(t, createDTO.Unit, createdItem.Unit)
	assert.Equal(t, createDTO.RecipeID, createdItem.RecipeID)
	assert.Equal(t, createDTO.ProductID, createdItem.ProductID)

	// Test Create - error: recipe not found
	nonExistentRecipeID := uuid.New()
	createDTO2 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    1.0,
		Unit:      "cup",
		RecipeID:  nonExistentRecipeID,
		ProductID: product.ID,
	}
	_, err = repo.Create(ctx, createDTO2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - error: product not found
	nonExistentProductID := uuid.New()
	createDTO3 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    1.0,
		Unit:      "cup",
		RecipeID:  recipe.ID,
		ProductID: nonExistentProductID,
	}
	_, err = repo.Create(ctx, createDTO3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - with zero values
	createDTO4 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    0.0,
		Unit:      "",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	createdItem4, err := repo.Create(ctx, createDTO4)
	require.NoError(t, err)
	assert.Equal(t, 0.0, createdItem4.Amount)
	assert.Equal(t, "", createdItem4.Unit)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
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

	// Create test recipe item
	createDTO := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    2.5,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test updating recipe item - update Amount and Unit
	updateDTO := &UpdateRecipe_itemDTO{
		ID:     item.ID,
		Amount: float64Ptr(5.0),
		Unit:   stringPtr("cup"),
	}
	updatedItem, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.NotNil(t, updatedItem)
	assert.Equal(t, *updateDTO.Amount, updatedItem.Amount)
	assert.Equal(t, *updateDTO.Unit, updatedItem.Unit)

	// Test updating recipe item - update RecipeID (create new item for this test)
	createDTO2 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    1.0,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	item2, err := repo.Create(ctx, createDTO2)
	require.NoError(t, err)

	newRecipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("New Recipe").
		SetKcal(600).
		Save(ctx)
	require.NoError(t, err)

	updateDTO2 := &UpdateRecipe_itemDTO{
		ID:       item2.ID,
		RecipeID: &newRecipe.ID,
		Amount:  float64Ptr(2.0), // Also update a field to ensure mutation has fields
	}
	updatedItem2, err := repo.Update(ctx, updateDTO2)
	require.NoError(t, err)
	assert.Equal(t, newRecipe.ID, updatedItem2.RecipeID)
	assert.Equal(t, 2.0, updatedItem2.Amount)

	// Test updating recipe item - update ProductID (create new item for this test)
	createDTO3 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    1.0,
		Unit:      "kg",
		RecipeID:  recipe.ID,
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

	updateDTO3 := &UpdateRecipe_itemDTO{
		ID:        item3.ID,
		ProductID: &newProduct.ID,
		Amount:    float64Ptr(3.0), // Also update a field to ensure mutation has fields
	}
	updatedItem3, err := repo.Update(ctx, updateDTO3)
	require.NoError(t, err)
	assert.Equal(t, newProduct.ID, updatedItem3.ProductID)
	assert.Equal(t, 3.0, updatedItem3.Amount)

	// Test no fields to update
	noUpdateDTO := &UpdateRecipe_itemDTO{ID: item.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())

	// Test recipe not found (create new item for this test)
	createDTO4 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    1.0,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	item4, err := repo.Create(ctx, createDTO4)
	require.NoError(t, err)

	nonExistentRecipeID := uuid.New()
	updateDTO4 := &UpdateRecipe_itemDTO{
		ID:       item4.ID,
		RecipeID: &nonExistentRecipeID,
		Amount:   float64Ptr(2.0), // Also update a field to ensure mutation has fields
	}
	_, err = repo.Update(ctx, updateDTO4)
	assert.Error(t, err)

	// Test product not found (create new item for this test)
	createDTO5 := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    1.0,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	item5, err := repo.Create(ctx, createDTO5)
	require.NoError(t, err)

	nonExistentProductID := uuid.New()
	updateDTO5 := &UpdateRecipe_itemDTO{
		ID:        item5.ID,
		ProductID: &nonExistentProductID,
		Amount:    float64Ptr(2.0), // Also update a field to ensure mutation has fields
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
	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
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

	// Create test recipe item
	createDTO := &CreateRecipe_itemDTO{
		ID:        uuid.New(),
		Amount:    2.5,
		Unit:      "kg",
		RecipeID:  recipe.ID,
		ProductID: product.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test deleting recipe item
	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)

	// Verify item is deleted
	_, err = repo.FindByID(ctx, item.ID)
	assert.Error(t, err)
}

