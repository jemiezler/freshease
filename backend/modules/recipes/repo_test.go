package recipes

import (
	"context"
	"testing"

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

	// Create test recipes
	instructions1 := "Cook for 30 minutes"
	recipe1, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Recipe One").
		SetNillableInstructions(&instructions1).
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	instructions2 := "Bake for 45 minutes"
	recipe2, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Recipe Two").
		SetNillableInstructions(&instructions2).
		SetKcal(600).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results
	foundIDs := make(map[uuid.UUID]bool)
	for _, recipe := range result {
		foundIDs[recipe.ID] = true
		assert.NotEmpty(t, recipe.Name)
		assert.Greater(t, recipe.Kcal, 0)
	}

	assert.True(t, foundIDs[recipe1.ID])
	assert.True(t, foundIDs[recipe2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test recipe
	instructions := "Test instructions"
	createdRecipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdRecipe.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdRecipe.ID, result.ID)
	assert.Equal(t, "Test Recipe", result.Name)
	assert.Equal(t, 500, result.Kcal)
	assert.NotNil(t, result.Instructions)
	assert.Equal(t, instructions, *result.Instructions)

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

	instructions := "New recipe instructions"
	dto := &CreateRecipeDTO{
		ID:           uuid.New(),
		Name:         "New Recipe",
		Instructions: &instructions,
		Kcal:         550,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Name, result.Name)
	assert.Equal(t, dto.Kcal, result.Kcal)
	assert.NotNil(t, result.Instructions)
	assert.Equal(t, instructions, *result.Instructions)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test recipe
	createdRecipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Original Recipe").
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Update recipe
	newName := "Updated Recipe"
	newKcal := 600
	newInstructions := "Updated instructions"
	dto := &UpdateRecipeDTO{
		ID:           createdRecipe.ID,
		Name:         &newName,
		Kcal:         &newKcal,
		Instructions: &newInstructions,
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdRecipe.ID, result.ID)
	assert.Equal(t, "Updated Recipe", result.Name)
	assert.Equal(t, 600, result.Kcal)
	assert.NotNil(t, result.Instructions)
	assert.Equal(t, newInstructions, *result.Instructions)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test recipe
	createdRecipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("To Delete").
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Delete recipe
	err = repo.Delete(ctx, createdRecipe.ID)
	require.NoError(t, err)

	// Verify recipe is deleted
	_, err = repo.FindByID(ctx, createdRecipe.ID)
	assert.Error(t, err)
}

