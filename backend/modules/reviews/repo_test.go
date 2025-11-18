package reviews

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

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test vendor and product
	vendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Test Product").
		SetSku("TEST-001").
		SetPrice(99.99).
		SetDescription("Test product").
		SetUnitLabel("kg").
		SetIsActive(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test reviews
	comment1 := "Great product!"
	review1, err := client.Review.Create().
		SetID(uuid.New()).
		SetRating(5).
		SetNillableComment(&comment1).
		AddUser(user).
		AddProduct(product).
		Save(ctx)
	require.NoError(t, err)

	comment2 := "Good quality"
	review2, err := client.Review.Create().
		SetID(uuid.New()).
		SetRating(4).
		SetNillableComment(&comment2).
		AddUser(user).
		AddProduct(product).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results
	foundIDs := make(map[uuid.UUID]bool)
	for _, review := range result {
		foundIDs[review.ID] = true
		assert.GreaterOrEqual(t, review.Rating, 1)
		assert.LessOrEqual(t, review.Rating, 5)
		assert.Equal(t, user.ID, review.UserID)
		assert.Equal(t, product.ID, review.ProductID)
	}

	assert.True(t, foundIDs[review1.ID])
	assert.True(t, foundIDs[review2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test vendor and product
	vendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Test Product").
		SetSku("TEST-001").
		SetPrice(99.99).
		SetDescription("Test product").
		SetUnitLabel("kg").
		SetIsActive(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test review
	comment := "Excellent product!"
	createdReview, err := client.Review.Create().
		SetID(uuid.New()).
		SetRating(5).
		SetNillableComment(&comment).
		AddUser(user).
		AddProduct(product).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdReview.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdReview.ID, result.ID)
	assert.Equal(t, 5, result.Rating)
	assert.NotNil(t, result.Comment)
	assert.Equal(t, comment, *result.Comment)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, product.ID, result.ProductID)

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

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test vendor and product
	vendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Test Product").
		SetSku("TEST-001").
		SetPrice(99.99).
		SetDescription("Test product").
		SetUnitLabel("kg").
		SetIsActive(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	comment := "New review comment"
	now := time.Now()
	dto := &CreateReviewDTO{
		ID:        uuid.New(),
		Rating:    5,
		Comment:   &comment,
		UserID:    user.ID,
		ProductID: product.ID,
		CreatedAt: &now,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Rating, result.Rating)
	assert.NotNil(t, result.Comment)
	assert.Equal(t, comment, *result.Comment)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, product.ID, result.ProductID)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test vendor and product
	vendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Test Product").
		SetSku("TEST-001").
		SetPrice(99.99).
		SetDescription("Test product").
		SetUnitLabel("kg").
		SetIsActive(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test review
	createdReview, err := client.Review.Create().
		SetID(uuid.New()).
		SetRating(3).
		AddUser(user).
		AddProduct(product).
		Save(ctx)
	require.NoError(t, err)

	// Update review
	newRating := 5
	newComment := "Updated comment"
	dto := &UpdateReviewDTO{
		ID:      createdReview.ID,
		Rating:  &newRating,
		Comment: &newComment,
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdReview.ID, result.ID)
	assert.Equal(t, 5, result.Rating)
	assert.NotNil(t, result.Comment)
	assert.Equal(t, newComment, *result.Comment)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test vendor and product
	vendor, err := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Test Product").
		SetSku("TEST-001").
		SetPrice(99.99).
		SetDescription("Test product").
		SetUnitLabel("kg").
		SetIsActive(true).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create test review
	createdReview, err := client.Review.Create().
		SetID(uuid.New()).
		SetRating(5).
		AddUser(user).
		AddProduct(product).
		Save(ctx)
	require.NoError(t, err)

	// Delete review
	err = repo.Delete(ctx, createdReview.ID)
	require.NoError(t, err)

	// Verify review is deleted
	_, err = repo.FindByID(ctx, createdReview.ID)
	assert.Error(t, err)
}

