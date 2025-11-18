package meal_plan_items

import (
	"context"
	"testing"
	"time"

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

	mealPlan, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

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

	// Create meal plan items
	day1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	item1, err := client.Meal_plan_item.Create().
		SetID(uuid.New()).
		SetDay(day1).
		SetSlot("breakfast").
		SetMealPlan(mealPlan).
		SetRecipe(recipe1).
		Save(ctx)
	require.NoError(t, err)

	day2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	item2, err := client.Meal_plan_item.Create().
		SetID(uuid.New()).
		SetDay(day2).
		SetSlot("lunch").
		SetMealPlan(mealPlan).
		SetRecipe(recipe2).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	items, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// Verify items are returned
	itemMap := make(map[uuid.UUID]*GetMeal_plan_itemDTO)
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

	mealPlan, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Create test meal plan item
	day := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createDTO := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day,
		Slot:       "breakfast",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test finding existing meal plan item
	foundItem, err := repo.FindByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, foundItem.ID)
	assert.Equal(t, item.Day, foundItem.Day)
	assert.Equal(t, item.Slot, foundItem.Slot)
	assert.Equal(t, item.MealPlanID, foundItem.MealPlanID)
	assert.Equal(t, item.RecipeID, foundItem.RecipeID)

	// Test meal plan item not found
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

	mealPlan, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Test creating new meal plan item
	day := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createDTO := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day,
		Slot:       "breakfast",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	createdItem, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.NotNil(t, createdItem)
	assert.Equal(t, createDTO.ID, createdItem.ID)
	assert.Equal(t, createDTO.Day, createdItem.Day)
	assert.Equal(t, createDTO.Slot, createdItem.Slot)
	assert.Equal(t, createDTO.MealPlanID, createdItem.MealPlanID)
	assert.Equal(t, createDTO.RecipeID, createdItem.RecipeID)

	// Test Create - error: meal plan not found
	nonExistentMealPlanID := uuid.New()
	day2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	createDTO2 := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day2,
		Slot:       "lunch",
		MealPlanID: nonExistentMealPlanID,
		RecipeID:   recipe.ID,
	}
	_, err = repo.Create(ctx, createDTO2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test Create - error: recipe not found
	nonExistentRecipeID := uuid.New()
	day3 := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	createDTO3 := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day3,
		Slot:       "dinner",
		MealPlanID: mealPlan.ID,
		RecipeID:   nonExistentRecipeID,
	}
	_, err = repo.Create(ctx, createDTO3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
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

	mealPlan, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Create test meal plan item
	day := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createDTO := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day,
		Slot:       "breakfast",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test updating meal plan item - update Day and Slot
	newDay := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	newSlot := "lunch"
	updateDTO := &UpdateMeal_plan_itemDTO{
		ID:   item.ID,
		Day:  &newDay,
		Slot: &newSlot,
	}
	updatedItem, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.NotNil(t, updatedItem)
	assert.Equal(t, *updateDTO.Day, updatedItem.Day)
	assert.Equal(t, *updateDTO.Slot, updatedItem.Slot)

	// Test updating meal plan item - update MealPlanID (create new item for this test)
	day2 := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	createDTO2 := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day2,
		Slot:       "dinner",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item2, err := repo.Create(ctx, createDTO2)
	require.NoError(t, err)

	newMealPlan, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	newDay2 := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)
	updateDTO2 := &UpdateMeal_plan_itemDTO{
		ID:         item2.ID,
		MealPlanID: &newMealPlan.ID,
		Day:        &newDay2, // Also update a field to ensure mutation has fields
	}
	updatedItem2, err := repo.Update(ctx, updateDTO2)
	require.NoError(t, err)
	assert.Equal(t, newMealPlan.ID, updatedItem2.MealPlanID)
	assert.Equal(t, newDay2, updatedItem2.Day)

	// Test updating meal plan item - update RecipeID (create new item for this test)
	day3 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	createDTO3 := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day3,
		Slot:       "breakfast",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item3, err := repo.Create(ctx, createDTO3)
	require.NoError(t, err)

	instructions2 := "New recipe instructions"
	newRecipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("New Recipe").
		SetNillableInstructions(&instructions2).
		SetKcal(600).
		Save(ctx)
	require.NoError(t, err)

	newDay3 := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	updateDTO3 := &UpdateMeal_plan_itemDTO{
		ID:       item3.ID,
		RecipeID: &newRecipe.ID,
		Day:      &newDay3, // Also update a field to ensure mutation has fields
	}
	updatedItem3, err := repo.Update(ctx, updateDTO3)
	require.NoError(t, err)
	assert.Equal(t, newRecipe.ID, updatedItem3.RecipeID)
	assert.Equal(t, newDay3, updatedItem3.Day)

	// Test no fields to update
	noUpdateDTO := &UpdateMeal_plan_itemDTO{ID: item.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())

	// Test meal plan not found (create new item for this test)
	day4 := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)
	createDTO4 := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day4,
		Slot:       "lunch",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item4, err := repo.Create(ctx, createDTO4)
	require.NoError(t, err)

	nonExistentMealPlanID := uuid.New()
	newDay4 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	updateDTO4 := &UpdateMeal_plan_itemDTO{
		ID:         item4.ID,
		MealPlanID: &nonExistentMealPlanID,
		Day:        &newDay4, // Also update a field to ensure mutation has fields
	}
	_, err = repo.Update(ctx, updateDTO4)
	assert.Error(t, err)

	// Test recipe not found (create new item for this test)
	day5 := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)
	createDTO5 := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day5,
		Slot:       "dinner",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item5, err := repo.Create(ctx, createDTO5)
	require.NoError(t, err)

	nonExistentRecipeID := uuid.New()
	newDay5 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	updateDTO5 := &UpdateMeal_plan_itemDTO{
		ID:       item5.ID,
		RecipeID: &nonExistentRecipeID,
		Day:      &newDay5, // Also update a field to ensure mutation has fields
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

	mealPlan, err := client.Meal_plan.Create().
		SetID(uuid.New()).
		SetWeekStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	instructions := "Test instructions"
	recipe, err := client.Recipe.Create().
		SetID(uuid.New()).
		SetName("Test Recipe").
		SetNillableInstructions(&instructions).
		SetKcal(500).
		Save(ctx)
	require.NoError(t, err)

	// Create test meal plan item
	day := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createDTO := &CreateMeal_plan_itemDTO{
		ID:         uuid.New(),
		Day:        day,
		Slot:       "breakfast",
		MealPlanID: mealPlan.ID,
		RecipeID:   recipe.ID,
	}
	item, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)

	// Test deleting meal plan item
	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)

	// Verify item is deleted
	_, err = repo.FindByID(ctx, item.ID)
	assert.Error(t, err)
}

