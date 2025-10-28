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

	// Create test product categories
	category1, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits and vegetables").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	category2, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Vegetables").
		SetDescription("Fresh vegetables").
		SetSlug("vegetables").
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created categories
	foundIDs := make(map[uuid.UUID]bool)
	for _, cat := range result {
		foundIDs[cat.ID] = true
		assert.NotEmpty(t, cat.Name)
		assert.NotEmpty(t, cat.Description)
		assert.NotEmpty(t, cat.Slug)
	}

	assert.True(t, foundIDs[category1.ID])
	assert.True(t, foundIDs[category2.ID])
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test product category
	createdCategory, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits and vegetables").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdCategory.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdCategory.ID, result.ID)
	assert.Equal(t, "Fruits", result.Name)
	assert.Equal(t, "Fresh fruits and vegetables", result.Description)
	assert.Equal(t, "fruits", result.Slug)

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

	dto := &CreateProductCategoryDTO{
		ID:          uuid.New(),
		Name:        "Fruits",
		Description: "Fresh fruits and vegetables",
		Slug:        "fruits",
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID) // ID should be generated
	assert.Equal(t, dto.Name, result.Name)
	assert.Equal(t, dto.Description, result.Description)
	assert.Equal(t, dto.Slug, result.Slug)

	// Verify it was actually created in the database
	dbCategory, err := client.Product_category.Get(ctx, result.ID)
	require.NoError(t, err)
	assert.Equal(t, dto.Name, dbCategory.Name)
	assert.Equal(t, dto.Description, dbCategory.Description)
	assert.Equal(t, dto.Slug, dbCategory.Slug)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test product category
	createdCategory, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits and vegetables").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	// Test Update - full update
	dto := &UpdateProductCategoryDTO{
		ID:          createdCategory.ID,
		Name:        stringPtr("Updated Fruits"),
		Description: stringPtr("Updated description"),
		Slug:        stringPtr("updated-fruits"),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdCategory.ID, result.ID)
	assert.Equal(t, "Updated Fruits", result.Name)
	assert.Equal(t, "Updated description", result.Description)
	assert.Equal(t, "updated-fruits", result.Slug)

	// Verify it was actually updated in the database
	dbCategory, err := client.Product_category.Get(ctx, createdCategory.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Fruits", dbCategory.Name)
	assert.Equal(t, "Updated description", dbCategory.Description)
	assert.Equal(t, "updated-fruits", dbCategory.Slug)

	// Test Update - partial update (only name)
	dto2 := &UpdateProductCategoryDTO{
		ID:   createdCategory.ID,
		Name: stringPtr("Partial Update"),
	}

	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, createdCategory.ID, result2.ID)
	assert.Equal(t, "Partial Update", result2.Name)
	assert.Equal(t, "Updated description", result2.Description) // Should remain unchanged
	assert.Equal(t, "updated-fruits", result2.Slug)             // Should remain unchanged

	// Test Update - no fields to update
	dto3 := &UpdateProductCategoryDTO{
		ID: createdCategory.ID,
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

	// Create test product category
	createdCategory, err := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits and vegetables").
		SetSlug("fruits").
		Save(ctx)
	require.NoError(t, err)

	// Test Delete - success
	err = repo.Delete(ctx, createdCategory.ID)
	require.NoError(t, err)

	// Verify it was actually deleted from the database
	_, err = client.Product_category.Get(ctx, createdCategory.ID)
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

	// Create multiple product categories
	dto1 := &CreateProductCategoryDTO{
		ID:          uuid.New(),
		Name:        "Fruits",
		Description: "Fresh fruits and vegetables",
		Slug:        "fruits",
	}

	dto2 := &CreateProductCategoryDTO{
		ID:          uuid.New(),
		Name:        "Vegetables",
		Description: "Fresh vegetables",
		Slug:        "vegetables",
	}

	category1, err := repo.Create(ctx, dto1)
	require.NoError(t, err)

	category2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)

	// List all categories
	allCategories, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, allCategories, 2)

	// Get specific category
	retrievedCategory, err := repo.FindByID(ctx, category1.ID)
	require.NoError(t, err)
	assert.Equal(t, category1.ID, retrievedCategory.ID)
	assert.Equal(t, dto1.Name, retrievedCategory.Name)

	// Update category
	updateDTO := &UpdateProductCategoryDTO{
		ID:   category1.ID,
		Name: stringPtr("Updated Fruits"),
	}

	updatedCategory, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, "Updated Fruits", updatedCategory.Name)
	assert.Equal(t, dto1.Description, updatedCategory.Description)
	assert.Equal(t, dto1.Slug, updatedCategory.Slug)

	// Delete one category
	err = repo.Delete(ctx, category1.ID)
	require.NoError(t, err)

	// Verify only one category remains
	remainingCategories, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, remainingCategories, 1)
	assert.Equal(t, category2.ID, remainingCategories[0].ID)
}
