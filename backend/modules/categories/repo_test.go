package categories

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

func TestEntRepo_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test categories
	category1, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Category One").
		SetSlug("category-one").
		Save(ctx)
	require.NoError(t, err)

	category2, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Category Two").
		SetSlug("category-two").
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results contain our created categories
	foundIDs := make(map[uuid.UUID]bool)
	for _, category := range result {
		foundIDs[category.ID] = true
		assert.NotEmpty(t, category.Name)
		assert.NotEmpty(t, category.Slug)
	}

	assert.True(t, foundIDs[category1.ID])
	assert.True(t, foundIDs[category2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test category
	createdCategory, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Test Category").
		SetSlug("test-category").
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdCategory.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdCategory.ID, result.ID)
	assert.Equal(t, "Test Category", result.Name)
	assert.Equal(t, "test-category", result.Slug)

	// Test FindByID - not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestEntRepo_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	now := time.Now()
	dto := &CreateCategoryDTO{
		ID:        uuid.New(),
		Name:      "New Category",
		Slug:      "new-category",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Name, result.Name)
	assert.Equal(t, dto.Slug, result.Slug)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test category
	createdCategory, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("Original Category").
		SetSlug("original-category").
		Save(ctx)
	require.NoError(t, err)

	// Update category
	newName := "Updated Category"
	newSlug := "updated-category"
	dto := &UpdateCategoryDTO{
		ID:        createdCategory.ID,
		Name:      &newName,
		Slug:      &newSlug,
		UpdatedAt: time.Now(),
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdCategory.ID, result.ID)
	assert.Equal(t, "Updated Category", result.Name)
	assert.Equal(t, "updated-category", result.Slug)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test category
	createdCategory, err := client.Category.Create().
		SetID(uuid.New()).
		SetName("To Delete").
		SetSlug("to-delete").
		Save(ctx)
	require.NoError(t, err)

	// Delete category
	err = repo.Delete(ctx, createdCategory.ID)
	require.NoError(t, err)

	// Verify category is deleted
	_, err = repo.FindByID(ctx, createdCategory.ID)
	assert.Error(t, err)
}

